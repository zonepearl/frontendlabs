# C — The Complete Field Guide (Beginner → Expert)

> A practical, example-driven path through C: what the language actually
> gives you, why it's shaped the way it is, and how to read and write C that
> solves real problems — from a CLI tool to a TCP server to the internals of
> SQLite and Redis.
>
> Every concept ships with runnable code, a "Real-world example" or "War
> story," and — at the end of each Part — a full mini-project you build and
> run today. Read it once top to bottom; keep it as a lookup guide after.

---

## How to use this guide

- **Beginner (Part I):** the toolchain, types, control flow, functions,
  arrays/pointers, strings, structs, the preprocessor. The three trickiest
  sections (pointers, strings, structs/macros) each end with a **🔎
  Checkpoint** — a handful of predict-the-output/spot-the-bug/write-it-
  yourself questions with collapsible answers, so you can confirm you're
  ready before moving on instead of finding out three sections later. Part
  I then ends with two terminal mini-projects.
- **Intermediate (Part II):** dynamic memory, multi-file programs, file
  I/O, debugging/sanitizers, a minimal test harness. Ends with a vector
  library + CSV parser, and a JSON parser.
- **Advanced (Part III):** data structures from scratch, intrusive
  containers (the Linux-kernel trick), threads and atomics, sockets,
  `_Generic`. Ends with a thread pool and a multi-client chat server.
- **Real-world builds (Part IV):** production CLIs, a text editor, a Unix
  shell, a bump allocator.
- **Security-focused C (Part V):** undefined behavior, safe string
  handling, hardening flags, fuzzing, static analysis.
- **Open-source walkthroughs (Part VI):** SQLite, Redis, curl, cJSON/jq,
  Linux-kernel-style intrusive lists, uthash/klib — the same concepts you
  just learned, read out of real, battle-tested codebases.
- **Expert (Part VII):** the C11 memory model, writing a real allocator,
  ABI and API design, cross-compilation and build systems, profiling.
  Ends with a capstone: a mini Redis-like key-value store with a wire
  protocol.
- **Writing safe, leak-free C (Part VIII):** every memory-leak shape you'll
  actually hit in production (not just "forgot to `free`"), a defensive-
  programming checklist, and a bad-code/good-code catalog of the mistakes
  real C programmers make — read this even if you skip everything else.
- **Hardware & driver programming (Part IX):** the senior-engineer layer —
  talking to real hardware from user space (GPIO/I2C/SPI/serial),
  memory-mapped registers and `volatile`, writing and loading your first
  Linux kernel module, character device drivers, interrupts, and the
  concurrency rules that change once your code runs in kernel context.
  Ends with a capstone: one sensor, driven two ways.
- **Signals, non-local jumps, and mapped files (Part X):** the three
  general-purpose systems topics every "expert C" checklist expects that
  the rest of this guide only mentions in passing — `signal`/`sigaction`,
  `setjmp`/`longjmp`, and `mmap` for real files (as opposed to Part VII's
  anonymous allocator mapping or Part IX's hardware-register mapping).
- **Appendices:** C17→C23 feature diff, the UB cheat sheet, compiler flags,
  glossary, further reading.

Conventions:
- Code targets **C17 (ISO/IEC 9899:2018)** as the safe default and calls
  out **C23 (ISO/IEC 9899:2024)** features explicitly where used (`nullptr`,
  `typeof`, `#embed`, `constexpr`, `[[attributes]]`, `_BitInt`, binary
  literals). Toolchain used throughout: recent Clang/Apple Clang and GCC
  releases. Newer compilers accept `-std=c23`; some older releases use
  `-std=c2x` for the same in-progress standard mode. Where a snippet needs
  something GCC/Clang-specific beyond the standard (a builtin, an attribute),
  it's labeled.
- "**Real-world example**" boxes ground a concept in something you'll
  actually write. "**War story**" boxes describe a real class of incident
  — the kind that shows up in CVE databases and postmortems — caused by
  getting that concept wrong.
- Every code block is meant to compile. Where a snippet is *wrong on
  purpose* (to demonstrate a bug or UB), it's labeled `// BUG` or
  `// UB — do not do this`.

---

## Table of contents

**Part I — Foundations: the C way**
1. What C is, why it still matters, and which C you're learning
2. The toolchain: the four-stage pipeline, and your first program
3. Types and the sizes the standard actually guarantees
4. Operators, control flow, and the details beginners skip
5. Functions, the stack, and header files
6. Arrays, pointers, and the arithmetic that explains everything else
7. Strings in C: buffers, ownership, and why `gets()` was removed
8. Structs, unions, enums, `typedef`, and the preprocessor
9. Mini projects: a unit converter and a word-frequency counter

**Part II — Intermediate: memory, modules, and tooling**
10. Dynamic memory: the `malloc` family, ownership, and the bugs it enables
11. Pointers to pointers, function pointers, pointers to arrays
12. File I/O: buffering, binary vs. text mode, `errno`
13. Multi-file programs: headers, `static`, translation units, and `make`
14. Debugging and sanitizers: gdb/lldb, Valgrind, ASan/UBSan
15. A minimal test harness for C
16. Mini projects: a dynamic-array library + CSV parser, and a JSON parser

**Part III — Advanced: data structures, concurrency, networking**
17. Bit manipulation and bit-fields
18. Data structures from scratch: linked lists, stacks/queues, hash tables, trees
19. Intrusive containers: the Linux kernel's `container_of` trick
20. Concurrency: POSIX threads, mutexes/condvars, C11 `<threads.h>` and atomics
21. Sockets: a TCP echo server, then a tiny HTTP server
22. Variadic functions and `_Generic`
23. Mini projects: a thread pool, and a `select()`-based chat server

**Part IV — Real-world builds**
24. Building production CLIs: argument parsing, config precedence, exit codes
25. Capstone: a `kilo`-style terminal text editor
26. Capstone: a Unix mini-shell with pipes and redirection
27. Capstone: a bump/arena allocator

**Part V — Security-focused C**
28. Undefined behavior: the list every C programmer must memorize
29. Buffer overflows, safe string handling, and CERT C rules
30. Hardening: compiler flags, canaries, ASLR, `_FORTIFY_SOURCE`
31. Fuzzing a real parsing bug with libFuzzer
32. Static analysis: clang-tidy, cppcheck, the Clang Static Analyzer

**Part VI — Open-source package walkthroughs**
33. SQLite: the amalgamation, the VFS layer, opcode-based execution
34. Redis: the `ae.c` event loop, `sds` dynamic strings, object encoding
35. curl: the easy/multi handle design and callback-driven I/O
36. cJSON: recursive-descent parsing and tagged unions
37. uthash and klib: generics via macros

**Part VII — Expert: internals, allocators, API design**
38. The C11 memory model: atomics, ordering, lock-free basics
39. Writing a real allocator: free lists, coalescing, and `dlmalloc`'s ideas
40. ABI, linking, and designing a stable C API
41. Cross-compilation and build systems at scale (CMake/Meson)
42. Profiling: `perf`, cache-aware layout, and false sharing
43. Capstone: a mini Redis-like key-value store with a wire protocol

**Part VIII — Writing safe, leak-free C**
44. Memory leaks, in depth: every shape they actually take in production
45. Writing safe C: a defensive-programming checklist
46. How NOT to write C: a bad-code/good-code anti-pattern catalog

**Part IX — Hardware & driver programming**
47. Two worlds: user-space hardware I/O vs. writing a kernel driver
48. Talking to hardware from user space: GPIO, I2C, SPI, and serial
49. Memory-mapped registers: `volatile`, `/dev/mem`, and memory barriers
50. Your first Linux kernel module: hello world, the Makefile, insmod/dmesg
51. Character device drivers: `file_operations`, `cdev`, and `copy_to/from_user`
52. Interrupts and kernel concurrency: top/bottom halves, spinlocks vs. mutexes
53. Debugging drivers: `printk`/dmesg, an oops backtrace, and QEMU+gdb
54. Capstone: one sensor, driven two ways — user-space ioctl vs. a kernel driver

**Part X — Signals, non-local jumps, and memory-mapped files**
55. Signals: handling asynchronous events safely
56. Non-local jumps: `setjmp`/`longjmp` and why they're not exceptions
57. Memory-mapped files: `mmap` for real files, `msync`, and shared mappings

**Appendices**
- A. C17 → C23: what's new, feature by feature
- B. Undefined-behavior and pitfalls cheat sheet
- C. Compiler-flags cheat sheet
- D. Glossary
- E. Further reading and real codebases worth reading

---

# Part I — Foundations: the C way

> The language itself, in order: the toolchain, then every piece a C program
> is made of — types, control flow, functions, arrays and pointers, strings,
> structs, and the preprocessor. Three checkpoint sections (pointers,
> strings, structs/macros) let you confirm you're solid before moving on
> instead of finding out three sections later. Two terminal mini-projects
> close out the Part before Part II adds dynamic memory on top.

## 1. What C is, why it still matters, and which C you're learning

### 1.1 The problem C solved, and still solves

C was designed by Dennis Ritchie at Bell Labs (1969–1973) to rewrite Unix
in something more portable than assembly. Its entire design philosophy is
one sentence: **give the programmer direct, predictable control over
memory and machine instructions, and trust them with it.** Nothing is
inserted behind your back — no garbage collector, no hidden bounds checks,
no implicit allocations, no exceptions unwinding a stack you didn't ask to
unwind. Every cost is visible in the source.

That's why, 50+ years later, C is still common near the bottom of the
stack: Linux and many embedded kernels; important parts of Windows;
language runtimes such as CPython, the JVM, V8, and Go's early runtime;
SQLite; Redis; curl; OpenSSL; git's core; PostgreSQL. If a program needs
to run close to the hardware, on a tiny device, or as a reference
implementation other languages bind against, C is often one of the first
languages considered.

### 1.2 What C deliberately does not have

No classes, no exceptions, no garbage collector, no built-in strings
(a "string" is a convention: a pointer to `char` plus a `'\0'`
terminator), no bounds checking on arrays, no function overloading (until
`_Generic` in C11, which is a macro trick, not real overloading), no
namespaces. This is not primitiveness — it's the reason a C ABI is the
*lingua franca* every other language calls into (Python's `ctypes`, Rust's
`extern "C"`, Go's `cgo`, Java's JNI all speak C calling conventions,
because nothing simpler exists to agree on).

### 1.3 Which C you're learning: standards and what changed

| Standard | Year | What it added that you'll actually use |
|---|---|---|
| C89/C90 (ANSI C) | 1989 | The baseline everyone assumes: function prototypes, `void*`, `const` |
| C99 | 1999 | `//` comments, `stdint.h` fixed-width types, variable-length arrays (VLAs), designated initializers, `for (int i = ...)`, `inline`, `_Bool`, compound literals |
| C11 | 2011 | `<stdatomic.h>`, `<threads.h>`, `_Generic`, `_Static_assert`, anonymous structs/unions, `_Noreturn` |
| C17 | 2018 | Bugfix release, no new features — **the safest "modern C" target today** |
| **C23** | 2024 | `nullptr`, `bool`/`true`/`false` as keywords (no `<stdbool.h>` needed), `constexpr`, `typeof`/`typeof_unqual`, `#embed`, `[[attributes]]` (`[[nodiscard]]`, `[[maybe_unused]]`, `[[deprecated]]`), binary literals (`0b1010`), digit separators (`1'000'000`), `_BitInt(N)`, removal of implicit-`int` and K&R-style function definitions |

**Which one should you target?** C17 for anything that must compile on
older distros (RHEL, embedded toolchains) or interop with a wide range of
compilers. C23 when you control the toolchain — it removes several classes
of footgun (`nullptr` vs. `NULL`'s type ambiguity, `bool` finally being a
real keyword) with zero downside on GCC 13+/Clang 17+. This guide defaults
to C17 syntax and calls out C23 alternatives inline.

```bash
# Check what your compiler actually supports
clang --version
gcc --version
clang -std=c23 -dM -E -x c /dev/null | grep __STDC_VERSION__
# 202311L  == C23. 201710L == C17. 201112L == C11.
```

### 1.4 War story

CVE-2014-0160 — **Heartbleed** — was a missing bounds check in OpenSSL's
TLS heartbeat extension: the code read `payload_length` bytes from a
buffer using a length supplied by the *attacker*, with no check that the
buffer actually contained that many bytes. One `memcpy` with an
attacker-controlled length leaked private keys, session cookies, and
passwords from millions of servers for two years before anyone noticed.
The fix was three lines. The lesson that shaped an entire industry's
approach to C (Part V of this guide) was: **in C, the compiler will never
save you from reading past a buffer — only you, your tests, your
sanitizers, and your code review can.**

---

## 2. The toolchain: the four-stage pipeline, and your first program

### 2.1 The four stages, made visible

"Compiling" a C file is actually four separate steps, and understanding
them explains almost every confusing compiler error you'll ever see.

```bash
cat > hello.c <<'EOF'
#include <stdio.h>

int main(void) {
    printf("hello, c\n");
    return 0;
}
EOF

# 1. Preprocessing: expand #include, #define, #if — pure text substitution
clang -E hello.c -o hello.i
wc -l hello.i          # thousands of lines, all from stdio.h's chain of includes

# 2. Compilation: C source -> assembly for the target architecture
clang -S hello.i -o hello.s
head -20 hello.s

# 3. Assembly: assembly -> machine code, produces an object file
clang -c hello.s -o hello.o
file hello.o            # Mach-O 64-bit object arm64 (or ELF 64-bit on Linux)

# 4. Linking: resolve symbols (printf lives in libc, not your object file)
clang hello.o -o hello
./hello                 # hello, c

# In practice you just run:
clang -std=c17 -Wall -Wextra -o hello hello.c
```

**Real-world example.** "Undefined reference to `foo`" is always a
*linking* error — the compiler understood your code fine, but no object
file or library provides a symbol named `foo`. "Implicit declaration of
function `foo`" is a *compilation* warning from stage 2 — you called
something the compiler hasn't seen a prototype for (usually a missing
`#include`). Knowing which stage produced an error tells you where to
look: linker errors mean check your `-l`/`-L` flags and whether you
compiled every `.c` file; compiler errors mean check your headers and
syntax.

### 2.2 The flags you should always use

```bash
clang -std=c17 -Wall -Wextra -Wpedantic -Werror \
      -fsanitize=address,undefined -g -O0 \
      -o prog prog.c
```

| Flag | Why |
|---|---|
| `-std=c17` (or `-std=c23`) | Pin the language version; don't rely on "whatever GNU extensions the compiler defaults to" |
| `-Wall -Wextra` | Turns on the warnings that catch real bugs: uninitialized reads, signed/unsigned comparison, unused variables |
| `-Wpedantic` | Warns on non-standard extensions — important if you need portability |
| `-Werror` | Warnings fail the build. Excellent for CI and most project code; you may relax it for exploratory/compiler-portability builds |
| `-fsanitize=address,undefined` | Catches buffer overflows, use-after-free, and UB at *runtime* — Part V covers this in depth |
| `-g -O0` | Debug symbols, no optimization, so a debugger's line numbers match your source |

**A portability trap `-std=c17` sets for you on Linux.** Strict `-std=c17`
defines `__STRICT_ANSI__`, and on glibc that *hides* POSIX declarations
(`strdup`, `fork`, `execvp`, `getopt`, `pthread_*`'s POSIX-only helpers,
`pipe`, `dup2`, BSD socket calls) unless a feature-test macro is set —
so code that calls them fails with "implicit declaration of function"
under `-Werror`, even though the exact same command compiles cleanly on
macOS (Apple's libc doesn't gate these the same way, which is why every
snippet in this guide compiled clean on the machine it was tested on).
Any time a later section's code touches POSIX (Part III's threads and
sockets, Part IV's `fork`/`termios`/`getopt`), compile it with either
`-std=gnu17` (GNU's C17, POSIX/BSD extensions visible by default) or
`-std=c17 -D_POSIX_C_SOURCE=200809L -D_DEFAULT_SOURCE` on Linux. This
guide keeps `-std=c17` in the command lines for the pure-ISO-C sections
(Parts I–II) and switches to `-std=gnu17` from Part III onward for
exactly this reason.

### 2.3 `make`, at the level you actually need on day one

```makefile
CC      = clang
CFLAGS  = -std=c17 -Wall -Wextra -Wpedantic -Werror -g
SRC     = main.c
BIN     = prog

$(BIN): $(SRC)
	$(CC) $(CFLAGS) -o $(BIN) $(SRC)

clean:
	rm -f $(BIN)

.PHONY: clean
```

Section 13 goes deep on multi-file `make`; this is enough to stop typing
the full `clang` command by hand.

---

## 3. Types and the sizes the standard actually guarantees

### 3.1 The trap: "int is 4 bytes" is a habit, not a guarantee

The C standard only guarantees *minimum* ranges, not exact sizes:
`char` ≥ 8 bits, `short`/`int` ≥ 16 bits, `long` ≥ 32 bits, `long long` ≥
64 bits, and `sizeof(char) ≤ sizeof(short) ≤ sizeof(int) ≤ sizeof(long) ≤
sizeof(long long)`. On the platforms you use daily this happens to be 1/2/
4/8/8 bytes, but code that assumes `int` is exactly 32 bits everywhere has
already broken on real hardware (16-bit `int` on some embedded targets;
`long` is 4 bytes on Windows LLP64 but 8 bytes on Linux/macOS LP64 for the
*same* "64-bit" build).

**The fix, always: use `<stdint.h>` for anything where the size matters.**

```c
#include <stdint.h>
#include <stdio.h>

int main(void) {
    int32_t  a = 42;          // exactly 32 bits, everywhere, or it fails to compile
    uint64_t b = 1ULL << 40;  // exactly 64 bits, unsigned
    size_t   n = sizeof(a);   // the *correct* type for sizes/counts — unsigned, platform width
    intptr_t p;               // an integer wide enough to hold a pointer, for pointer arithmetic tricks

    printf("%d %llu %zu\n", a, (unsigned long long)b, n);
    return 0;
}
```

| Type | Guarantee | Use it for |
|---|---|---|
| `int8_t`/`uint8_t` … `int64_t`/`uint64_t` | exact width | wire protocols, file formats, anywhere size matters |
| `size_t` | unsigned, wide enough for any object's size | array indices, `sizeof`, loop counters over memory |
| `ssize_t` (POSIX) | signed version of `size_t` | return values that need a `-1` error sentinel (`read()`, `write()`) |
| `ptrdiff_t` | signed, holds the result of pointer subtraction | pointer arithmetic results |
| `intptr_t`/`uintptr_t` | can hold a `void*` round-trip | rare — casting pointers to integers for bit tricks or logging |

### 3.2 `printf`/`scanf` format specifiers, correctly

```c
#include <stdint.h>
#include <inttypes.h>
#include <stdio.h>

int main(void) {
    size_t   n = 100;
    int64_t  big = -123456789012345;
    uint32_t u = 4000000000u;

    printf("n=%zu big=%" PRId64 " u=%" PRIu32 "\n", n, big, u);
    // %zu for size_t, PRId64/PRIu32 macros from <inttypes.h> for fixed-width types —
    // never guess "%ld" for int64_t; it's `long` on Linux but `long long` on Windows.
    return 0;
}
```

**War story.** Mismatched `printf` format specifiers are UB, not just
"ugly" — `printf("%d", some_int64)` reads only 4 bytes off a variadic
argument that's actually 8, so it *sometimes* reads the second half of
your argument, garbage stack memory, or the wrong subsequent argument
entirely, depending on the ABI's argument-passing register layout on that
call. `-Wformat` (in `-Wall`) catches this at compile time — never disable
it.

### 3.3 Signed integer overflow is undefined behavior — unsigned is not

```c
#include <limits.h>
#include <stdio.h>

int main(void) {
    int x = INT_MAX;
    // int y = x + 1;      // UB — compiler may assume this never happens and
                            // optimize based on that assumption (see Part V)
    unsigned int u = UINT_MAX;
    unsigned int v = u + 1; // well-defined: wraps to 0, per modular arithmetic rules
    printf("%u\n", v);      // 0
    return 0;
}
```

This single asymmetry (signed overflow = UB, unsigned overflow = defined
wraparound) is behind a large fraction of real-world security bugs and is
covered in depth in Part V.

---

## 4. Operators, control flow, and the details beginners skip

### 4.1 `switch` falls through — on purpose, and that's both a feature and a trap

```c
#include <stdio.h>

void describe(int day) {
    switch (day) {
        case 6:
        case 7:
            printf("weekend\n");
            break;
        case 1: case 2: case 3: case 4: case 5:
            printf("weekday\n");
            break;
        default:
            printf("invalid\n");
    }
}
```

Forgetting a `break` is one of the most common beginner bugs; `-Wextra`
warns on implicit fall-through unless you mark it intentional (C17:
a comment; C23: `[[fallthrough]];`).

```c
switch (x) {
    case 1:
        do_thing();
        [[fallthrough]];   // C23 attribute — silences the warning, documents intent
    case 2:
        do_other_thing();
        break;
}
```

### 4.2 `goto` — the feature everyone tells you to avoid, that real C code uses constantly

C has no exceptions and no RAII. The idiomatic way to clean up multiple
resources on an error path *without* deeply nested `if`s is a single
forward `goto` to a cleanup label — this is the dominant pattern in the
Linux kernel, SQLite, and most production C:

```c
int process_file(const char *path) {
    FILE *f = NULL;
    char *buf = NULL;
    int rc = -1;

    f = fopen(path, "r");
    if (!f) goto cleanup;

    buf = malloc(4096);
    if (!buf) goto cleanup;

    if (fread(buf, 1, 4096, f) == 0) goto cleanup;

    rc = 0; // success
cleanup:
    free(buf);
    if (f) fclose(f);
    return rc;
}
```

This is *not* the "goto considered harmful" spaghetti Dijkstra warned
about — it only ever jumps forward, to one place, and it's the closest
thing C has to a `finally` block. Part IV and VI show this exact pattern
throughout real codebases.

### 4.3 Short-circuit evaluation is a control-flow tool, not just a boolean trick

```c
// Real-world example: guard a dereference
if (ptr != NULL && ptr->value > 0) {  // right side never evaluates if ptr is NULL
    use(ptr);
}

// Common idiom: use || for "do this or bail"
FILE *f = fopen(path, "r");
if (!f) {
    perror("fopen");
    return 1;
}
```

---

## 5. Functions, the stack, and header files

### 5.1 What actually happens on a function call

```c
#include <stdio.h>

int add(int a, int b) {
    int result = a + b;   // lives in add's stack frame
    return result;
}

int main(void) {
    int sum = add(3, 4);  // a new stack frame is pushed, args copied in (by value),
                           // frame popped on return, result copied out
    printf("%d\n", sum);
    return 0;
}
```

Every local variable and every argument (in C, **all arguments are passed
by value** — there is no pass-by-reference; you simulate it by passing a
pointer) lives on the **stack**, a region that grows and shrinks
automatically with function calls. It's fast (just moving a stack
pointer) but finite (typically 1–8MB per thread) and disappears the
instant the function returns — which is exactly why returning a pointer to
a local variable is one of the first bugs every C programmer writes:

```c
// BUG — undefined behavior
int *make_answer(void) {
    int answer = 42;
    return &answer;   // answer's stack frame is gone the moment this returns
}
```

Clang/GCC with `-Wall` will warn: `warning: address of stack memory
associated with local variable 'answer' returned`. The fix is either to
return by value, take an output pointer from the caller, or heap-allocate
(Section 10).

### 5.2 Header files: declarations, not definitions

```c
// point.h
#ifndef POINT_H
#define POINT_H

typedef struct {
    double x, y;
} Point;

double point_distance(Point a, Point b);

#endif // POINT_H
```

```c
// point.c
#include "point.h"
#include <math.h>

double point_distance(Point a, Point b) {
    double dx = a.x - b.x, dy = a.y - b.y;
    return sqrt(dx * dx + dy * dy);
}
```

The `#ifndef`/`#define`/`#endif` triple is an **include guard** — it stops
a header being pasted twice into one translation unit if two `.c` files
both include something that includes `point.h`. C23 also supports
`#pragma once`, which every modern compiler already accepted as an
extension; the include guard remains the portable choice for library code.

---

## 6. Arrays, pointers, and the arithmetic that explains everything else

### 6.1 An array *is* a contiguous block; a pointer is *not* an array

```c
#include <stdio.h>

int main(void) {
    int arr[5] = {10, 20, 30, 40, 50};
    int *p = arr;              // arrays decay to a pointer to their first element

    printf("%d %d\n", arr[2], *(p + 2));   // identical: 30 30
    printf("%zu\n", sizeof(arr));          // 20 — sizeof knows arr is 5 ints
    printf("%zu\n", sizeof(p));            // 8  — sizeof(p) is just a pointer's size

    return 0;
}
```

**This is the single most important fact in C**: `arr[i]` is defined as
`*(arr + i)`. Pointer arithmetic is scaled by the pointee's type — `p + 1`
moves forward `sizeof(*p)` bytes, not 1 byte. This is why `void*`
arithmetic isn't standard C (the compiler doesn't know the scale) and why
`sizeof` on a decayed pointer is a classic bug: the moment an array is
passed to a function, it decays to a pointer and `sizeof` inside that
function returns the pointer's size (8), not the array's.

```c
void bad_len(int arr[]) {
    printf("%zu\n", sizeof(arr));  // 8, NOT the caller's array size — classic beginner bug
}
```

### 6.2 Pointer arithmetic, walked through

```c
#include <stdio.h>

int main(void) {
    int nums[] = {1, 2, 3, 4, 5};
    int *begin = nums;
    int *end = nums + 5;   // one past the last element — legal to compute, not to dereference

    for (int *p = begin; p != end; p++) {
        printf("%d ", *p);
    }
    printf("\n");

    ptrdiff_t count = end - begin;  // 5 — subtracting pointers gives element count
    printf("%td\n", count);
    return 0;
}
```

This `begin`/`end` idiom is exactly how C++ iterators and Go slices are
modeled under the hood — C just makes the pointer arithmetic explicit.

### 6.3 `const` correctness with pointers — read right to left

```c
int x = 5;
const int *p1 = &x;        // pointer to const int:   *p1 = 6 is an error, p1 = &y is fine
int *const p2 = &x;        // const pointer to int:    *p2 = 6 is fine, p2 = &y is an error
const int *const p3 = &x;  // const pointer to const int: neither is fine
```

**Real-world example.** Function signatures use `const` to document (and
have the compiler enforce) "this function reads but does not modify" —
this is why almost every standard-library string function that only reads
takes `const char *`:

```c
size_t my_strlen(const char *s) {   // promise: won't write through s
    size_t n = 0;
    while (s[n] != '\0') n++;
    return n;
}
```

### 🔎 Checkpoint: pointers and arrays

Answer these before moving on — this section is the one most beginners
need to revisit, and the fastest way to know if you actually need to is
to try these first.

**1. Predict the output**, then check yourself:
```c
#include <stdio.h>
int main(void) {
    int arr[4] = {10, 20, 30, 40};
    int *p = arr + 1;
    printf("%d %d\n", *p, p[1]);
    return 0;
}
```
<details><summary>Answer</summary>

`20 30`. `p` points at `arr[1]` (value `20`), so `*p` is `20`. `p[1]` is
defined as `*(p + 1)`, and `p` is already `arr + 1`, so `p[1]` reaches
`arr[2]`, which is `30`.
</details>

**2. Spot the bug:**
```c
void print_len(int arr[]) {
    printf("%zu\n", sizeof(arr) / sizeof(arr[0]));
}
```
<details><summary>Answer</summary>

`arr` decayed to `int *` the moment it crossed the function boundary
(§6.1), so `sizeof(arr)` is the pointer's size (8 on a 64-bit machine),
not the caller's array size — this division computes `8 / 4 = 2`,
always, no matter how big the array actually was. There is no fix that
keeps this signature; the only correct fix is to pass the length in
separately: `void print_len(int *arr, size_t len)`. This exact mistake
is real enough that Clang/GCC warn about it by default (§46.8).
</details>

**3. Write it yourself:** write `int sum_array(const int *arr, size_t len)`
that returns the sum of `len` elements, then call it correctly on a
stack-allocated array.
<details><summary>One correct answer</summary>

```c
#include <stdio.h>

int sum_array(const int *arr, size_t len) {
    int total = 0;
    for (size_t i = 0; i < len; i++) total += arr[i];
    return total;
}

int main(void) {
    int nums[] = {1, 2, 3, 4, 5};
    printf("%d\n", sum_array(nums, sizeof(nums) / sizeof(nums[0])));  // 15
    return 0;
}
```
Note `sizeof(nums) / sizeof(nums[0])` at the *call site* — this only
works here because `nums` hasn't decayed yet; it's a real array in the
scope where `sizeof` sees it. Try moving that computation into a helper
function that takes `int nums[]` as a parameter and predict what breaks.
</details>

**4. Conceptual — why does this compile, and what does it print?**
```c
#include <stdio.h>
int main(void) {
    int arr[4] = {10, 20, 30, 40};
    printf("%d\n", 2[arr]);   // yes, this is legal C — see the answer
    return 0;
}
```
<details><summary>Answer</summary>

It prints `30`, identically to `arr[2]`. `arr[2]` is defined as
`*(arr + 2)`; addition is commutative, so `*(arr + 2)` and `*(2 + arr)`
are the same expression, and `2[arr]` means exactly `*(2 + arr)`. This
isn't a special case in the grammar — `a[b]` *is* `*(a + b)` for any
pointer/integer pair, in either order. It's a well-known trivia question
precisely because it exposes how literally the array-indexing rule from
§6.1 should be taken.
</details>

---

## 7. Strings in C: buffers, ownership, and why `gets()` was removed

### 7.1 A "string" is a convention, not a type

```c
#include <string.h>
#include <stdio.h>

int main(void) {
    char greeting[20] = "hello";   // 6 bytes used: 'h''e''l''l''o''\0', 14 unused
    printf("%zu %zu\n", strlen(greeting), sizeof(greeting)); // 5  20

    strcat(greeting, ", world");   // caller's job to ensure the buffer is big enough
    printf("%s\n", greeting);
    return 0;
}
```

`strlen` walks memory until it finds `'\0'` — there is no length prefix,
no bounds. Every buffer-related CVE in C's history traces back to this one
design decision, made in 1972 to save a byte per string on a PDP-11.

### 7.2 The functions to prefer, and why

| Avoid | Because | Prefer |
|---|---|---|
| `gets()` | No way to bound the read; **removed from the C11 standard entirely** | `fgets(buf, sizeof buf, stdin)` |
| `strcpy(dst, src)` | No bound on `src`'s length vs. `dst`'s capacity | `strlcpy` (BSD/macOS) or `snprintf(dst, n, "%s", src)` |
| `strcat(dst, src)` | Same problem, compounded (walks `dst` first) | `strlcat`, or track length yourself |
| `sprintf(buf, ...)` | No bound on output length | `snprintf(buf, sizeof buf, ...)` — always |
| `scanf("%s", buf)` | Unbounded read into `buf` | `scanf("%19s", buf)` (explicit width) or `fgets` + parse |

```c
#include <stdio.h>
#include <string.h>

int main(void) {
    char name[32];
    fgets(name, sizeof name, stdin);          // safe: stops at 31 chars + '\0'
    name[strcspn(name, "\n")] = '\0';         // fgets keeps the trailing newline — strip it

    char msg[64];
    snprintf(msg, sizeof msg, "hello, %s!", name);  // always bounded, always '\0'-terminated
    printf("%s\n", msg);
    return 0;
}
```

**War story.** `strcpy`/`strcat`/`sprintf` without bounds checking are
exactly the family of functions behind decades of stack-buffer-overflow
CVEs, starting with the 1988 Morris Worm (which exploited `gets()` in
`fingerd`) through to modern IoT firmware CVEs today. `-D_FORTIFY_SOURCE=2`
(Section 30) adds runtime bounds checks to these calls when the compiler
can determine buffer sizes, but the real fix is: **default to the bounded
variant by default.**

### 7.3 C23: better ergonomics, same underlying model

C23 doesn't change the string model (still no built-in string type), but
it removes friction: `nullptr` replaces the type-unsafe `NULL` macro,
`typeof` lets you write generic string helpers without `_Generic`
boilerplate, and `[[nodiscard]]` on a function like a custom
`safe_strdup` forces callers to check the returned pointer for
allocation failure.

```c
// C23
#include <stdlib.h>
#include <string.h>

[[nodiscard]] char *safe_strdup(const char *s) {
    size_t n = strlen(s) + 1;
    char *copy = malloc(n);
    if (copy) memcpy(copy, s, n);
    return copy;
}
```

### 🔎 Checkpoint: strings and buffers

**1. Predict what happens** — will these two both "work," and if not, why
does one look fine while quietly being wrong?
```c
#include <string.h>
int main(void) {
    char buf6[6];
    strcpy(buf6, "hello");   // "hello" is 5 letters + '\0' = 6 bytes

    char buf5[5];
    strcpy(buf5, "hello");   // same source string, smaller destination
    return 0;
}
```
<details><summary>Answer</summary>

The first is correct: `"hello"` needs exactly 6 bytes (`h e l l o \0`)
and `buf6` has exactly 6. The second overflows `buf5` by one byte.

Try compiling this with this guide's own recommended flags from §2.2
(`-Wall -Wextra -Werror`) and something interesting happens: **it
doesn't compile at all.** On a modern hardened toolchain (Apple's Clang,
or GCC/Clang with glibc's `_FORTIFY_SOURCE` — §30), the compiler can see
at compile time that a 6-byte literal can't fit in a 5-byte array, warns
`-Wfortify-source`, and `-Werror` turns that into a hard build failure
before the buggy code ever runs. Drop `-Werror` and it becomes just a
warning — the code compiles, but `strcpy` has been silently rewritten
into a checked `__strcpy_chk` call that **aborts the process
immediately at runtime** with a diagnostic, instead of letting the
overflow happen. Either way, this exact overflow never gets a chance to
silently corrupt memory — but only because the compiler could prove it
at compile time. The classic "silently corrupts whatever memory sits
after `buf5`, might even look fine in testing" danger from §29.1 is
still completely real the moment the source string comes from somewhere
the compiler can't see in advance, like user input:
`strcpy(buf5, some_variable)` — no compiler on earth can prove that one
safe or unsafe ahead of time, so no such protection kicks in. Same bug,
same fix (`snprintf`/`strlcpy`, §7.2) — this compile-time special case
just happens to be caught for you,
where a runtime-determined one still isn't.
</details>

**2. Spot the bug:**
```c
#include <stdio.h>
#include <string.h>
int main(void) {
    char line[256];
    fgets(line, sizeof line, stdin);
    if (strcmp(line, "quit") == 0) {
        printf("goodbye\n");
    }
    return 0;
}
```
The user types `quit` and presses enter, but `"goodbye"` never prints. Why?
<details><summary>Answer</summary>

`fgets` keeps the trailing newline in the buffer (§7.2's `fgets` note) —
`line` actually holds `"quit\n"`, not `"quit"`, so `strcmp` never
matches. The fix is `line[strcspn(line, "\n")] = '\0';` right after the
`fgets` call, exactly as shown in §7.2's `snprintf` example.
</details>

**3. Write it yourself:** write `void safe_concat(char *dst, size_t
dst_cap, const char *a, const char *b)` that writes `a` followed by `b`
into `dst`, truncating safely if `dst` is too small — never overflowing,
no matter how long `a`/`b` are.
<details><summary>One correct answer</summary>

```c
#include <stdio.h>

void safe_concat(char *dst, size_t dst_cap, const char *a, const char *b) {
    snprintf(dst, dst_cap, "%s%s", a, b);   // snprintf always bounds to dst_cap, truncating cleanly
}

int main(void) {
    char buf[16];
    safe_concat(buf, sizeof buf, "hello, ", "world! this part gets truncated safely");
    printf("%s\n", buf);
    return 0;
}
```
This is the entire lesson of §7.2 in one line: reach for `snprintf` by
default any time you're building a string into a fixed buffer, instead
of reasoning about `strcat`/`strcpy`'s lengths by hand.
</details>

**4. Conceptual:** why can't you use `sizeof` to get the length of a
`char *` variable the way you can on a `char[]` array — and what's the
one function that *does* give you that length?
<details><summary>Answer</summary>

`sizeof` on a pointer variable gives the pointer's own size (8 bytes on
a 64-bit machine) — the same array-decay fact from §6.1's checkpoint,
not the length of what it points to. C has no way to recover "how many
bytes were allocated here" from a bare pointer; the only way to know a
C string's length is `strlen`, which walks memory until it finds the
`'\0'` terminator (§7.1) — an `O(n)` scan, not a stored value, unlike
Redis's `sds` strings in §34, which keep a length field precisely to
avoid this.
</details>

---

## 8. Structs, unions, enums, `typedef`, and the preprocessor

### 8.1 Structs: memory layout you can reason about

```c
#include <stdio.h>
#include <stddef.h>

typedef struct {
    char  kind;     // 1 byte
    int   value;    // 4 bytes
    short flags;    // 2 bytes
} Item;             // sizeof is NOT 1+4+2=7 — see below

int main(void) {
    printf("sizeof(Item) = %zu\n", sizeof(Item));               // likely 12, not 7
    printf("offsetof value = %zu\n", offsetof(Item, value));    // likely 4, not 1
    return 0;
}
```

The compiler inserts **padding** so each field lands on an address that's
a multiple of its own alignment requirement (an `int` usually needs
4-byte alignment). Reordering fields largest-to-smallest often shrinks a
struct:

```c
typedef struct {
    int   value;   // 4 bytes, 4-byte aligned
    short flags;   // 2 bytes
    char  kind;    // 1 byte
    // 1 byte padding to round the whole struct to a multiple of 4
} ItemPacked;       // sizeof == 8, not 12 — same three fields, better layout
```

**Real-world example.** This is exactly the optimization database and
game engines do to arrays of millions of structs: reordering fields from
worst-case to best-case padding can shrink memory footprint by 30%+ and
improve cache-line utilization (Section 42 goes deep on this).

### 8.2 Unions: one block of memory, many interpretations

```c
#include <stdio.h>
#include <stdint.h>

typedef union {
    uint32_t as_int;
    uint8_t  as_bytes[4];
} Word;

int main(void) {
    Word w;
    w.as_int = 0x01020304;
    // On a little-endian machine (x86, Apple Silicon in LE mode):
    printf("%02x %02x %02x %02x\n",
           w.as_bytes[0], w.as_bytes[1], w.as_bytes[2], w.as_bytes[3]);
    // prints 04 03 02 01 — this IS how you detect endianness in real code
    return 0;
}
```

This union-based byte inspection is legal, portable *behavior*-wise (type
punning through a union is explicitly allowed in C, unlike C++), and is
exactly how real network code and file-format parsers implement portable
byte-order detection and manual serialization.

### 8.3 Tagged unions: C's answer to "sum types" / enums-with-data

```c
typedef enum { SHAPE_CIRCLE, SHAPE_RECT } ShapeKind;

typedef struct {
    ShapeKind kind;
    union {
        struct { double radius; } circle;
        struct { double w, h; } rect;
    };  // C11 anonymous union — no member name needed to reach circle/rect
} Shape;

double area(const Shape *s) {
    switch (s->kind) {
        case SHAPE_CIRCLE: return 3.14159265 * s->circle.radius * s->circle.radius;
        case SHAPE_RECT:   return s->rect.w * s->rect.h;
    }
    return 0.0;
}
```

This *is* the pattern Part VI shows powering cJSON's `cJSON` struct and
countless AST/value representations in real parsers and interpreters —
languages with real sum types (Rust `enum`, Swift `enum`, Haskell ADTs)
compile down to exactly this layout.

### 8.4 The preprocessor: macros are text substitution, not functions

```c
#define SQUARE(x) ((x) * (x))     // parenthesize EVERYTHING — always

int main(void) {
    int a = 5;
    int b = SQUARE(a + 1);        // expands to ((a + 1) * (a + 1)) = 36, correct
    // Without the inner parens: ((a + 1 * a + 1)) = wrong answer
    (void)b;

    int i = 2;
    int bad = SQUARE(i++);        // ((i++) * (i++)) — undefined behavior: i modified twice
                                   // between sequence points. NEVER pass side effects to macros.
    (void)bad;
    return 0;
}
```

Prefer `static inline` functions over macros wherever you don't need
text-substitution tricks — they're type-checked, debuggable, and don't
double-evaluate arguments:

```c
static inline int square(int x) { return x * x; }  // safe, and the compiler will still inline it
```

### 🔎 Checkpoint: structs, unions, and macros

**1. Predict `sizeof`, then reorder to shrink it:**
```c
#include <stdio.h>
struct A { char a; double b; char c; };
struct B { double b; char a; char c; };   // same three fields, different order

int main(void) {
    printf("%zu %zu\n", sizeof(struct A), sizeof(struct B));
    return 0;
}
```
<details><summary>Answer</summary>

`sizeof(struct A)` is `24`: `a` (1 byte) needs 7 bytes of padding before
`b` can start on an 8-byte boundary, then `b` (8 bytes), then `c`
(1 byte) needs 7 more bytes of padding so the *whole struct's* size is a
multiple of its largest member's alignment (8). `sizeof(struct B)` is
`16`: putting the 8-byte `double` first means `a` and `c` (1 byte each)
can share the remaining space with only 6 bytes of padding at the very
end. Same three fields, 8 bytes saved per struct just from reordering —
exactly §8.1's "largest to smallest" rule, and exactly the optimization
§42.2 scales up to millions of structs.
</details>

**2. Spot the bug** — what does this actually print, and why is it
probably not what you'd guess?
```c
#include <stdio.h>
#define MIN(a,b) a < b ? a : b

int main(void) {
    int x = MIN(1,2) + 3;
    printf("%d\n", x);
    return 0;
}
```
<details><summary>Answer</summary>

It prints `1`, not `4`. `MIN(1,2) + 3` expands, by pure text
substitution (§8.4), to `1 < 2 ? 1 : 2 + 3`. The ternary operator `?:`
binds *looser* than `+`, so this parses as `(1 < 2) ? 1 : (2 + 3)` — the
condition is true, so the whole expression is `1`, and the `+ 3` silently
attached itself to the wrong branch. The fix is the same lesson as
§46.10: parenthesize the *entire* macro body, not just each parameter:
`#define MIN(a,b) ((a) < (b) ? (a) : (b))`.
</details>

**3. Write it yourself:** using §8.3's tagged-union pattern, define a
`Value` type that can hold an `int`, a `float`, or a `const char *`,
tagged by a `ValueKind` enum, plus a `print_value` function that
switches on the tag.
<details><summary>One correct answer</summary>

```c
#include <stdio.h>

typedef enum { VAL_INT, VAL_FLOAT, VAL_STRING } ValueKind;

typedef struct {
    ValueKind kind;
    union {
        int intval;
        float floatval;
        const char *strval;
    };
} Value;

void print_value(const Value *v) {
    switch (v->kind) {
        case VAL_INT:    printf("int: %d\n", v->intval); break;
        case VAL_FLOAT:  printf("float: %f\n", v->floatval); break;
        case VAL_STRING: printf("string: %s\n", v->strval); break;
    }
}

int main(void) {
    Value a = { .kind = VAL_INT, .intval = 42 };
    Value b = { .kind = VAL_STRING, .strval = "hi" };
    print_value(&a);
    print_value(&b);
    return 0;
}
```
This is the exact shape cJSON uses for a JSON value (§36) and your own
Section 16.3 JSON parser's `JsonValue` — the same three-line pattern
(tag, anonymous union, switch) is how C represents "one of several
types" everywhere it comes up.
</details>

**4. Conceptual:** what actually goes wrong if you delete the
`#ifndef POINT_H` / `#define POINT_H` / `#endif` lines from a header
that two different `.c` files both `#include` (possibly indirectly,
through another header)?
<details><summary>Answer</summary>

The preprocessor is pure text substitution (§8.4) — without an include
guard, if the same header gets `#include`d twice into *one* translation
unit, every declaration in it gets pasted in twice. For a function
prototype that's harmless (redeclaring the same prototype is legal), but
for a `struct`/`typedef` definition it's a hard compile error:
*"redefinition of struct Point."* The include guard makes the second
`#include` of the same file expand to nothing, because `POINT_H` is
already defined the second time the preprocessor reaches the
`#ifndef` line.
</details>

---

## 9. Mini projects: a unit converter and a word-frequency counter

### 9.1 Mini project — CLI unit converter

Goal: apply everything above (parsing `argv`, `strtod`, `switch`,
`snprintf`) in one small, real tool.

```c
// unitconv.c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
    const char *from, *to;
    double factor;   // multiply "from" units by factor to get "to" units
} Conversion;

static const Conversion TABLE[] = {
    {"km", "mi", 0.621371},
    {"mi", "km", 1.60934},
    {"kg", "lb", 2.20462},
    {"lb", "kg", 0.453592},
    {"c",  "f",  0},  // handled specially below
    {"f",  "c",  0},
};
#define TABLE_LEN (sizeof(TABLE) / sizeof(TABLE[0]))

int main(int argc, char *argv[]) {
    if (argc != 4) {
        fprintf(stderr, "usage: %s <value> <from> <to>\n", argv[0]);
        return 1;
    }

    char *end;
    double value = strtod(argv[1], &end);
    if (*end != '\0') {
        fprintf(stderr, "error: '%s' is not a number\n", argv[1]);
        return 1;
    }
    const char *from = argv[2], *to = argv[3];

    if (strcmp(from, "c") == 0 && strcmp(to, "f") == 0) {
        printf("%.2f\n", value * 9.0 / 5.0 + 32.0);
        return 0;
    }
    if (strcmp(from, "f") == 0 && strcmp(to, "c") == 0) {
        printf("%.2f\n", (value - 32.0) * 5.0 / 9.0);
        return 0;
    }

    for (size_t i = 0; i < TABLE_LEN; i++) {
        if (strcmp(TABLE[i].from, from) == 0 && strcmp(TABLE[i].to, to) == 0) {
            printf("%.4f\n", value * TABLE[i].factor);
            return 0;
        }
    }

    fprintf(stderr, "error: no conversion from '%s' to '%s'\n", from, to);
    return 1;
}
```

```bash
clang -std=c17 -Wall -Wextra -Werror -o unitconv unitconv.c
./unitconv 10 km mi     # 6.2137
./unitconv 100 c f      # 212.00
```

### 9.2 Mini project — word-frequency counter

Goal: `fgets`, `strtok`, a fixed-size hash-free lookup via linear scan
(you'll replace this with a real hash table in Section 18), and dynamic
memory ownership done correctly for the first time.

```c
// wordfreq.c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>

typedef struct {
    char *word;
    int   count;
} Entry;

#define MAX_WORDS 10000

static Entry entries[MAX_WORDS];
static int   num_entries = 0;

static int record(const char *word) {
    for (int i = 0; i < num_entries; i++) {
        if (strcmp(entries[i].word, word) == 0) {
            entries[i].count++;
            return 1;
        }
    }
    if (num_entries >= MAX_WORDS) return 1;
    entries[num_entries].word = strdup(word);   // heap copy — word buffer gets reused/mutated
    if (!entries[num_entries].word) return 0;
    entries[num_entries].count = 1;
    num_entries++;
    return 1;
}

static int by_count_desc(const void *a, const void *b) {
    const Entry *ea = a, *eb = b;
    return (eb->count > ea->count) - (eb->count < ea->count);  // negative/zero/positive
}

int main(void) {
    char line[1024];
    while (fgets(line, sizeof line, stdin)) {
        for (char *p = line; *p; p++) *p = (char)tolower((unsigned char)*p);
        char *tok = strtok(line, " \t\n.,!?;:\"'()");
        while (tok) {
            if (!record(tok)) { perror("strdup"); return 1; }
            tok = strtok(NULL, " \t\n.,!?;:\"'()");
        }
    }

    qsort(entries, (size_t)num_entries, sizeof(Entry), by_count_desc);

    for (int i = 0; i < num_entries && i < 20; i++) {
        printf("%5d  %s\n", entries[i].count, entries[i].word);
    }

    for (int i = 0; i < num_entries; i++) free(entries[i].word);  // own what you strdup
    return 0;
}
```

```bash
clang -std=gnu17 -Wall -Wextra -Werror -o wordfreq wordfreq.c   # gnu17: wordfreq.c calls
                                                                  # strdup(), a POSIX function
                                                                  # strict c17 hides on Linux/glibc
curl -s https://www.gutenberg.org/files/1342/1342-0.txt | ./wordfreq | head
```

Note the ownership rule established here and used for the rest of the
guide: **whoever calls `strdup`/`malloc` is responsible for the matching
`free`.** `record()` allocates; `main()`'s final loop frees. This
convention — and what happens when you break it — is the entire subject
of Part II.

---

# Part II — Intermediate: memory, modules, and tooling

> Part I's programs lived comfortably inside one function's worth of stack
> memory. This Part is where C stops being safe by default: dynamic
> allocation and the ownership rules that go with it, splitting a program
> across multiple files, file I/O, and the debugging tools (sanitizers, a
> minimal test harness) you need once a bug is no longer obvious from
> reading the code. Ends with a vector library, a CSV parser, and a JSON
> parser — real tools, not toy examples.

## 10. Dynamic memory: the `malloc` family, ownership, and the bugs it enables

### 10.1 The four functions and what they actually do

```c
#include <stdlib.h>

void *malloc(size_t size);              // allocate size bytes, uninitialized
void *calloc(size_t n, size_t size);    // allocate n*size bytes, zeroed, overflow-checked
void *realloc(void *ptr, size_t size);  // resize, may move, contents preserved up to min(old,new)
void  free(void *ptr);                  // return memory to the allocator
```

```c
#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int *arr = malloc(10 * sizeof(int));
    if (!arr) { perror("malloc"); return 1; }   // ALWAYS check — allocation can fail

    for (int i = 0; i < 10; i++) arr[i] = i * i;

    int *bigger = realloc(arr, 20 * sizeof(int));
    if (!bigger) {           // realloc failing does NOT free the original
        free(arr);            // so you still must free the original on failure
        return 1;
    }
    arr = bigger;             // only reassign after confirming success

    for (int i = 10; i < 20; i++) arr[i] = i * i;
    for (int i = 0; i < 20; i++) printf("%d ", arr[i]);
    printf("\n");

    free(arr);
    return 0;
}
```

**Why `calloc(n, size)` over `malloc(n * size)`:** `calloc` checks for
multiplication overflow internally before allocating; `n * size` computed
by hand can silently overflow and allocate a tiny buffer that you then
write far past the end of. This exact bug class (integer overflow feeding
an allocation size) is CWE-190 and shows up in real CVEs regularly.

### 10.2 The four classic memory bugs, each demonstrated and fixed

```c
// 1. Memory leak — allocated, never freed
void leak(void) {
    char *buf = malloc(100);
    (void)buf;
}   // buf goes out of scope; the 100 bytes are now unreachable and unfreeable

// 2. Use-after-free — reading/writing memory you already returned
void use_after_free(void) {
    char *buf = malloc(100);
    free(buf);
    buf[0] = 'x';        // UB — the allocator may have already reused this memory
}

// 3. Double free — freeing the same pointer twice corrupts the allocator's internal metadata
void double_free(void) {
    char *buf = malloc(100);
    free(buf);
    free(buf);            // UB — often crashes, sometimes exploitable
}

// 4. Buffer overflow — writing past what was allocated
void overflow(void) {
    char *buf = malloc(10);
    buf[10] = 'x';        // UB — index 10 is one past the 10 valid bytes (0..9)
    free(buf);
}
```

Section 14 shows how AddressSanitizer catches every one of these
instantly, with a stack trace pointing at the exact allocation and the
exact bad access — this is why no serious C project ships without running
its test suite under ASan.

### 10.3 Ownership: the discipline that replaces a garbage collector

C has no automatic memory management, so every project needs an explicit
**ownership convention** — a rule for who is responsible for freeing what.
The three that cover almost all real code:

1. **Caller allocates, caller frees** (most common): a function fills a
   buffer the caller provided.
2. **Callee allocates, caller frees** (`strdup`, `malloc`-returning
   functions): document this in the function's name or a comment —
   `create_`, `_new`, `_dup` prefixes/suffixes are a strong convention.
3. **Ownership transfer**: a function takes a pointer and becomes
   responsible for freeing it later (e.g., adding a node to a linked list
   that owns its nodes).

```c
// Convention: functions named *_new return heap memory the caller must
// free with the matching *_free. This pairing is the backbone of every
// C library's public API — see Part VI's SQLite/Redis/curl walkthroughs.
typedef struct { char *data; size_t len; } Buffer;

Buffer *buffer_new(size_t capacity) {
    Buffer *b = malloc(sizeof(Buffer));
    if (!b) return NULL;
    b->data = malloc(capacity);
    if (!b->data) { free(b); return NULL; }
    b->len = 0;
    return b;
}

void buffer_free(Buffer *b) {
    if (!b) return;      // free(NULL) is always safe — mirror that here
    free(b->data);
    free(b);
}
```

**Real-world example.** This exact `_new`/`_free` pairing, with a NULL
check at the top of `_free` so `foo_free(NULL)` is always safe, is the
pattern you'll see in Section 33's SQLite walkthrough (`sqlite3_close`),
Section 34's Redis walkthrough (`sdsfree`), and Section 35's curl
walkthrough (`curl_easy_cleanup`).

---

## 11. Pointers to pointers, function pointers, pointers to arrays

### 11.1 Pointers to pointers: modifying a pointer through a function

```c
#include <stdio.h>
#include <stdlib.h>

int allocate(int **out, int value) {
    *out = malloc(sizeof(int));   // to change the CALLER's pointer, you need its address
    if (!*out) return 0;
    **out = value;
    return 1;
}

int main(void) {
    int *p = NULL;
    if (!allocate(&p, 42)) { perror("malloc"); return 1; }
    printf("%d\n", *p);
    free(p);
    return 0;
}
```

This is also exactly how `argv` works: `int main(int argc, char *argv[])`
— `argv` is `char **`, an array of pointers to `char` (strings), because
`main` needs a variable number of variable-length strings, which is
precisely what a pointer-to-pointer expresses.

### 11.2 Function pointers: C's mechanism for callbacks and polymorphism

```c
#include <stdio.h>
#include <stdlib.h>

int compare_asc(const void *a, const void *b) {
    int ia = *(const int *)a, ib = *(const int *)b;
    return (ia > ib) - (ia < ib);
}
int compare_desc(const void *a, const void *b) {
    int ia = *(const int *)a, ib = *(const int *)b;
    return (ib > ia) - (ib < ia);
}

int main(void) {
    int nums[] = {5, 2, 8, 1, 9};
    int n = sizeof(nums) / sizeof(nums[0]);

    int (*cmp)(const void *, const void *) = compare_asc;  // a variable holding a function
    qsort(nums, (size_t)n, sizeof(int), cmp);

    for (int i = 0; i < n; i++) printf("%d ", nums[i]);
    printf("\n");

    qsort(nums, (size_t)n, sizeof(int), compare_desc);
    for (int i = 0; i < n; i++) printf("%d ", nums[i]);
    printf("\n");
    return 0;
}
```

`qsort`'s comparator argument is the same mechanism used for every
callback-driven C API you'll meet in Part VI — curl's
`CURLOPT_WRITEFUNCTION`, Redis's command dispatch table, and every event
loop's "call this function when X happens."

A typedef makes function-pointer types readable:

```c
typedef int (*Comparator)(const void *, const void *);

void sort_ints(int *arr, size_t n, Comparator cmp) {
    qsort(arr, n, sizeof(int), cmp);
}
```

### 11.3 Pointers to arrays vs. arrays of pointers — the syntax that trips everyone up

```c
int  (*p1)[5];   // p1 is a pointer to an array of 5 ints
int  *p2[5];     // p2 is an array of 5 pointers to int

int matrix[3][4];
int (*row)[4] = matrix;      // row points to a whole 4-int row; row+1 skips 4 ints
```

Read declarations from the variable name outward, following precedence:
`[]` and `()` bind tighter than `*`, so `int *p2[5]` groups as
`int *(p2[5])` — an array, of pointers.

---

## 12. File I/O: buffering, binary vs. text mode, `errno`

### 12.1 `stdio` buffering — why your output sometimes "doesn't show up"

```c
#include <stdio.h>
#include <unistd.h>

int main(void) {
    printf("about to crash...");   // stdout is line-buffered (tty) or fully buffered (pipe/file) —
                                    // this may sit in a buffer, unflushed, if stdout isn't a tty
    fflush(stdout);                 // force it out now
    // *(int *)0 = 1;               // simulate a crash — without the fflush above,
                                     // the message might never appear
    return 0;
}
```

`stderr` is unbuffered by default specifically so error messages appear
immediately even if the process crashes right after — this is why
diagnostics go to `stderr`, not `stdout`.

### 12.2 Reading a whole file, correctly

```c
#include <stdio.h>
#include <stdlib.h>

char *read_whole_file(const char *path, size_t *out_len) {
    FILE *f = fopen(path, "rb");            // "b" matters: on Windows, text mode translates
    if (!f) return NULL;                    // \r\n <-> \n, corrupting binary data

    fseek(f, 0, SEEK_END);
    long size = ftell(f);
    if (size < 0) { fclose(f); return NULL; }
    rewind(f);

    char *buf = malloc((size_t)size + 1);
    if (!buf) { fclose(f); return NULL; }

    size_t read = fread(buf, 1, (size_t)size, f);
    fclose(f);

    if (read != (size_t)size) { free(buf); return NULL; }
    buf[size] = '\0';
    if (out_len) *out_len = read;
    return buf;    // caller owns this — must free()
}
```

### 12.3 `errno`: check it right after the call that can fail, nothing else

```c
#include <errno.h>
#include <stdio.h>
#include <string.h>

int main(void) {
    FILE *f = fopen("/does/not/exist", "r");
    if (!f) {
        // errno is set by the failing call; a call in between (even printf) can clobber it
        fprintf(stderr, "fopen failed: %s\n", strerror(errno));
        return 1;
    }
    fclose(f);
    return 0;
}
```

`perror("fopen")` is the shorthand for the same thing:
`perror` prints `"fopen: No such file or directory\n"` using the current
`errno` automatically.

---

## 13. Multi-file programs: headers, `static`, translation units, and `make`

### 13.1 `static` at file scope: C's only real access-control keyword

```c
// counter.c
#include "counter.h"

static int count = 0;              // static at file scope = "private to this .c file";
                                    // no other translation unit can link against `count`
static void log_change(void) {     // static function = private helper, invisible outside this file
    // ...
}

int counter_increment(void) {
    count++;
    log_change();
    return count;
}
```

```c
// counter.h
#ifndef COUNTER_H
#define COUNTER_H
int counter_increment(void);
#endif
```

This `static`-for-privacy convention is exactly how C fakes
"encapsulation": the header exposes only what callers need
(`counter_increment`), everything else is invisible outside the `.c` file.

### 13.2 A real multi-file `make` setup

```
project/
├── Makefile
├── include/
│   └── counter.h
└── src/
    ├── counter.c
    └── main.c
```

```makefile
CC      = clang
CFLAGS  = -std=c17 -Wall -Wextra -Wpedantic -Werror -g -Iinclude
SRCS    = $(wildcard src/*.c)
OBJS    = $(SRCS:src/%.c=build/%.o)
BIN     = build/app

all: $(BIN)

$(BIN): $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^

build/%.o: src/%.c | build
	$(CC) $(CFLAGS) -c $< -o $@

build:
	mkdir -p build

clean:
	rm -rf build

.PHONY: all clean
```

Each `.c` file compiles to its own `.o` **independently** (that's a
"translation unit") — the linker then combines them, resolving calls
between them by symbol name. This is why one file's compile error doesn't
require recompiling everything else, and why `static` symbols in one file
are invisible to the linker when resolving another.

---

## 14. Debugging and sanitizers: gdb/lldb, Valgrind, ASan/UBSan

### 14.1 AddressSanitizer — catch memory bugs the instant they happen

```bash
clang -std=c17 -g -fsanitize=address -o buggy buggy.c
./buggy
```

```c
// buggy.c
#include <stdlib.h>
int main(void) {
    int *arr = malloc(5 * sizeof(int));
    arr[5] = 1;   // one past the end
    free(arr);
    return 0;
}
```

Running this under ASan prints a precise report: `heap-buffer-overflow`,
the exact line of the bad write, and the exact line of the `malloc` that
created the buffer — versus a plain build, which might not crash at all
(the overflow lands in unused allocator padding) and ship the bug to
production. **Every project in this guide should be tested under
`-fsanitize=address,undefined` before you trust it.**

### 14.2 UndefinedBehaviorSanitizer — catch UB, not just memory bugs

```bash
clang -std=c17 -g -fsanitize=undefined -o ubtest ubtest.c
```

```c
#include <limits.h>
int main(void) {
    int x = INT_MAX;
    return x + 1;   // UBSan: "signed integer overflow: 2147483647 + 1 cannot be
                     //         represented in type 'int'"
}
```

### 14.3 Valgrind — the other lens (great for leak detection on Linux)

```bash
valgrind --leak-check=full --show-leak-kinds=all ./buggy
```

Valgrind runs your binary on a software CPU and tracks every allocation;
it's slower than ASan but doesn't require recompiling, and its leak
reports (`definitely lost`, `still reachable`) are the standard vocabulary
for memory-leak triage in C/C++ shops.

### 14.4 A five-minute gdb/lldb session

```bash
clang -g -O0 -o buggy buggy.c
lldb ./buggy     # or: gdb ./buggy
(lldb) run
(lldb) bt                 # backtrace — where did it crash?
(lldb) frame select 1     # move up the stack
(lldb) print arr[3]       # inspect a value
(lldb) break main.c:12    # set a breakpoint
(lldb) continue
```

**Real-world example.** The standard triage loop in production C: a core
dump or crash report comes in → load it in gdb/lldb (`gdb ./app core`) →
`bt` for the backtrace → walk frames with `frame select N` → inspect
variables. This is the exact workflow used to debug the segfaults you'll
intentionally cause in Section 10.2's exercises, and the real ones you'll
cause by accident for the rest of your C career.

---

## 15. A minimal test harness for C

C has no built-in test framework; most real projects either pull in
something like Unity/Check/CMocka, or — for smaller libraries — write 30
lines of macros once and reuse them forever. Here's that 30 lines,
because understanding it demystifies every "real" framework:

```c
// test.h
#ifndef TEST_H
#define TEST_H
#include <stdio.h>

static int tests_run = 0, tests_failed = 0;

#define ASSERT_EQ(actual, expected) do {                                    \
    tests_run++;                                                            \
    if ((actual) != (expected)) {                                           \
        tests_failed++;                                                     \
        fprintf(stderr, "FAIL %s:%d: expected %ld, got %ld\n",              \
                __FILE__, __LINE__, (long)(expected), (long)(actual));      \
    }                                                                        \
} while (0)

#define TEST_SUMMARY() do {                                                  \
    printf("%d/%d tests passed\n", tests_run - tests_failed, tests_run);     \
    return tests_failed ? 1 : 0;                                             \
} while (0)

#endif
```

```c
// test_add.c
#include "test.h"

int add(int a, int b) { return a + b; }

int main(void) {
    ASSERT_EQ(add(2, 2), 4);
    ASSERT_EQ(add(-1, 1), 0);
    ASSERT_EQ(add(0, 0), 1);   // intentionally wrong, to see a FAIL line
    TEST_SUMMARY();
}
```

`do { ... } while (0)` wrapping a multi-statement macro is the standard
trick that makes `ASSERT_EQ(x, y);` behave like a single statement inside
an `if` without a trailing-semicolon or dangling-else bug — you'll see
this exact pattern in nearly every C macro that expands to more than one
expression.

---

## 16. Mini projects: a dynamic-array library + CSV parser, and a JSON parser

### 16.1 Mini project — a generic dynamic array (vector)

This is the pattern every C container library uses (Section 37's uthash,
Section 19's intrusive lists, and any "`vec.h`"-style single-header
library on GitHub): a macro that generates a typed struct + functions,
because C has no real generics before you reach for `void*` + casts or
`_Generic`.

```c
// vec.h — a growable array, macro-generated per type
#ifndef VEC_H
#define VEC_H
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#define VEC_DEFINE(T, Name)                                                  \
typedef struct { T *data; size_t len, cap; } Name;                           \
                                                                               \
static inline void Name##_init(Name *v) { v->data = NULL; v->len = v->cap = 0; } \
                                                                               \
static inline int Name##_push(Name *v, T item) {                             \
    if (v->len == v->cap) {                                                  \
        if (v->cap > SIZE_MAX / 2) return 0;                                  \
        size_t new_cap = v->cap ? v->cap * 2 : 4;                            \
        if (new_cap > SIZE_MAX / sizeof(T)) return 0;                         \
        T *new_data = realloc(v->data, new_cap * sizeof(T));                 \
        if (!new_data) return 0;                                             \
        v->data = new_data;                                                  \
        v->cap = new_cap;                                                    \
    }                                                                        \
    v->data[v->len++] = item;                                                \
    return 1;                                                                \
}                                                                             \
	                                                                               \
static inline void Name##_free(Name *v) { free(v->data); v->data = NULL; v->len = v->cap = 0; }

#endif
```

```c
// main.c
#include <stdio.h>
#include "vec.h"

VEC_DEFINE(int, IntVec)

int main(void) {
    IntVec v;
    IntVec_init(&v);
    for (int i = 0; i < 10; i++) {
        if (!IntVec_push(&v, i * i)) { perror("IntVec_push"); IntVec_free(&v); return 1; }
    }
    for (size_t i = 0; i < v.len; i++) printf("%d ", v.data[i]);
    printf("\n");
    IntVec_free(&v);
    return 0;
}
```

Doubling the capacity on growth (`cap * 2`) gives **amortized O(1)**
`push` — the classic result every dynamic-array implementation (C++
`std::vector`, Go slices, Python lists) relies on: most pushes are O(1),
occasional resizes are O(n), and the total cost over n pushes is O(n).

### 16.2 Mini project — a CSV parser using the vector above

```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "vec.h"

VEC_DEFINE(char *, StrVec)

static int split_line(const char *line, char delim, StrVec *out) {
    const char *start = line;
    for (const char *p = line; ; p++) {
        if (*p == delim || *p == '\0' || *p == '\n') {
            size_t len = (size_t)(p - start);
            char *field = malloc(len + 1);
            if (!field) return 0;
            memcpy(field, start, len);
            field[len] = '\0';
            if (!StrVec_push(out, field)) { free(field); return 0; }
            start = p + 1;
            if (*p == '\0' || *p == '\n') break;
        }
    }
    return 1;
}

int main(void) {
    char line[1024];
    while (fgets(line, sizeof line, stdin)) {
        StrVec fields;
        StrVec_init(&fields);
        if (!split_line(line, ',', &fields)) {
            fprintf(stderr, "out of memory\n");
            for (size_t i = 0; i < fields.len; i++) free(fields.data[i]);
            StrVec_free(&fields);
            return 1;
        }

        for (size_t i = 0; i < fields.len; i++) {
            printf("[%zu]=%s ", i, fields.data[i]);
            free(fields.data[i]);   // each field was malloc'd in split_line
        }
        printf("\n");
        StrVec_free(&fields);
    }
    return 0;
}
```

This is a **teaching** CSV parser — it doesn't handle quoted fields with
embedded commas. Real CSV parsing needs a small state machine; that's
exactly the technique the JSON parser below and Section 36's cJSON
walkthrough use for a harder grammar.

### 16.3 Mini project — a minimal JSON parser (recursive descent)

This is the single most valuable exercise in the whole guide for
"reading real C" — Section 36 walks through cJSON's actual source right
after this, and you'll recognize every function.

```c
// json.h
#ifndef JSON_H
#define JSON_H
#include <stddef.h>   // size_t

typedef enum { JSON_NULL, JSON_BOOL, JSON_NUMBER, JSON_STRING, JSON_ARRAY, JSON_OBJECT } JsonType;

typedef struct JsonValue JsonValue;
typedef struct { char *key; JsonValue *value; } JsonMember;

struct JsonValue {
    JsonType type;
    union {
        int    boolean;
        double number;
        char  *string;
        struct { JsonValue **items; size_t count; } array;
        struct { JsonMember *members; size_t count; } object;
    };
};

JsonValue *json_parse(const char *text);
void       json_free(JsonValue *v);

#endif
```

```c
// json.c — a compact recursive-descent parser
#include "json.h"
#include <stdlib.h>
#include <string.h>
#include <ctype.h>

typedef struct { const char *p; } Parser;

static void skip_ws(Parser *ps) { while (isspace((unsigned char)*ps->p)) ps->p++; }

static JsonValue *parse_value(Parser *ps);

static JsonValue *new_value(JsonType t) {
    JsonValue *v = calloc(1, sizeof(JsonValue));
    if (!v) return NULL;
    v->type = t;
    return v;
}

static char *parse_raw_string(Parser *ps) {
    ps->p++; // skip opening quote
    const char *start = ps->p;
    while (*ps->p && *ps->p != '"') ps->p++;   // no escape handling — kept minimal on purpose
    size_t len = (size_t)(ps->p - start);
    char *s = malloc(len + 1);
    if (!s) return NULL;
    memcpy(s, start, len);
    s[len] = '\0';
    if (*ps->p == '"') ps->p++;
    return s;
}

static JsonValue *parse_string(Parser *ps) {
    JsonValue *v = new_value(JSON_STRING);
    if (!v) return NULL;
    v->string = parse_raw_string(ps);
    if (!v->string) { free(v); return NULL; }
    return v;
}

static JsonValue *parse_number(Parser *ps) {
    char *end;
    double n = strtod(ps->p, &end);
    ps->p = end;
    JsonValue *v = new_value(JSON_NUMBER);
    if (!v) return NULL;
    v->number = n;
    return v;
}

static JsonValue *parse_array(Parser *ps) {
    JsonValue *v = new_value(JSON_ARRAY);
    if (!v) return NULL;
    ps->p++; skip_ws(ps);   // skip '['
    if (*ps->p == ']') { ps->p++; return v; }

    size_t cap = 4;
    v->array.items = malloc(cap * sizeof(JsonValue *));
    if (!v->array.items) { free(v); return NULL; }
    while (1) {
        skip_ws(ps);
        if (v->array.count == cap) {
            cap *= 2;
            JsonValue **new_items = realloc(v->array.items, cap * sizeof(JsonValue *));
            if (!new_items) { json_free(v); return NULL; }
            v->array.items = new_items;
        }
        JsonValue *item = parse_value(ps);
        if (!item) { json_free(v); return NULL; }
        v->array.items[v->array.count++] = item;
        skip_ws(ps);
        if (*ps->p == ',') { ps->p++; continue; }
        break;
    }
    skip_ws(ps);
    if (*ps->p == ']') ps->p++;
    return v;
}

static JsonValue *parse_object(Parser *ps) {
    JsonValue *v = new_value(JSON_OBJECT);
    if (!v) return NULL;
    ps->p++; skip_ws(ps);   // skip '{'
    if (*ps->p == '}') { ps->p++; return v; }

    size_t cap = 4;
    v->object.members = malloc(cap * sizeof(JsonMember));
    if (!v->object.members) { free(v); return NULL; }
    while (1) {
        skip_ws(ps);
        char *key = parse_raw_string(ps);
        if (!key) { json_free(v); return NULL; }
        skip_ws(ps);
        if (*ps->p == ':') ps->p++;
        skip_ws(ps);
        JsonValue *val = parse_value(ps);
        if (!val) { free(key); json_free(v); return NULL; }

        if (v->object.count == cap) {
            cap *= 2;
            JsonMember *new_members = realloc(v->object.members, cap * sizeof(JsonMember));
            if (!new_members) { free(key); json_free(val); json_free(v); return NULL; }
            v->object.members = new_members;
        }
        v->object.members[v->object.count].key = key;
        v->object.members[v->object.count].value = val;
        v->object.count++;

        skip_ws(ps);
        if (*ps->p == ',') { ps->p++; continue; }
        break;
    }
    skip_ws(ps);
    if (*ps->p == '}') ps->p++;
    return v;
}

static JsonValue *parse_value(Parser *ps) {
    skip_ws(ps);
    switch (*ps->p) {
        case '"': return parse_string(ps);
        case '[': return parse_array(ps);
        case '{': return parse_object(ps);
        case 't': ps->p += 4; { JsonValue *v = new_value(JSON_BOOL); if (v) v->boolean = 1; return v; }
        case 'f': ps->p += 5; { JsonValue *v = new_value(JSON_BOOL); if (v) v->boolean = 0; return v; }
        case 'n': ps->p += 4; return new_value(JSON_NULL);
        default:  return parse_number(ps);
    }
}

JsonValue *json_parse(const char *text) {
    Parser ps = { .p = text };
    return parse_value(&ps);
}

void json_free(JsonValue *v) {
    if (!v) return;
    switch (v->type) {
        case JSON_STRING:
            free(v->string);
            break;
        case JSON_ARRAY:
            for (size_t i = 0; i < v->array.count; i++) json_free(v->array.items[i]);
            free(v->array.items);
            break;
        case JSON_OBJECT:
            for (size_t i = 0; i < v->object.count; i++) {
                free(v->object.members[i].key);
                json_free(v->object.members[i].value);
            }
            free(v->object.members);
            break;
        default: break;
    }
    free(v);
}
```

Recursive descent works because JSON's grammar is itself recursive
(`value` can contain `array`s and `object`s which contain more `value`s)
— each grammar rule becomes one function, and each function calls the
ones for the rules nested inside it. **This is not a toy technique**:
it's the same broad structure used by many hand-written parsers,
including cJSON's (Section 36) and SQLite's SQL tokenizer/parser
(Section 33). This compact example is intentionally incomplete: it omits
escape sequences, exponent notation, literal validation, depth limits, and
full syntax errors so the recursive-descent shape stays visible.

---

# Part III — Advanced: data structures, concurrency, networking

> You can allocate memory and organize a program across files; this Part is
> where you build the structures that actually hold data — linked lists,
> trees, hash tables, and the Linux-kernel's intrusive-container trick —
> then add threads, atomics, and raw sockets on top. It ends with a thread
> pool and a multi-client chat server that exercises concurrency and
> networking together, not in isolation.

## 17. Bit manipulation and bit-fields

```c
#include <stdio.h>
#include <stdint.h>

#define FLAG_READ    (1u << 0)
#define FLAG_WRITE   (1u << 1)
#define FLAG_EXECUTE (1u << 2)

int main(void) {
    uint32_t perms = FLAG_READ | FLAG_WRITE;   // set two flags

    if (perms & FLAG_WRITE) printf("writable\n");   // test a flag
    perms &= ~FLAG_WRITE;                            // clear a flag
    perms ^= FLAG_EXECUTE;                           // toggle a flag

    printf("%u\n", perms);

    // Real-world example: this is exactly how chmod-style permission bits,
    // network protocol flag fields, and hardware register access all work.
    return 0;
}
```

Bit-fields let a struct pack sub-byte fields, common in wire-format and
hardware-register structs (layout is implementation-defined in *order*,
which is why real protocol code prefers explicit shifting/masking over
bit-fields for anything that crosses machine boundaries):

```c
typedef struct {
    unsigned int is_admin : 1;
    unsigned int level    : 4;   // 0-15
    unsigned int reserved : 3;
} PermissionByte;   // fits in 1 byte total
```

---

## 18. Data structures from scratch: linked lists, stacks/queues, hash tables, trees

### 18.1 Singly linked list

```c
#include <stdio.h>
#include <stdlib.h>

typedef struct Node {
    int value;
    struct Node *next;
} Node;

Node *node_new(int value) {
    Node *n = malloc(sizeof(Node));
    if (!n) return NULL;
    n->value = value;
    n->next = NULL;
    return n;
}

int list_push_front(Node **head, int value) {
    Node *n = node_new(value);
    if (!n) return 0;
    n->next = *head;
    *head = n;             // pointer-to-pointer: rewrites the caller's head
    return 1;
}

void list_free(Node *head) {
    while (head) {
        Node *next = head->next;
        free(head);
        head = next;
    }
}

int main(void) {
    Node *head = NULL;
    for (int i = 1; i <= 5; i++) {
        if (!list_push_front(&head, i)) { perror("malloc"); list_free(head); return 1; }
    }
    for (Node *n = head; n; n = n->next) printf("%d ", n->value);
    printf("\n");
    list_free(head);
    return 0;
}
```

### 18.2 A real hash table (separate chaining, `djb2` hash)

```c
// hashtable.h/.c — string-keyed hash table with chaining
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct HTEntry {
    char *key;
    int   value;
    struct HTEntry *next;   // chain for collisions
} HTEntry;

typedef struct {
    HTEntry **buckets;
    size_t    num_buckets;
    size_t    count;
} HashTable;

static unsigned long djb2(const char *s) {
    unsigned long hash = 5381;
    int c;
    while ((c = (unsigned char)*s++)) hash = ((hash << 5) + hash) + (unsigned long)c; // hash*33 + c
    return hash;
}

static char *xstrdup(const char *s) {
    size_t len = strlen(s) + 1;
    char *copy = malloc(len);
    if (!copy) return NULL;
    memcpy(copy, s, len);
    return copy;
}

HashTable *ht_new(size_t num_buckets) {
    if (num_buckets == 0) return NULL;
    HashTable *ht = malloc(sizeof(HashTable));
    if (!ht) return NULL;
    ht->buckets = calloc(num_buckets, sizeof(HTEntry *));
    if (!ht->buckets) { free(ht); return NULL; }
    ht->num_buckets = num_buckets;
    ht->count = 0;
    return ht;
}

int ht_set(HashTable *ht, const char *key, int value) {
    unsigned long idx = djb2(key) % ht->num_buckets;
    for (HTEntry *e = ht->buckets[idx]; e; e = e->next) {
        if (strcmp(e->key, key) == 0) { e->value = value; return 1; }  // update existing
    }
    HTEntry *e = malloc(sizeof(HTEntry));
    if (!e) return 0;
    e->key = xstrdup(key);
    if (!e->key) { free(e); return 0; }
    e->value = value;
    e->next = ht->buckets[idx];   // insert at head of chain
    ht->buckets[idx] = e;
    ht->count++;
    return 1;
}

int ht_get(HashTable *ht, const char *key, int *out) {
    unsigned long idx = djb2(key) % ht->num_buckets;
    for (HTEntry *e = ht->buckets[idx]; e; e = e->next) {
        if (strcmp(e->key, key) == 0) { *out = e->value; return 1; }
    }
    return 0;
}

void ht_free(HashTable *ht) {
    if (!ht) return;
    for (size_t i = 0; i < ht->num_buckets; i++) {
        HTEntry *e = ht->buckets[i];
        while (e) { HTEntry *next = e->next; free(e->key); free(e); e = next; }
    }
    free(ht->buckets);
    free(ht);
}

int main(void) {
    HashTable *ht = ht_new(16);
    if (!ht) { perror("ht_new"); return 1; }
    if (!ht_set(ht, "apples", 5) ||
        !ht_set(ht, "oranges", 3) ||
        !ht_set(ht, "apples", 7)) {   // update
        perror("ht_set");
        ht_free(ht);
        return 1;
    }

    int v;
    if (ht_get(ht, "apples", &v)) printf("apples = %d\n", v);   // 7
    if (!ht_get(ht, "bananas", &v)) printf("bananas not found\n");

    ht_free(ht);
    return 0;
}
```

`djb2` (Dan Bernstein's hash) is real, used-in-production code — a
version of it or `FNV-1a` appears in countless small C libraries because
it's a handful of lines and good enough for non-adversarial keys.
**Real-world example.** For untrusted keys (e.g., HTTP header names from
the network), a predictable hash lets an attacker force worst-case O(n)
chains — a real DoS class called "hash flooding" (behind CVE-2011-4815
and similar). Production hash tables (Python's `dict`, Redis's `dict`)
seed the hash randomly per-process specifically to prevent this.

### 18.3 A binary search tree, for completeness

```c
typedef struct BSTNode {
    int value;
    struct BSTNode *left, *right;
} BSTNode;

BSTNode *bst_insert(BSTNode *root, int value) {
    if (!root) {
        BSTNode *n = malloc(sizeof(BSTNode));
        if (!n) return NULL;
        n->value = value; n->left = n->right = NULL;
        return n;
    }
    if (value < root->value) root->left  = bst_insert(root->left, value);
    else                     root->right = bst_insert(root->right, value);
    return root;
}

void bst_inorder(BSTNode *root, void (*visit)(int)) {
    if (!root) return;
    bst_inorder(root->left, visit);
    visit(root->value);
    bst_inorder(root->right, visit);
}
```

---

## 19. Intrusive containers: the Linux kernel's `container_of` trick

Every data structure so far *contains* its payload (`Node.value`).
Intrusive containers flip this: the **payload struct embeds the list
node**, and the container macro computes back from the embedded node's
address to the containing struct's address. This is how the Linux kernel
implements *one* generic doubly-linked-list type
(`struct list_head`) reused for thousands of unrelated structs, with zero
`void*` casts and zero per-type boilerplate.

```c
#include <stdio.h>
#include <stddef.h>   // offsetof

// The kernel's actual technique, simplified:
#define container_of(ptr, type, member) \
    ((type *)((char *)(ptr) - offsetof(type, member)))

typedef struct list_head {
    struct list_head *next, *prev;
} list_head;

static void list_init(list_head *h) { h->next = h->prev = h; }

static void list_add(list_head *new_node, list_head *head) {
    new_node->next = head->next;
    new_node->prev = head;
    head->next->prev = new_node;
    head->next = new_node;
}

// A real payload struct that EMBEDS the generic node:
typedef struct {
    int         pid;
    char        name[32];
    list_head   node;      // <-- the intrusive part
} Task;

int main(void) {
    list_head tasks;
    list_init(&tasks);

    Task t1 = { .pid = 101, .name = "init" };
    Task t2 = { .pid = 205, .name = "sshd" };
    list_add(&t1.node, &tasks);
    list_add(&t2.node, &tasks);

    for (list_head *p = tasks.next; p != &tasks; p = p->next) {
        Task *t = container_of(p, Task, node);   // walk the generic list, recover the real struct
        printf("pid=%d name=%s\n", t->pid, t->name);
    }
    return 0;
}
```

`container_of` works because `offsetof(Task, node)` is a compile-time
constant — the byte distance from the start of `Task` to its `node`
field — so subtracting that many bytes from `node`'s address always
recovers the enclosing `Task*`, regardless of what type embeds `node`.
This is literally `include/linux/list.h` in the Linux kernel source, used
for the process list, the list of open files, wait queues, and hundreds
of other kernel structures — one list implementation, zero per-type
code duplication, and (unlike the `void*`-based generic containers in
Section 16) **no extra allocation or indirection**: the node lives inside
the struct that needs it.

---

## 20. Concurrency: POSIX threads, mutexes/condvars, C11 `<threads.h>` and atomics

### 20.1 Threads and the race condition they immediately create

```c
#include <pthread.h>
#include <stdio.h>

static int counter = 0;   // shared, unprotected — this is the bug

void *increment(void *arg) {
    (void)arg;
    for (int i = 0; i < 100000; i++) counter++;   // read-modify-write, not atomic
    return NULL;
}

int main(void) {
    pthread_t t1, t2;
    pthread_create(&t1, NULL, increment, NULL);
    pthread_create(&t2, NULL, increment, NULL);
    pthread_join(t1, NULL);
    pthread_join(t2, NULL);

    printf("%d\n", counter);   // almost never 200000 — a classic data race
    return 0;
}
```

Compile with `-fsanitize=thread` to have ThreadSanitizer point at the
exact racing lines instead of guessing from a wrong number.

### 20.2 The fix: a mutex

```c
#include <pthread.h>
#include <stdio.h>

static int counter = 0;
static pthread_mutex_t lock = PTHREAD_MUTEX_INITIALIZER;

void *increment(void *arg) {
    (void)arg;
    for (int i = 0; i < 100000; i++) {
        pthread_mutex_lock(&lock);
        counter++;              // now protected — only one thread inside at a time
        pthread_mutex_unlock(&lock);
    }
    return NULL;
}

int main(void) {
    pthread_t t1, t2;
    pthread_create(&t1, NULL, increment, NULL);
    pthread_create(&t2, NULL, increment, NULL);
    pthread_join(t1, NULL);
    pthread_join(t2, NULL);
    printf("%d\n", counter);   // always 200000
    return 0;
}
```

### 20.3 Producer/consumer with a condition variable

```c
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>

#define CAP 5
static int buf[CAP];
static int head = 0, tail = 0, count = 0;
static pthread_mutex_t mtx = PTHREAD_MUTEX_INITIALIZER;
static pthread_cond_t  not_full  = PTHREAD_COND_INITIALIZER;
static pthread_cond_t  not_empty = PTHREAD_COND_INITIALIZER;

void produce(int value) {
    pthread_mutex_lock(&mtx);
    while (count == CAP) pthread_cond_wait(&not_full, &mtx);  // sleep until there's room
    buf[tail] = value; tail = (tail + 1) % CAP; count++;
    pthread_cond_signal(&not_empty);                          // wake a waiting consumer
    pthread_mutex_unlock(&mtx);
}

int consume(void) {
    pthread_mutex_lock(&mtx);
    while (count == 0) pthread_cond_wait(&not_empty, &mtx);   // sleep until there's data
    int value = buf[head]; head = (head + 1) % CAP; count--;
    pthread_cond_signal(&not_full);
    pthread_mutex_unlock(&mtx);
    return value;
}
```

`pthread_cond_wait` atomically unlocks the mutex and sleeps, then
re-locks it before returning — this is *why* it takes the mutex as an
argument, and it's the exact primitive behind every blocking queue in
every language runtime (Java's `wait`/`notify`, Python's
`threading.Condition`).

### 20.4 C11 atomics: lock-free counters

```c
#include <stdatomic.h>
#include <pthread.h>
#include <stdio.h>

static atomic_int counter = 0;

void *increment(void *arg) {
    (void)arg;
    for (int i = 0; i < 100000; i++) {
        atomic_fetch_add(&counter, 1);   // one indivisible hardware instruction, no mutex
    }
    return NULL;
}

int main(void) {
    pthread_t t1, t2;
    pthread_create(&t1, NULL, increment, NULL);
    pthread_create(&t2, NULL, increment, NULL);
    pthread_join(t1, NULL);
    pthread_join(t2, NULL);
    printf("%d\n", atomic_load(&counter));   // always 200000, no lock needed
    return 0;
}
```

Atomics are faster than a mutex for simple counters (no syscall, no
thread blocking) but they only protect *one* memory location at a time —
the moment an invariant spans two variables, you need a real lock.
Section 38 covers the memory-ordering semantics (`memory_order_relaxed`
vs. `acquire`/`release` vs. `seq_cst`) that `atomic_fetch_add`'s default
sequential consistency hides from you here.

---

## 21. Sockets: a TCP echo server, then a tiny HTTP server

### 21.1 A TCP echo server (POSIX sockets, the API every language's networking library wraps)

```c
// echo_server.c
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>

int main(void) {
    int server_fd = socket(AF_INET, SOCK_STREAM, 0);
    int opt = 1;
    setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

    struct sockaddr_in addr = {0};
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(9999);          // host-to-network byte order — always for wire formats

    bind(server_fd, (struct sockaddr *)&addr, sizeof(addr));
    listen(server_fd, 10);
    printf("listening on :9999\n");

    while (1) {
        struct sockaddr_in client_addr;
        socklen_t len = sizeof(client_addr);
        int client_fd = accept(server_fd, (struct sockaddr *)&client_addr, &len);
        if (client_fd < 0) continue;

        char buf[1024];
        ssize_t n;
        while ((n = read(client_fd, buf, sizeof buf)) > 0) {
            write(client_fd, buf, (size_t)n);   // echo it straight back
        }
        close(client_fd);
    }
    return 0;
}
```

```bash
clang -std=gnu17 -Wall -Wextra -o echo_server echo_server.c   # gnu17: BSD sockets are a POSIX
                                                                 # API, hidden under strict c17
./echo_server &
printf "hello\n" | nc localhost 9999
```

This one function call, `accept()`, is the reason the server above can
only handle one client at a time — while it's blocked reading from
`client_fd`, no other client can connect. Section 23's mini-project fixes
exactly this with `select()`, and Section 34's Redis walkthrough shows
the production answer (a single-threaded event loop).

### 21.2 A minimal HTTP/1.0 server on the same socket API

```c
// http_server.c — serves one static string to every request
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>

int main(void) {
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    int opt = 1;
    setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof opt);

    struct sockaddr_in addr = { .sin_family = AF_INET, .sin_port = htons(8080),
                                 .sin_addr.s_addr = INADDR_ANY };
    bind(fd, (struct sockaddr *)&addr, sizeof addr);
    listen(fd, 16);
    printf("http://localhost:8080\n");

    while (1) {
        int client = accept(fd, NULL, NULL);
        char req[2048] = {0};
        read(client, req, sizeof req - 1);   // real servers must parse Content-Length, not assume one read is the whole request

        const char *body = "hello from a C http server\n";
        char response[256];
        int len = snprintf(response, sizeof response,
            "HTTP/1.0 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %zu\r\n\r\n%s",
            strlen(body), body);
        write(client, response, (size_t)len);
        close(client);
    }
}
```

```bash
clang -std=gnu17 -Wall -Wextra -o http_server http_server.c
./http_server & curl -s http://localhost:8080
```

This is deliberately the simplest possible HTTP server to show the socket
API mapping directly onto the protocol you already know from a browser —
it's also deliberately *not production-ready* (single request at a time,
no request-size limits, no timeout handling); Section 34's Redis
walkthrough and Section 43's capstone show what closes those gaps.

---

## 22. Variadic functions and `_Generic`

### 22.1 Variadic functions (`printf`'s mechanism)

```c
#include <stdio.h>
#include <stdarg.h>

double sum(int count, ...) {
    va_list args;
    va_start(args, count);       // must start right after the last named parameter
    double total = 0;
    for (int i = 0; i < count; i++) {
        total += va_arg(args, double);  // caller and callee MUST agree on the type here —
                                         // there is no runtime check
    }
    va_end(args);
    return total;
}

int main(void) {
    printf("%.1f\n", sum(3, 1.0, 2.0, 3.0));  // 6.0
    return 0;
}
```

Because the compiler cannot check that the arguments you pass match the
types `va_arg` requests, variadic functions are the other classic UB
source in C after strings — this is exactly the mechanism behind Section
1.4's format-string mismatch bug.

### 22.2 `_Generic`: compile-time dispatch by type (C11)

```c
#include <stdio.h>

#define print_value(x) _Generic((x),           \
    int: print_int,                             \
    double: print_double,                       \
    char *: print_string                        \
)(x)

void print_int(int x)       { printf("int: %d\n", x); }
void print_double(double x) { printf("double: %f\n", x); }
void print_string(char *x)  { printf("string: %s\n", x); }

int main(void) {
    print_value(42);        // "int: 42" — resolved at COMPILE time, zero runtime cost
    print_value(3.14);      // "double: 3.140000"
    print_value("hi");      // "string: hi"
    return 0;
}
```

`_Generic` is how `<tgmath.h>`'s `sqrt`/`sin`/`cos` pick the `float` vs.
`double` vs. `long double` variant, and it's the closest C gets to
function overloading — but it's resolved entirely at compile time by the
*static* type of the expression, never at runtime, which is why it can't
replace a real `void*`-based runtime-polymorphic API like `qsort`'s.

---

## 23. Mini projects: a thread pool, and a `select()`-based chat server

### 23.1 Mini project — a thread pool

The pattern behind every real worker-pool (web server request handlers,
Section 26's shell running background jobs, Redis's background I/O
threads): a fixed number of worker threads pull tasks off a shared queue
protected by a mutex + condition variable.

```c
// threadpool.c
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>

#define NUM_THREADS 4
#define QUEUE_CAP   64

typedef void (*Task)(void *);

typedef struct { Task fn; void *arg; } Job;

static Job      queue[QUEUE_CAP];
static int      q_head = 0, q_tail = 0, q_count = 0;
static pthread_mutex_t q_lock = PTHREAD_MUTEX_INITIALIZER;
static pthread_cond_t  q_not_empty = PTHREAD_COND_INITIALIZER;
static int      shutdown_flag = 0;

void pool_submit(Task fn, void *arg) {
    pthread_mutex_lock(&q_lock);
    queue[q_tail] = (Job){ fn, arg };
    q_tail = (q_tail + 1) % QUEUE_CAP;
    q_count++;
    pthread_cond_signal(&q_not_empty);
    pthread_mutex_unlock(&q_lock);
}

void *worker(void *arg) {
    int id = *(int *)arg;
    while (1) {
        pthread_mutex_lock(&q_lock);
        while (q_count == 0 && !shutdown_flag) pthread_cond_wait(&q_not_empty, &q_lock);
        if (q_count == 0 && shutdown_flag) { pthread_mutex_unlock(&q_lock); break; }

        Job job = queue[q_head];
        q_head = (q_head + 1) % QUEUE_CAP;
        q_count--;
        pthread_mutex_unlock(&q_lock);

        job.fn(job.arg);   // run the task OUTSIDE the lock — never hold a lock during work
    }
    printf("worker %d exiting\n", id);
    return NULL;
}

void print_job(void *arg) {
    printf("job %d running on thread %lu\n", *(int *)arg, (unsigned long)pthread_self());
    free(arg);
}

int main(void) {
    pthread_t threads[NUM_THREADS];
    int ids[NUM_THREADS];
    for (int i = 0; i < NUM_THREADS; i++) {
        ids[i] = i;
        pthread_create(&threads[i], NULL, worker, &ids[i]);
    }

    for (int i = 0; i < 20; i++) {
        int *job_id = malloc(sizeof(int));
        if (!job_id) { perror("malloc"); break; }
        *job_id = i;
        pool_submit(print_job, job_id);
    }

    pthread_mutex_lock(&q_lock);
    shutdown_flag = 1;
    pthread_cond_broadcast(&q_not_empty);   // wake ALL workers so they can see shutdown_flag
    pthread_mutex_unlock(&q_lock);

    for (int i = 0; i < NUM_THREADS; i++) pthread_join(threads[i], NULL);
    return 0;
}
```

### 23.2 Mini project — a `select()`-based multi-client chat server

This fixes Section 21.1's "only one client at a time" limitation without
threads at all — a single-threaded event loop, the same core idea behind
`libuv` (Node.js), `libevent`, and Redis's `ae.c` (Section 34).

```c
// chat_server.c
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/select.h>

#define MAX_CLIENTS 32

int main(void) {
    int server_fd = socket(AF_INET, SOCK_STREAM, 0);
    int opt = 1;
    setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof opt);

    struct sockaddr_in addr = { .sin_family = AF_INET, .sin_port = htons(9000),
                                 .sin_addr.s_addr = INADDR_ANY };
    bind(server_fd, (struct sockaddr *)&addr, sizeof addr);
    listen(server_fd, 10);

    int clients[MAX_CLIENTS] = {0};   // 0 means "empty slot"
    printf("chat server on :9000\n");

    while (1) {
        fd_set readfds;
        FD_ZERO(&readfds);
        FD_SET(server_fd, &readfds);
        int max_fd = server_fd;

        for (int i = 0; i < MAX_CLIENTS; i++) {
            if (clients[i] > 0) {
                FD_SET(clients[i], &readfds);
                if (clients[i] > max_fd) max_fd = clients[i];
            }
        }

        select(max_fd + 1, &readfds, NULL, NULL, NULL);   // block until ANY fd is ready

        if (FD_ISSET(server_fd, &readfds)) {               // new connection waiting
            int new_fd = accept(server_fd, NULL, NULL);
            for (int i = 0; i < MAX_CLIENTS; i++) {
                if (clients[i] == 0) { clients[i] = new_fd; break; }
            }
        }

        for (int i = 0; i < MAX_CLIENTS; i++) {
            if (clients[i] > 0 && FD_ISSET(clients[i], &readfds)) {
                char buf[512];
                ssize_t n = read(clients[i], buf, sizeof buf);
                if (n <= 0) {                               // client disconnected
                    close(clients[i]);
                    clients[i] = 0;
                    continue;
                }
                for (int j = 0; j < MAX_CLIENTS; j++) {      // broadcast to everyone else
                    if (clients[j] > 0 && clients[j] != clients[i]) {
                        write(clients[j], buf, (size_t)n);
                    }
                }
            }
        }
    }
}
```

`select()` scales poorly past a few hundred descriptors (it's O(n) per
call and has a hard `FD_SETSIZE` limit) — production event loops use
`epoll` (Linux) or `kqueue` (BSD/macOS) instead, which is exactly what
Redis's `ae.c` abstracts over in Section 34.

---

# Part IV — Real-world builds

> Theory is done — Parts I–III gave you everything this Part uses. From here
> it's full, opinionated builds rather than isolated examples: a production
> CLI, a text editor, a Unix shell, and a bump allocator, each combining
> memory management, data structures, and I/O into something closer to what
> you'd actually ship.

## 24. Building production CLIs: argument parsing, config precedence, exit codes

### 24.1 `getopt` — the standard way to parse flags

```c
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char *argv[]) {
    int verbose = 0;
    char *output = "out.txt";
    int opt;

    while ((opt = getopt(argc, argv, "vo:h")) != -1) {
        switch (opt) {
            case 'v': verbose = 1; break;
            case 'o': output = optarg; break;   // optarg is set by getopt for options that take a value
            case 'h':
                printf("usage: %s [-v] [-o file] [input...]\n", argv[0]);
                return 0;
            default:
                fprintf(stderr, "usage: %s [-v] [-o file] [input...]\n", argv[0]);
                return 2;   // convention: 2 for usage errors, 1 for runtime errors, 0 for success
        }
    }

    for (int i = optind; i < argc; i++) {   // optind marks where positional args start
        printf("input: %s\n", argv[i]);
    }

    printf("verbose=%d output=%s\n", verbose, output);
    return 0;
}
```

### 24.2 Config precedence — the convention many serious CLIs follow

```
defaults  <  config file  <  environment variables  <  command-line flags
(lowest precedence)                                    (highest precedence)
```

```c
#include <stdlib.h>
#include <string.h>

const char *resolve_output_dir(const char *cli_flag) {
    if (cli_flag) return cli_flag;                 // 1. explicit flag wins
    const char *env = getenv("APP_OUTPUT_DIR");
    if (env) return env;                            // 2. then environment
    return "./output";                              // 3. then a hardcoded default
    // a real tool would also check a config file here, between env and default
}
```

**Real-world example.** This precedence idea is why `git -c
user.name=X commit` overrides `~/.gitconfig`, which overrides git's
built-in defaults. Many serious CLIs (curl, docker, kubectl) use some
version of this layered-override convention, because it lets a script
override one setting for one invocation without touching global config.

### 24.3 Exit codes are part of your program's API

```c
// Conventionally (see <sysexits.h> on BSD/macOS for the fuller list):
// 0   = success
// 1   = general runtime error
// 2   = usage error (bad arguments)
// 126 = command found but not executable
// 127 = command not found
// 128+N = terminated by signal N
```

Shell scripts and CI pipelines branch on your exit code — returning `0`
after a caught, logged error (instead of `1`) is a real and common bug
class that silently turns failures into "success" in automation.

---

## 25. Capstone: a `kilo`-style terminal text editor

`kilo` (by Salvatore Sanfilippo, the same author as Redis) is a ~1000-line
text editor that fits in a single C file and is a rite of passage for
learning raw terminal I/O. Here is the core mechanism — putting the
terminal into **raw mode** — which is the part every tutorial-follower
finds most magical:

```c
#include <termios.h>
#include <unistd.h>
#include <ctype.h>
#include <stdio.h>
#include <stdlib.h>

static struct termios orig_termios;

void disable_raw_mode(void) {
    tcsetattr(STDIN_FILENO, TCSAFLUSH, &orig_termios);   // restore on exit
}

void enable_raw_mode(void) {
    tcgetattr(STDIN_FILENO, &orig_termios);
    atexit(disable_raw_mode);          // guarantee terminal is restored even on early exit

    struct termios raw = orig_termios;
    raw.c_lflag &= ~(ECHO | ICANON | ISIG | IEXTEN);  // no echo, read byte-by-byte, no Ctrl-C/Z, no Ctrl-V
    raw.c_iflag &= ~(IXON | ICRNL | BRKINT | INPCK | ISTRIP);
    raw.c_oflag &= ~(OPOST);           // don't translate \n to \r\n on output
    raw.c_cc[VMIN] = 0;                // read() returns as soon as any input is available
    raw.c_cc[VTIME] = 1;               // ...or after 100ms with nothing, whichever first

    tcsetattr(STDIN_FILENO, TCSAFLUSH, &raw);
}

int main(void) {
    enable_raw_mode();
    printf("Raw mode on. Press 'q' to quit.\r\n");

    char c;
    while (read(STDIN_FILENO, &c, 1) == 1 && c != 'q') {
        if (c == '\r') { printf("\r\n"); continue; }
        if (iscntrl((unsigned char)c)) printf("%d\r\n", c);
        else printf("%d ('%c')\r\n", c, c);
    }
    return 0;   // disable_raw_mode runs automatically via atexit
}
```

By default, your terminal is in **canonical mode**: it buffers a whole
line and does its own editing (backspace, etc.) before your program ever
sees a byte, and Ctrl-C sends `SIGINT` before you can intercept it. Raw
mode turns all of that off so *your program* becomes the terminal driver
— this is exactly what `vim`, `nano`, `less`, and every TUI application
does on startup, and exactly why they all restore the terminal on exit
(via `atexit`, as here, or a signal handler) — a crash that skips this
step is the classic "my terminal is broken, why does nothing echo" bug
you fix by blindly typing `reset` and pressing enter.

From here, a full `kilo` adds: reading the whole file into a line array
(the dynamic-array pattern from Section 16), rendering a screen buffer
with ANSI escape codes, and mapping arrow-key escape sequences
(`\x1b[A`, `\x1b[B`, …) to cursor movement — all built from primitives
you already have.

---

## 26. Capstone: a Unix mini-shell with pipes and redirection

```c
// minishell.c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/wait.h>
#include <fcntl.h>

#define MAX_ARGS 64

static char **tokenize(char *line) {
    char **args = malloc(MAX_ARGS * sizeof(char *));
    if (!args) return NULL;
    int i = 0;
    char *tok = strtok(line, " \t\n");
    while (tok && i < MAX_ARGS - 1) { args[i++] = tok; tok = strtok(NULL, " \t\n"); }
    args[i] = NULL;
    return args;
}

static void run_command(char **args) {
    if (!args[0]) return;

    // handle output redirection: `cmd args... > file`
    char *outfile = NULL;
    for (int i = 0; args[i]; i++) {
        if (strcmp(args[i], ">") == 0) { outfile = args[i + 1]; args[i] = NULL; break; }
    }

    pid_t pid = fork();     // duplicate the current process — pid==0 in the child
    if (pid == 0) {
        if (outfile) {
            int fd = open(outfile, O_WRONLY | O_CREAT | O_TRUNC, 0644);
            if (fd < 0) { perror("open"); exit(1); }
            if (dup2(fd, STDOUT_FILENO) < 0) { perror("dup2"); close(fd); exit(1); }
            close(fd);
        }
        execvp(args[0], args);         // replace THIS process's image with the new program
        perror("execvp");              // only reached if execvp itself failed
        exit(127);
    } else if (pid > 0) {
        int status;
        waitpid(pid, &status, 0);      // parent blocks until the child finishes
    } else {
        perror("fork");
    }
}

int main(void) {
    char line[1024];
    while (1) {
        printf("minishell> ");
        fflush(stdout);
        if (!fgets(line, sizeof line, stdin)) break;    // EOF (Ctrl-D)
        if (strcmp(line, "\n") == 0) continue;

        char **args = tokenize(line);
        if (!args) { perror("malloc"); return 1; }
        if (args[0] == NULL) { free(args); continue; }
        if (strcmp(args[0], "exit") == 0) { free(args); break; }
        if (strcmp(args[0], "cd") == 0) {                // cd must run in the shell itself,
            const char *dir = args[1] ? args[1] : getenv("HOME");
            if (!dir || chdir(dir) != 0) perror("cd");     // not a child — a child's chdir wouldn't
            free(args); continue;                          // affect the parent shell's directory
        }

        run_command(args);
        free(args);
    }
    return 0;
}
```

`fork()` + `execvp()` is *the* Unix process-creation model, and it's why
`cd` and `exit` must be **built into the shell** rather than external
programs — a child process's working directory change or `exit()` call
has zero effect on its parent. Piping (`cmd1 | cmd2`) extends this exact
pattern: `pipe()` creates a connected read/write fd pair, then each child
`dup2`s one end onto its stdin or stdout before `execvp` — the same
`dup2` trick used for the `>` redirection above, just wired to the other
child's fd instead of a file.

---

## 27. Capstone: a bump/arena allocator

Real programs that allocate thousands of short-lived objects (parsers,
per-request web-server state, game-engine per-frame allocations) often
skip `malloc`/`free` entirely per-object and use an **arena**: allocate
one big block up front, hand out slices of it by just bumping a pointer,
and free the *entire arena* at once when the whole batch of work is done.

```c
// arena.h
#ifndef ARENA_H
#define ARENA_H
#include <stddef.h>

typedef struct {
    char  *base;
    size_t size;
    size_t offset;
} Arena;

Arena arena_new(size_t size);
void *arena_alloc(Arena *a, size_t size);
void  arena_reset(Arena *a);     // "free everything" — just rewind the pointer
void  arena_destroy(Arena *a);

#endif
```

```c
// arena.c
#include "arena.h"
#include <stdlib.h>
#include <stdint.h>

Arena arena_new(size_t size) {
    return (Arena){ .base = malloc(size), .size = size, .offset = 0 };
}

void *arena_alloc(Arena *a, size_t size) {
    if (!a->base || a->offset > SIZE_MAX - 7) return NULL;
    size_t aligned = (a->offset + 7) & ~(size_t)7;   // 8-byte align every allocation
    if (aligned > a->size || size > a->size - aligned) return NULL;  // out of arena space
    void *ptr = a->base + aligned;
    a->offset = aligned + size;
    return ptr;
}

void arena_reset(Arena *a) { a->offset = 0; }         // O(1) "free" of every allocation at once

void arena_destroy(Arena *a) { free(a->base); a->base = NULL; }
```

```c
#include <stdio.h>
#include "arena.h"

int main(void) {
    Arena a = arena_new(1024 * 1024);
    if (!a.base) { perror("arena_new"); return 1; }

    for (int batch = 0; batch < 3; batch++) {
        for (int i = 0; i < 1000; i++) {
            int *n = arena_alloc(&a, sizeof(int));
            if (!n) { fprintf(stderr, "arena exhausted\n"); arena_destroy(&a); return 1; }
            *n = i;
        }
        printf("batch %d done, resetting\n", batch);
        arena_reset(&a);       // one call frees all 1000 ints from this batch — no per-object free
    }

    arena_destroy(&a);
    return 0;
}
```

There is **no `arena_free(ptr)` for individual allocations** — that's the
entire point: you trade fine-grained control for O(1) bulk deallocation
and zero fragmentation. This is precisely why compilers (each function's
AST nodes freed all at once after codegen), game engines (per-frame
allocations reset every frame), and request-scoped web-server code
(everything freed when the HTTP response is sent) all use this pattern
instead of individually `free()`-ing thousands of small objects.

---

# Part V — Security-focused C

> A different angle on everything so far: not building features, but not
> shooting yourself in the foot. Undefined behavior is the one topic every C
> programmer needs memorized rather than just referenced — the rest of this
> Part (safe string handling, hardening flags, fuzzing, static analysis) is
> the tooling that catches what memorization alone misses.

## 28. Undefined behavior: the list every C programmer must memorize

Undefined behavior (UB) means the standard places **no requirement** on
what happens — not "it crashes," but "the compiler is allowed to assume
it never happens and optimize accordingly," which is often more
surprising than a crash.

| UB | Example | Why it matters |
|---|---|---|
| Signed integer overflow | `INT_MAX + 1` | Compiler may assume `x + 1 > x` always and eliminate an overflow check you wrote |
| Out-of-bounds access | `arr[10]` on a 10-element array (valid indices 0–9) | No bounds checking exists; you silently corrupt adjacent memory |
| Use of uninitialized value | `int x; printf("%d", x);` | Value is indeterminate; sanitizers catch it, plain runs may "work" until they don't |
| NULL pointer dereference | `*(int*)NULL` | Often segfaults, but the compiler may also assume it's unreachable and delete surrounding code |
| Use-after-free / dangling pointer | Section 10.2 | Silent corruption, sometimes exploitable |
| Data race | Section 20.1 | No defined result; can tear values, reorder observably |
| Modifying a `const` object | casting away `const` and writing | May write to read-only memory (segfault) or silently do nothing, depending on placement |
| Missing `return` in non-`void` function | falls off the end | Caller gets a garbage value |
| Strict aliasing violation | reading a `float*` through an `int*` | Compiler may reorder/omit the access assuming the pointers can't alias |
| Shift by ≥ the type's width | `1 << 32` on a 32-bit `int` | Undefined behavior, not "0" as many assume |

**War story.** The signed-overflow row above is not theoretical: a real
class of Linux kernel bugs came from code like `if (ptr + len < ptr)` to
detect pointer-arithmetic overflow — because pointer overflow (and
signed overflow generally) is UB, compilers began optimizing this
*entire check away*, reasoning "overflow can't happen, so this comparison
is always false." Code that had defended against overflow for years
silently stopped defending against it after a compiler upgrade, with no
warning. The fix across the industry was to switch to
`__builtin_add_overflow` (GCC/Clang builtin) or unsigned arithmetic with
explicit wraparound semantics, both of which are *defined* behavior the
compiler cannot optimize away.

```c
#include <stdio.h>

int safe_add(int a, int b, int *result) {
    return __builtin_add_overflow(a, b, result);  // returns 1 if overflow occurred, defined behavior
}

int main(void) {
    int result;
    if (safe_add(2147483647, 1, &result)) {
        printf("overflow detected, not UB\n");
    }
    return 0;
}
```

---

## 29. Buffer overflows, safe string handling, and CERT C rules

### 29.1 Anatomy of a stack buffer overflow

```c
#include <stdio.h>
#include <string.h>

void vulnerable(const char *input) {
    char buf[16];
    strcpy(buf, input);   // BUG — no bound; if input is longer than 15 chars,
                            // this writes past buf, into adjacent stack memory,
                            // potentially overwriting the saved return address
}
```

On the stack, `buf` sits near the **saved return address** — the location
the CPU jumps to when `vulnerable` returns. A long enough `input`
overwrites that address with attacker-chosen bytes, and when the function
returns, execution jumps wherever the attacker wanted. This is the
classic "stack smashing" exploit primitive behind decades of remote-code-
execution CVEs, and exactly why:

```c
void safer(const char *input) {
    char buf[16];
    snprintf(buf, sizeof buf, "%s", input);   // bounded and always null-terminated when size > 0
}
```

### 29.2 CERT C's short list that catches most real bugs

| Rule (abbreviated) | In practice |
|---|---|
| STR31-C: Guarantee storage for strings is sufficient | Compute needed size *before* writing, or use bounded functions |
| ARR30-C/ARR38-C: Don't form/use out-of-bounds pointers | Check indices; never trust external input as an index |
| MEM30-C/MEM31-C: Don't access freed memory; free exactly once | Null out pointers after `free` (defense in depth), see Section 10.2 |
| INT30-C/INT32-C: Ensure unsigned/signed arithmetic doesn't wrap unexpectedly | Use `__builtin_*_overflow` or check bounds before arithmetic |
| EXP34-C: Don't dereference null pointers | Check every `malloc`/`fopen`/lookup return before using it |
| FIO34-C: Distinguish EOF from a valid data value | Check `fgetc`'s return against `EOF` *before* casting to `char` |

### 29.3 A worked hardening example: bounded, checked input parsing

```c
#include <stdio.h>
#include <stdlib.h>
#include <limits.h>

// Turns a real CVE-shaped bug ("trust the length byte from the network")
// into checked, bounded code.
int parse_length_prefixed(const unsigned char *data, size_t data_len,
                           unsigned char *out, size_t out_cap) {
    if (data_len < 2) return -1;                       // not enough bytes for a length prefix
    unsigned int claimed_len = (unsigned)(data[0] << 8 | data[1]);  // 16-bit big-endian length

    if (claimed_len > data_len - 2) return -1;          // claimed length exceeds what we actually have
    if (claimed_len > out_cap) return -1;               // claimed length exceeds caller's buffer

    memcpy(out, data + 2, claimed_len);
    return (int)claimed_len;
}
```

Every one of those three `if` checks corresponds to a real historical
CVE class: trusting a length prefix without checking it against both the
actual buffer received *and* the destination buffer's capacity is the
Heartbleed pattern from Section 1.4, generalized.

---

## 30. Hardening: compiler flags, canaries, ASLR, `_FORTIFY_SOURCE`

```bash
clang -std=c17 -O2 \
      -D_FORTIFY_SOURCE=2 \
      -fstack-protector-strong \
      -fPIE -pie \
      -Wl,-z,relro,-z,now \
      -o hardened_app app.c
```

The linker flags above are for Linux/ELF toolchains. On macOS or Windows,
use the platform's equivalent hardening options instead of copying the
`-Wl,-z,...` part verbatim.

| Flag | Defense |
|---|---|
| `-D_FORTIFY_SOURCE=2` | Adds runtime bounds checks to `memcpy`/`strcpy`/`sprintf`/etc. when the compiler can determine the destination size — turns some overflows into an immediate abort instead of silent corruption |
| `-fstack-protector-strong` | Inserts a random "canary" value before the return address; if a buffer overflow overwrites it, the program detects the corruption before returning and aborts |
| `-fPIE -pie` | Position-independent executable — required for ASLR (address space layout randomization) to randomize the binary's own load address, not just libraries |
| `-Wl,-z,relro,-z,now` | Makes the GOT (global offset table) read-only after startup, closing off a common exploit technique (GOT overwrite) |

**Real-world example.** These are representative of the hardening options
Linux distributions such as Debian, Fedora, and Ubuntu enable through their
package build systems. They are not magic, but they often turn some
memory-corruption bugs from "attacker gets control" into "process aborts."

---

## 31. Fuzzing a real parsing bug with libFuzzer

Fuzzing feeds a function millions of randomly-mutated inputs, guided by
code-coverage feedback, looking for inputs that crash it or trip a
sanitizer. It is *the* technique that found the bug class behind
Heartbleed-style CVEs across the industry after the fact (Google's
OSS-Fuzz alone has found 10,000+ bugs in open-source C/C++ projects).

```c
// fuzz_target.c — libFuzzer calls this function with random byte arrays
#include <stdint.h>
#include <stddef.h>
#include <string.h>

int parse_length_prefixed(const unsigned char *data, size_t data_len,
                           unsigned char *out, size_t out_cap); // from Section 29.3

int LLVMFuzzerTestOneInput(const uint8_t *data, size_t size) {
    unsigned char out[256];
    parse_length_prefixed(data, size, out, sizeof out);
    return 0;
}
```

```bash
clang -std=c17 -g -fsanitize=fuzzer,address -o fuzz_target fuzz_target.c parser.c
./fuzz_target                     # runs indefinitely, mutating inputs, until it finds a crash
./fuzz_target crash-abc123        # re-run just the one input that crashed, to reproduce/debug
```

Feed it the version of `parse_length_prefixed` *without* the bounds
checks from Section 29.3 and libFuzzer will find the overflow in seconds,
saving the exact crashing input to disk — this before/after comparison
is the single most convincing five-minute demo of why input validation
matters in C.

---

## 32. Static analysis: clang-tidy, cppcheck, the Clang Static Analyzer

```bash
# Catches style + some real bugs (uninitialized use, modernize suggestions)
clang-tidy app.c -- -std=c17

# Path-sensitive analysis: finds null derefs, leaks, dead stores by
# actually simulating execution paths through your code
scan-build clang -std=c17 -c app.c

# A fast, dependency-free linter — great as a pre-commit hook
cppcheck --enable=all --std=c17 app.c
```

**Real-world example.** `scan-build` (the Clang Static Analyzer's driver)
is what will flag Section 10.2's memory leak *without running the
program at all* — it walks every path through `leak()`, sees that `buf`
is assigned from `malloc` and never reaches a `free` on any path, and
reports it at the exact allocation line. Static analysis and the
dynamic tools from Section 14/31 (ASan/UBSan/fuzzing) are complementary:
static analysis finds bugs on paths your tests never execute; dynamic
tools find bugs your static analyzer can't reason about (attacker-
controlled runtime values, complex aliasing). Serious C projects run
both, plus a real test suite, in CI on every commit.

---

# Part VI — Open-source package walkthroughs

> This Part takes every concept from Parts I–V and points at where it
> lives in real, widely-deployed C codebases. The goal is to make reading
> unfamiliar C — the actual day-to-day skill of an experienced C
> engineer — feel approachable. Snippets here are simplified/paraphrased
> to teach the *design*, not copy-pasted verbatim; read the linked real
> files for the authoritative current source.

## 33. SQLite: the amalgamation, the VFS layer, opcode-based execution

**What it is:** the most widely deployed database engine in the world
(every Android/iOS app, every browser, git, Python's stdlib all embed
it) — a single C library, no server process, one file on disk.

**Concept mapping:**

- **The amalgamation** (`sqlite3.c`, ~250K lines generated from the real
  multi-file source tree) is Section 13's multi-file build taken to its
  logical extreme *in reverse*: SQLite ships as one giant translation
  unit specifically so the compiler can inline and optimize across
  what would otherwise be file boundaries — a deliberate trade of
  Section 13's modularity for whole-program optimization, appropriate
  because SQLite is a leaf dependency, not code you edit daily.
- **The VFS (virtual file system) layer** is Section 11.2's function-
  pointer-table pattern, at the center of SQLite's entire portability
  story: every OS interaction (`open`, `read`, `write`, `lock`) goes
  through a struct of function pointers (`sqlite3_vfs`), so porting
  SQLite to a new OS means implementing one struct, not touching the
  other 250K lines:

```c
// simplified sketch of SQLite's real sqlite3_io_methods design
typedef struct sqlite3_io_methods {
    int (*xRead)(void *file, void *buf, int amt, long long offset);
    int (*xWrite)(void *file, const void *buf, int amt, long long offset);
    int (*xSync)(void *file, int flags);
    int (*xLock)(void *file, int lockType);
    int (*xUnlock)(void *file, int lockType);
    // ... this table IS the abstraction boundary between SQLite and any given OS
} sqlite3_io_methods;
```

- **SQL execution as a bytecode VM**: SQLite compiles SQL text into an
  opcode sequence (visible via `EXPLAIN`) and runs it on a small virtual
  machine — the same architecture as Section 22.2's `_Generic` dispatch
  generalized into a real interpreter loop (a giant `switch` over opcode
  values, each `case` implementing one instruction), and the same idea
  behind CPython's bytecode interpreter, just written in C instead of
  being one.
- **Ownership convention**: `sqlite3_open()`/`sqlite3_close()`,
  `sqlite3_prepare_v2()`/`sqlite3_finalize()` — Section 10.3's `_new`/
  `_free` pairing, applied consistently across SQLite's entire public
  API, is why every SQLite tutorial's first lesson is "always call
  `sqlite3_close`."

```bash
# Try it yourself:
sqlite3 :memory: "EXPLAIN SELECT 1;"   # see the actual opcodes for a trivial query
```

## 34. Redis: the `ae.c` event loop, `sds` dynamic strings, object encoding

**What it is:** an in-memory data-structure server, written by Salvatore
Sanfilippo (also `kilo`'s author, Section 25) — single-threaded for all
command execution, which sidesteps Section 20's entire locking problem
by design.

**Concept mapping:**

- **`ae.c`, Redis's event loop**, is Section 23.2's `select()`-based chat
  server, generalized: it abstracts over `epoll` (Linux), `kqueue`
  (macOS/BSD), and `select` (fallback) behind one API
  (`aeCreateFileEvent`, `aeMain`), exactly the way SQLite's VFS abstracts
  over filesystems. This single-threaded event loop, not a thread pool,
  is *why* Redis commands are atomic without any mutex: only one command
  ever executes at a time, so Section 20's entire race-condition class
  simply cannot occur in command logic.
- **`sds` (Simple Dynamic Strings)**, `sds.h`/`sds.c`, is Section 7's
  "a string is a convention" problem, solved once, everywhere in the
  codebase: an `sds` is a `char *` (so it's still compatible with every
  libc string function), but the bytes *before* that pointer hold a
  header with the length and allocated capacity — the same "hidden
  header before the returned pointer" trick real allocators use
  internally (Section 39).

```c
// simplified sketch of the real sds design
struct sdshdr {
    uint32_t len;     // used length — makes strlen() O(1) instead of O(n)
    uint32_t alloc;   // allocated capacity — makes "do we need to grow?" O(1)
    char     buf[];   // C99 flexible array member — the actual string bytes
};
// sdsnew() returns &hdr->buf[0], NOT &hdr — so callers use it exactly like char*,
// while sdslen(s) reaches backward: ((struct sdshdr*)(s - sizeof(struct sdshdr)))->len
```

  This flexible-array-member technique (`buf[]` as the struct's last
  member, with the actual allocation sized to `sizeof(struct sdshdr) +
  needed_bytes`) is standard C99 and appears constantly in real systems
  code for exactly this "one allocation, header + variable-length
  payload" shape.
- **Object encoding**: Redis's `robj` uses Section 8.3's tagged-union
  pattern — a `type` field plus an `encoding` field that says whether a
  "list" is currently a compact `ziplist` (few, small elements) or a
  full linked list (many elements), transparently upgrading as it grows.
  This is the real-world justification for tagged unions: pick the
  cheapest representation that fits the current data, dispatch on the
  tag everywhere you touch it.

## 35. curl: the easy/multi handle design and callback-driven I/O

**What it is:** the HTTP(S)/FTP/... client library behind nearly every
language's HTTP bindings and an enormous share of embedded devices.

**Concept mapping:**

- **Opaque handles**: `CURL *curl = curl_easy_init();` returns a pointer
  to a struct whose *definition* is private to curl's `.c` files —
  callers only ever see a `CURL *` and interact through
  `curl_easy_setopt`/`curl_easy_perform`/`curl_easy_cleanup`. This is
  Section 13's `static`-for-privacy idea taken to the API level: the
  struct layout can change between curl versions without breaking
  binaries linked against it, because no caller code ever reads a field
  by name — this is Section 40's "designing a stable ABI," and it's why
  curl has kept API/ABI compatibility for over two decades.
- **Callback-driven I/O**: `CURLOPT_WRITEFUNCTION` is exactly Section
  11.2's function-pointer parameter — you hand curl a function pointer,
  and curl calls it with each chunk of the response body as it arrives,
  instead of you calling `read()` in a loop yourself:

```c
size_t write_cb(char *ptr, size_t size, size_t nmemb, void *userdata) {
    FILE *out = userdata;
    return fwrite(ptr, size, nmemb, out);   // return value tells curl how many bytes you consumed
}

// curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_cb);
// curl_easy_setopt(curl, CURLOPT_WRITEDATA, out_file);
```

- **The multi interface** (`curl_multi_*`) is Section 21/23's event-loop
  idea applied to *outgoing* connections: instead of blocking on one
  request at a time (`curl_easy_perform`), the multi API lets one thread
  drive hundreds of concurrent requests through a `select()`/`poll()`-
  style loop (`curl_multi_fdset`), the same non-blocking-I/O pattern as
  Redis's `ae.c` and Section 23.2's chat server, applied to a client
  instead of a server.

## 36. cJSON: recursive-descent parsing and tagged unions

**What it is:** a small (~2500-line), extremely widely embedded JSON
library — if Section 16.3's mini JSON parser felt approachable, cJSON is
the natural next read: same recursive-descent structure, production-
hardened.

**Concept mapping, directly against your own Section 16.3 parser:**

- `cJSON` the struct is Section 8.3's tagged union, just flattened into
  one struct with a `type` int field and multiple payload fields instead
  of a `union` block — a deliberate simplicity trade-off (slightly more
  memory per node) for slightly simpler code than a real `union`.
- `cJSON_Parse` → `parse_value` → `parse_object`/`parse_array`/
  `parse_string`/`parse_number` is *exactly* your Section 16.3 function
  breakdown, function-for-function, because that decomposition is simply
  what JSON's grammar requires — reading cJSON's actual `cJSON_Parse*`
  functions after writing your own is one of the fastest ways to feel
  the jump from "I can write basic C" to "I can read real C."
- `cJSON_Delete` walking the tree and freeing children before parents is
  Section 16.3's `json_free`, and it's the same recursive-teardown shape
  as Section 18.1's `list_free` and Section 18.3's tree traversal —
  once you've written one recursive free function, you can read every
  other one.

## 37. uthash and klib: generics via macros

**What they are:** two of the most-vendored single-header C libraries on
GitHub — `uthash.h` adds a hash table to *any* existing struct via
macros; `klib`'s `khash.h`/`kvec.h` do the same for hash tables and
vectors with a slightly different macro style. Both solve exactly
Section 16.1's "C has no generics" problem, at library scale.

```c
// uthash: turn an existing struct into a hash-table entry by adding ONE field
#include "uthash.h"

typedef struct {
    int id;              // the key
    char name[32];
    UT_hash_handle hh;   // <-- makes this struct hashable; this IS Section 19's
                          //     intrusive-container idea, applied to a hash table
                          //     instead of a linked list
} User;

User *users = NULL;   // the hash table itself is just a NULL pointer to start

void add_user(User *u) { HASH_ADD_INT(users, id, u); }
User *find_user(int id) { User *u; HASH_FIND_INT(users, &id, u); return u; }
```

`UT_hash_handle hh` embedded in `User` is precisely Section 19's
`container_of` technique — `uthash.h`'s macros compute back from `hh`'s
address to the owning `User*`, the exact same trick as the Linux kernel's
`list_head`, just wrapping a hash table's bucket/chain pointers instead
of a doubly-linked list's `next`/`prev`. Once you've built Section 19's
intrusive list by hand, `uthash.h`'s ~1000 lines of macros stop looking
like magic and start looking like "the thing I already understand, at
production polish level."

---

# Part VII — Expert: internals, allocators, API design

## 38. The C11 memory model: atomics, ordering, lock-free basics

Section 20.4 used `atomic_fetch_add` with its default ordering
(`memory_order_seq_cst` — sequentially consistent, the easiest to reason
about, the most expensive on weakly-ordered hardware like ARM). Real
lock-free code often relaxes this deliberately for performance:

```c
#include <stdatomic.h>

atomic_int ready = 0;
int        payload = 0;   // NOT atomic — protected by the ordering below, not by its own type

void producer(void) {
    payload = 42;                                       // 1. write the data
    atomic_store_explicit(&ready, 1, memory_order_release); // 2. THEN publish "it's ready"
}

int consumer(void) {
    while (!atomic_load_explicit(&ready, memory_order_acquire)) { } // spin until published
    return payload;    // guaranteed to see 42 — the acquire/release pair forms a "happens-before" edge
}
```

`memory_order_release` on the store and `memory_order_acquire` on the
load form a synchronization point: everything the producer wrote
*before* the release-store (here, `payload = 42`) is guaranteed visible
to the consumer *after* its acquire-load sees the new value. Without this
pairing (e.g., using `memory_order_relaxed` on both), the compiler and
CPU are both free to reorder the two writes, and the consumer could
observe `ready == 1` while still reading `payload`'s old value — a bug
that's nearly impossible to reproduce in testing and shows up only under
production load on weakly-ordered hardware.

| Ordering | Cost | Use for |
|---|---|---|
| `memory_order_relaxed` | Cheapest | Independent counters/stats where ordering relative to other data doesn't matter |
| `memory_order_acquire`/`_release` | Moderate | Publishing a pointer/flag after finishing writes to the data it guards (the pattern above) |
| `memory_order_seq_cst` | Most expensive | Default/safe choice; use until you've measured a real bottleneck |

**Real-world example.** This exact acquire/release handshake is how
lock-free single-producer/single-consumer ring buffers (used in
audio-processing callbacks, kernel-bypass networking, and Redis-adjacent
high-throughput queues) hand data between threads without ever taking a
mutex — the mutex-based producer/consumer from Section 20.3, with the
lock replaced by two atomics and a lot more care.

## 39. Writing a real allocator: free lists, coalescing, and `dlmalloc`'s ideas

`malloc`/`free` (Section 10) aren't magic — they're a C library on top of
the OS's `mmap`/`sbrk` syscalls. Here's a simplified allocator using the
same free-list-with-coalescing idea as Doug Lea's `dlmalloc` (the design
most production allocators, including glibc's, descend from):

```c
#include <stddef.h>
#include <sys/mman.h>

typedef struct Block {
    size_t       size;    // payload size, not including this header
    int          free;
    struct Block *next;   // next block in address order (for coalescing)
} Block;

#define HEADER_SIZE sizeof(Block)
static Block *heap_start = NULL;

static Block *request_from_os(size_t size) {
    size_t total = size + HEADER_SIZE;
    Block *b = mmap(NULL, total, PROT_READ | PROT_WRITE,
                     MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    if (b == MAP_FAILED) return NULL;
    b->size = size;
    b->free = 0;
    b->next = NULL;
    return b;
}

void *my_malloc(size_t size) {
    Block *b = heap_start, *prev = NULL;
    while (b) {                                    // first-fit search of the free list
        if (b->free && b->size >= size) { b->free = 0; return (char *)b + HEADER_SIZE; }
        prev = b;
        b = b->next;
    }
    Block *new_block = request_from_os(size);       // nothing fit — ask the OS for more
    if (!new_block) return NULL;
    if (prev) prev->next = new_block; else heap_start = new_block;
    return (char *)new_block + HEADER_SIZE;
}

void my_free(void *ptr) {
    if (!ptr) return;
    Block *b = (Block *)((char *)ptr - HEADER_SIZE);  // walk back to the header — same trick as sds (Section 34)
    b->free = 1;

    // Coalesce adjacent free blocks so a big enough contiguous run is
    // available for future large allocations, instead of staying
    // fragmented into many small free blocks forever.
    Block *cur = heap_start;
    while (cur && cur->next) {
        if (cur->free && cur->next->free) {
            cur->size += HEADER_SIZE + cur->next->size;
            cur->next = cur->next->next;
        } else {
            cur = cur->next;
        }
    }
}
```

The `(char *)ptr - HEADER_SIZE` trick to recover the header from the
pointer `malloc` handed out is the *exact same pattern* as Section 34's
`sds` header and Section 19's `container_of` — "hide bookkeeping data
immediately before/around the pointer you give the caller" is one idea
that recurs across allocators, dynamic strings, and intrusive containers
because it lets the caller treat the returned pointer as if it were the
raw payload, with zero extra indirection.

**Real-world example.** First-fit (used above) is the simplest strategy;
real allocators use segregated free lists (separate lists per size
class, so small allocations never have to scan past large free blocks)
and per-thread arenas (so threads don't contend on one global free list)
— exactly the optimizations that separate `dlmalloc` from `ptmalloc2`
(glibc's allocator) from `jemalloc`/`tcmalloc` (used by Firefox/Chrome
and many high-throughput servers respectively).

## 40. ABI, linking, and designing a stable C API

### 40.1 Opaque pointers: the technique behind Section 35's `CURL*`

```c
// mylib.h — the public header. Note: no struct fields visible at all.
typedef struct MyContext MyContext;

MyContext *mycontext_create(void);
void       mycontext_set_option(MyContext *ctx, int key, int value);
void       mycontext_destroy(MyContext *ctx);
```

```c
// mylib.c — the private definition, invisible to anyone who only has the header
struct MyContext {
    int option_a;
    int option_b;
    void *internal_state;
};

MyContext *mycontext_create(void) { return calloc(1, sizeof(struct MyContext)); }
void mycontext_destroy(MyContext *ctx) { free(ctx); }
```

Because callers never see `struct MyContext`'s layout, you can add
fields, reorder them, or change `internal_state`'s type in a later
release **without breaking binaries already linked against the old
version** — only `mylib.c` needs to be recompiled. This is precisely
why curl, SQLite (`sqlite3*`), and OpenSSL (`SSL_CTX*`) all expose
opaque pointers instead of public structs: it's the C equivalent of an
interface, and it's the only realistic way to version a C library that
other people link against.

### 40.2 Symbol visibility: don't export what isn't API

```c
// Only symbols marked default-visibility (or not hidden) end up in the shared library's
// exported symbol table; everything else becomes a true implementation detail.
__attribute__((visibility("default")))
int mylib_public_function(void);

__attribute__((visibility("hidden")))
int mylib_internal_helper(void);   // callable within the library, invisible to linkers outside it
```

```bash
clang -fvisibility=hidden -shared -o libmylib.so mylib.c   # hide everything by default,
                                                             # opt individual symbols back in
```

**Real-world example.** This is why a well-built shared library's
`nm -D libfoo.so` shows only its intended public API, not every static
helper — smaller exported symbol tables load faster, can't be
accidentally depended on by other code, and let the library's authors
freely refactor internals (Section 40.1's whole point) without an API
review.

---

## 41. Cross-compilation and build systems at scale (CMake/Meson)

### 41.1 Cross-compiling — building for a target that isn't your machine

```bash
# Building for a Raspberry Pi (ARM) from an x86_64 or Apple Silicon dev machine:
aarch64-linux-gnu-gcc -std=c17 -o app app.c
file app   # ELF 64-bit LSB executable, ARM aarch64 — even though you're not on ARM
```

This works because a "compiler" is usefully thought of as three separable pieces — a
front end (parses C), a code generator (targets a specific instruction
set), and system headers/libraries for the *target* OS — and a
cross-compiler points the last two at the target instead of the host.
This is the basic shape behind embedded-Linux, Android NDK, and firmware
builds: developers often build on powerful workstations and produce
binaries for ARM/RISC-V devices that may be too small or inconvenient to
run the compiler itself.

### 41.2 CMake — the de facto standard for C/C++ projects beyond one Makefile

```cmake
# CMakeLists.txt
cmake_minimum_required(VERSION 3.20)
project(myapp C)

set(CMAKE_C_STANDARD 17)
set(CMAKE_C_STANDARD_REQUIRED ON)

add_compile_options(-Wall -Wextra -Wpedantic -Werror)

add_executable(myapp src/main.c src/util.c)
target_include_directories(myapp PRIVATE include)

enable_testing()
add_executable(myapp_tests tests/test_util.c src/util.c)
add_test(NAME UtilTests COMMAND myapp_tests)
```

```bash
cmake -B build -DCMAKE_BUILD_TYPE=Debug
cmake --build build
ctest --test-dir build
```

CMake's job is to *generate* the actual build files (Makefiles, Ninja
files, Xcode/Visual Studio projects) from one portable description.
Many large C projects use CMake, Autotools, Meson, or similar build
systems instead of one hand-written `Makefile`, because cross-platform
builds quickly need compiler, OS, dependency, and feature-detection logic.

---

## 42. Profiling: `perf`, cache-aware layout, and false sharing

### 42.1 Finding the actual bottleneck instead of guessing

```bash
perf record -g ./myapp
perf report                 # a sorted-by-time breakdown of where the CPU actually spent cycles
```

```bash
# macOS equivalent
sudo xctrace record --template 'Time Profiler' --launch -- ./myapp
```

**Real-world example.** "I think the bottleneck is the hash table" is a
guess; `perf report` showing 40% of samples inside `memcpy` called from
your logging function is evidence — profiling before optimizing is one of
the highest-leverage habits for avoiding wasted performance work.

### 42.2 Struct-of-arrays vs. array-of-structs — Section 8.1's padding lesson, at scale

```c
// Array-of-structs (AoS): natural, but wastes cache-line bandwidth if you
// only need one field across millions of elements
typedef struct { float x, y, z; int id; } ParticleAoS;
ParticleAoS particles[1000000];
// summing just .x touches every whole struct to use one field — wasted cache bandwidth

// Struct-of-arrays (SoA): every array is packed with only the field you need
typedef struct {
    float xs[1000000], ys[1000000], zs[1000000];
    int   ids[1000000];
} ParticlesSoA;
// summing xs[] walks the needed field contiguously — much better cache-line utilization
```

This is the kind of optimization physics engines, particle systems, and
columnar databases use, because CPUs fetch memory in cache-line chunks
(64 bytes is common on modern desktops/servers, but not guaranteed) —
Section 8.1's struct-padding awareness, scaled up to "how should a
million of these be laid out in memory."

### 42.3 False sharing — Section 20's atomics, sabotaged by layout

```c
// BUG (performance, not correctness): two threads updating DIFFERENT
// counters that happen to share a CPU cache line will invalidate each
// other's cache line on every write, serializing what should be
// independent, parallel work.
typedef struct {
    atomic_int counter_a;   // likely on the same 64-byte cache line as counter_b
    atomic_int counter_b;
} Counters;

// Fix: pad so each counter owns its own cache line
typedef struct {
    atomic_int counter_a;
    char       pad[60];      // pushes counter_b onto the next 64-byte cache line
    atomic_int counter_b;
} CountersPadded;
```

This is a real, measured effect (not a micro-optimization myth) — the
canonical demonstration is two threads each incrementing their own
`atomic_int` in a tight loop, which runs dramatically slower when the two
counters share a cache line than when padding separates them, purely
from cache-coherency traffic between cores.

---

## 43. Capstone: a mini Redis-like key-value store with a wire protocol

This capstone combines nearly every earlier Part: Section 21's sockets,
Section 18.2's hash table, Section 23.2's event-loop style, Section 16.3's
parsing technique (applied to a wire protocol instead of JSON), and
Section 40's API discipline.

**The protocol** — a simplified line-based version of Redis's real
RESP protocol:

```
SET key value\r\n   -> +OK\r\n
GET key\r\n         -> $value\r\n   (or $-1\r\n if missing)
DEL key\r\n         -> :1\r\n       (or :0\r\n if it wasn't present)
```

```c
// kvserver.c
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/select.h>
#include "hashtable_str.h"   // Section 18.2's hash table, adapted to store string values

#define MAX_CLIENTS 64
#define BUF_SIZE 4096

static void handle_command(HashTable *db, char *line, char *response, size_t resp_cap) {
    char cmd[16], key[128], value[BUF_SIZE];
    int n = sscanf(line, "%15s %127s %4095[^\r\n]", cmd, key, value);

    if (strcmp(cmd, "SET") == 0 && n == 3) {
        ht_set_str(db, key, value);
        snprintf(response, resp_cap, "+OK\r\n");
    } else if (strcmp(cmd, "GET") == 0 && n >= 2) {
        const char *val = ht_get_str(db, key);
        if (val) snprintf(response, resp_cap, "$%s\r\n", val);
        else     snprintf(response, resp_cap, "$-1\r\n");
    } else if (strcmp(cmd, "DEL") == 0 && n >= 2) {
        int deleted = ht_del_str(db, key);
        snprintf(response, resp_cap, ":%d\r\n", deleted);
    } else {
        snprintf(response, resp_cap, "-ERR unknown command\r\n");
    }
}

int main(void) {
    HashTable *db = ht_new(1024);

    int server_fd = socket(AF_INET, SOCK_STREAM, 0);
    int opt = 1;
    setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof opt);
    struct sockaddr_in addr = { .sin_family = AF_INET, .sin_port = htons(6380),
                                 .sin_addr.s_addr = INADDR_ANY };
    bind(server_fd, (struct sockaddr *)&addr, sizeof addr);
    listen(server_fd, 16);
    printf("mini-kv listening on :6380\n");

    int clients[MAX_CLIENTS] = {0};

    while (1) {
        fd_set readfds;
        FD_ZERO(&readfds);
        FD_SET(server_fd, &readfds);
        int max_fd = server_fd;
        for (int i = 0; i < MAX_CLIENTS; i++)
            if (clients[i] > 0) { FD_SET(clients[i], &readfds); if (clients[i] > max_fd) max_fd = clients[i]; }

        select(max_fd + 1, &readfds, NULL, NULL, NULL);

        if (FD_ISSET(server_fd, &readfds)) {
            int fd = accept(server_fd, NULL, NULL);
            for (int i = 0; i < MAX_CLIENTS; i++) if (clients[i] == 0) { clients[i] = fd; break; }
        }

        for (int i = 0; i < MAX_CLIENTS; i++) {
            if (clients[i] > 0 && FD_ISSET(clients[i], &readfds)) {
                char line[BUF_SIZE] = {0};
                ssize_t n = read(clients[i], line, sizeof line - 1);
                if (n <= 0) { close(clients[i]); clients[i] = 0; continue; }

                char response[BUF_SIZE];
                handle_command(db, line, response, sizeof response);
                write(clients[i], response, strlen(response));
            }
        }
    }
}
```

```bash
./kvserver &
printf 'SET name kumaran\r\n' | nc -q1 localhost 6380   # +OK
printf 'GET name\r\n'         | nc -q1 localhost 6380   # $kumaran
printf 'DEL name\r\n'         | nc -q1 localhost 6380   # :1
printf 'GET name\r\n'         | nc -q1 localhost 6380   # $-1
```

Where to take it next, each a direct application of an earlier Part:
persistence via an append-only log you replay on startup (a simplified
version of Redis's AOF); a background save via `fork()` (Section 26) so
the child writes a point-in-time snapshot while the parent keeps serving
requests, unaffected by the child's memory writes thanks to
copy-on-write; multiple databases via multiple hash tables; and
expiring keys via a min-heap of `(expire_time, key)` pairs checked once
per event-loop iteration. This is, concept for concept, how Redis itself
grew from a single-file prototype into the production system Section 34
walked through.

---

# Part VIII — Writing safe, leak-free C

> Parts I–VII taught you the mechanisms. This Part is the discipline
> layer on top: the leak shapes that actually show up in real bug
> trackers (rarely "forgot one `free`"), a checklist for writing new C
> defensively from the first line, and a bad-code/good-code catalog of
> the mistakes that keep recurring across beginner and intermediate C —
> each with the exact compiler flag or tool that would have caught it.

## 44. Memory leaks, in depth: every shape they actually take in production

Section 10.2 showed the textbook leak: allocate, never free. Real leaks
in real codebases almost never look that clean — they hide in error
paths, loops, and partial failures. Every example below is a pattern
you will meet again.

### 44.1 The error-path leak — the single most common real-world leak

A function acquires several resources in sequence; an early one succeeds,
a later one fails, and the function returns without releasing what it
already acquired.

```c
// BAD
typedef struct { FILE *log; char *buf1; char *buf2; } Ctx;

int setup_bad(Ctx *ctx) {
    ctx->log = fopen("/tmp/app.log", "w");
    if (!ctx->log) return -1;

    ctx->buf1 = malloc(1024);
    if (!ctx->buf1) return -1;      // LEAK: ctx->log is still open, never fclose'd

    ctx->buf2 = malloc(2048);
    if (!ctx->buf2) return -1;      // LEAK: ctx->log AND ctx->buf1 both leak here

    return 0;
}
```

```c
// GOOD — Section 4.2's single forward goto to one cleanup label
int setup_good(Ctx *ctx) {
    int rc = -1;
    ctx->log = NULL; ctx->buf1 = NULL; ctx->buf2 = NULL;

    ctx->log = fopen("/tmp/app.log", "w");
    if (!ctx->log) goto cleanup;

    ctx->buf1 = malloc(1024);
    if (!ctx->buf1) goto cleanup;

    ctx->buf2 = malloc(2048);
    if (!ctx->buf2) goto cleanup;

    rc = 0;
cleanup:
    if (rc != 0) {
        free(ctx->buf1);              // free(NULL) is always safe (Section 10.3) —
        free(ctx->buf2);               // no need to check which ones actually got allocated
        if (ctx->log) fclose(ctx->log);
    }
    return rc;
}
```

Every resource this function might hold gets exactly one release path,
reached from every failure point — this is *why* Section 4.2 introduced
`goto cleanup` as idiomatic C rather than an anti-pattern: the
alternative (repeating the right subset of cleanup calls before every
`return`) is what actually produces error-path leaks, because it's easy
to add a fourth resource later and forget to update three of the four
early-return blocks.

### 44.2 The lost-reference leak

```c
// BAD
void lost_reference_bad(void) {
    char *buf = malloc(100);
    buf = malloc(200);   // the 100-byte block's address is gone — nothing points to it anymore,
    free(buf);            // so it can never be freed for the rest of the program's life
}

// GOOD
void lost_reference_good(void) {
    char *buf = malloc(100);
    free(buf);             // release it BEFORE the variable that names it is reassigned
    buf = malloc(200);
    free(buf);
}
```

This is the leak hiding behind innocuous-looking code like
`result = process(result)` where `process` allocates a new buffer and
the caller forgot the old `result` needed freeing first — very common
in string-building loops (`s = concat(s, next_word)`, repeated).

### 44.3 The partial-cleanup leak in an array of allocations

```c
// BAD
char **alloc_rows_bad(int n, int width) {
    char **rows = malloc((size_t)n * sizeof(char *));
    for (int i = 0; i < n; i++) {
        rows[i] = malloc((size_t)width);
        if (!rows[i]) {
            free(rows);      // LEAK: rows[0..i-1] were allocated and are now unreachable —
            return NULL;      // freeing the array of POINTERS did not free what they pointed to
        }
    }
    return rows;
}

// GOOD
char **alloc_rows_good(int n, int width) {
    char **rows = malloc((size_t)n * sizeof(char *));
    if (!rows) return NULL;
    for (int i = 0; i < n; i++) {
        rows[i] = malloc((size_t)width);
        if (!rows[i]) {
            for (int j = 0; j < i; j++) free(rows[j]);  // unwind exactly what succeeded so far
            free(rows);
            return NULL;
        }
    }
    return rows;
}
```

This is the array version of 44.1's error-path leak, and it generalizes
to any structure built incrementally: a linked list built node-by-node,
a hash table resized bucket-by-bucket — the rule is always "the cleanup
path must free every element that succeeded before the one that failed,
not just the outer container."

### 44.4 The `realloc`-failure leak (revisiting Section 10.1)

```c
// BAD
void grow_bad(char **buf, size_t new_size) {
    *buf = realloc(*buf, new_size);   // if realloc fails, it returns NULL — and the ORIGINAL
}                                       // block's address, which *buf held, is now overwritten
                                        // and gone: a leak, not just "buf is now NULL"

// GOOD
int grow_good(char **buf, size_t new_size) {
    char *tmp = realloc(*buf, new_size);
    if (!tmp) return -1;    // *buf is untouched and still valid — free it if the caller gives up,
    *buf = tmp;               // or keep using the old, smaller buffer and try again later
    return 0;
}
```

### 44.5 Leaks aren't just `malloc` — every acquire needs a matching release

```c
// BAD — leaks a file descriptor and a mutex lock, not "memory," but the
// same resource-lifetime discipline applies
int read_config_bad(const char *path, pthread_mutex_t *lock) {
    pthread_mutex_lock(lock);
    FILE *f = fopen(path, "r");
    if (!f) return -1;              // LEAK: lock is never unlocked — every future
                                      // caller of this function deadlocks forever
    // ... read ...
    fclose(f);
    pthread_mutex_unlock(lock);
    return 0;
}
```

A held mutex that's never unlocked, a socket that's never `close`d, a
`mmap`ed region that's never `munmap`ped — all are "leaks" in the sense
that matters: a finite resource acquired and never released. The
ownership discipline from Section 10.3 (one `_new` for one `_free`, one
`lock` for one `unlock`, one `open` for one `close`) is the same fix for
all of them.

### 44.6 Finding leaks you already wrote

```bash
# Linux: the standard leak-triage tool, no recompilation needed
valgrind --leak-check=full --show-leak-kinds=all ./app

# Clang/GCC: LeakSanitizer, built into ASan on Linux (bundle with -fsanitize=address)
clang -std=gnu17 -g -fsanitize=address -o app app.c
ASAN_OPTIONS=detect_leaks=1 ./app     # prints each leak's allocation stack trace on exit
```

Valgrind's leak-kind vocabulary is worth knowing: `definitely lost`
(nothing in the program points to this block anymore — a real leak,
exactly Sections 44.1–44.3), `still reachable` (a global/static pointer
still points to it at exit — usually fine, e.g. a cache freed by process
exit instead of explicit code), and `possibly lost` (a pointer to the
*interior* of the block still exists, but not to its start — often a
custom allocator or intrusive-pointer pattern, worth a second look).

---

## 45. Writing safe C: a defensive-programming checklist

This is the consolidated discipline behind every "GOOD" example in this
guide — treat it as a pre-commit checklist for new C code, not a wall of
theory to read once.

| # | Rule | Where this guide showed why |
|---|---|---|
| 1 | Check every return value that can fail (`malloc`, `fopen`, `realloc`, syscalls) before using the result | Section 10.1, 12.3 |
| 2 | Bound every write into a fixed-size buffer — `snprintf`, not `sprintf`; `strlcpy`/manual length tracking, not `strcpy`/`strcat` | Section 7.2, 29 |
| 3 | Never pass externally influenced data as a `printf`-family format string | Section 46.4 below |
| 4 | Validate lengths from untrusted sources against *both* the source buffer's actual size and the destination's capacity, before any `memcpy` | Section 29.3 |
| 5 | Give every heap allocation exactly one owner and one release path; use `goto cleanup` for multi-resource functions | Section 4.2, 10.3, 44.1 |
| 6 | Null out or scope-limit pointers after `free`; never reuse a freed pointer | Section 10.2, 28 |
| 7 | Use `<stdint.h>` exact-width types for anything with a size that matters (protocols, file formats); never assume `int` is 32 bits | Section 3.1 |
| 8 | Prefer unsigned arithmetic with defined wraparound, or `__builtin_*_overflow`, over unchecked signed arithmetic where overflow is possible | Section 3.3, 28 |
| 9 | Initialize every variable at its point of declaration; never read a value before it's assigned | Section 28, 46.9 |
| 10 | Never return a pointer to a local (stack) variable | Section 5.1, 46.7 |
| 11 | Compile with `-Wall -Wextra -Wpedantic -Werror` from the first commit, not added later | Section 2.2, 46 (repeatedly) |
| 12 | Run the test suite under ASan/UBSan (and, on Linux, LeakSanitizer/Valgrind) before merging, not just before a release | Section 14, 44.6 |
| 13 | Keep mutable state `static`-private or instance-scoped; avoid file-scope globals a second thread or a re-entrant call path could race on | Section 13.1, 20.1, 46.13 |
| 14 | Treat every value crossing a trust boundary (network, file, CLI args, environment) as hostile until validated | Section 29, 31 |
| 15 | For anything parsing untrusted input, fuzz it — a parser without a fuzz target is a parser you haven't tested against the inputs that will actually break it | Section 31 |

### A worked example: applying the checklist to one function

```c
// Reads a length-prefixed record from an untrusted socket buffer.
// Every numbered comment below maps to a checklist rule above.
int read_record(const unsigned char *data, size_t data_len,
                 char *out, size_t out_cap) {
    if (data_len < 2) return -1;                              // (4) enough bytes for a length prefix?

    uint16_t claimed_len;
    memcpy(&claimed_len, data, sizeof claimed_len);            // (7) exact-width type for a wire field
    claimed_len = ntohs(claimed_len);                          // network byte order, always, for wire data

    size_t remaining = data_len - 2;
    if (claimed_len > remaining) return -1;                    // (4) trust neither side alone
    if ((size_t)claimed_len >= out_cap) return -1;             // (2) bound against the DESTINATION too

    memcpy(out, data + 2, claimed_len);
    out[claimed_len] = '\0';
    return (int)claimed_len;
}
```

This is Section 29.3's example again, but read it now as a checklist
walkthrough rather than a one-off fix: every line exists because of a
specific rule above, not because "more checks are always better" — a
defensively-written function should be able to point at *why* each
check exists.

---

## 46. How NOT to write C: a bad-code/good-code anti-pattern catalog

Every entry below compiles (the guide's opening promise holds here too —
"BAD" blocks compile, several with a warning your flags should be
catching; run them past `-Wall -Wextra` yourself to see your compiler
already knows about most of these).

### 46.1 Ignoring a return value that reports failure

```c
// BAD
FILE *f = fopen("config.json", "r");
fread(buf, 1, size, f);     // if fopen failed, f is NULL — this dereferences NULL

// GOOD
FILE *f = fopen("config.json", "r");
if (!f) { perror("fopen"); return -1; }
fread(buf, 1, size, f);
```

### 46.2 Unbounded string functions

```c
// BAD
char name[32];
strcpy(name, user_supplied);     // no bound — a longer input overflows `name` (Section 29.1)

// GOOD
char name[32];
snprintf(name, sizeof name, "%s", user_supplied);
```

### 46.3 Trusting `scanf("%s", ...)` with no width limit

```c
// BAD
char name[32];
scanf("%s", name);          // unbounded — identical risk to strcpy above

// GOOD
char name[32];
scanf("%31s", name);        // explicit width, leaving room for the '\0'
```

### 46.4 Format-string injection

```c
// BAD
void log_bad(const char *user_input) {
    printf(user_input);     // if user_input contains "%s"/"%n", this reads or WRITES
}                             // arbitrary memory through printf's own varargs machinery

// GOOD
void log_good(const char *user_input) {
    printf("%s", user_input);   // user_input is always DATA, never interpreted as a format string
}
```

Modern compilers already catch this one for you: Clang and GCC both emit
`-Wformat-security` (part of `-Wall`) — `printf(user_input)` triggers
*"format string is not a string literal (potentially insecure)"* at
compile time, which is exactly why checklist rule 11 says turn that
warning into a hard build failure with `-Werror`.

### 46.5 Shelling out with unsanitized input

```c
// BAD
char cmd[256];
snprintf(cmd, sizeof cmd, "grep %s /var/log/app.log", user_pattern);
system(cmd);   // user_pattern = "foo; rm -rf ~" runs a second, attacker-chosen command

// GOOD — avoid the shell entirely; call the tool directly with fork/exec (Section 26)
pid_t pid = fork();
if (pid == 0) {
    execlp("grep", "grep", user_pattern, "/var/log/app.log", (char *)NULL);
    _exit(127);
}
// (parent: waitpid as in Section 26)
```

`execlp`/`execvp` pass each argument as a separate, non-reinterpreted
string — there is no shell in between to reinterpret `;`, `|`, `` ` ``,
or `$()` — which is the entire class of "shell injection" vulnerability
`system()`/`popen()` are exposed to whenever any part of the command
line comes from outside the program.

### 46.6 Off-by-one and unsigned-underflow loop bounds

```c
// BAD
void print_reversed_bad(const char *s) {
    size_t len = strlen(s);
    for (size_t i = len - 1; i >= 0; i--) {   // size_t is UNSIGNED: when s is "", len-1
        putchar(s[i]);                          // wraps to SIZE_MAX, and `i >= 0` is ALWAYS true
    }                                            // — this reads far out of bounds, forever
}

// GOOD
void print_reversed_good(const char *s) {
    size_t len = strlen(s);
    for (size_t i = len; i > 0; i--) {
        putchar(s[i - 1]);
    }
}
```

### 46.7 Returning a pointer to a local variable

```c
// BAD
int *make_answer_bad(void) {
    int answer = 42;
    return &answer;      // answer's stack frame is gone the instant this function returns
}

// GOOD — either return by value, or have the caller own the storage
int make_answer_good(void) { return 42; }

void make_answer_out_param(int *out) { *out = 42; }   // when you need more than one value back
```

### 46.8 `sizeof` on a decayed array parameter

```c
// BAD
void print_len_bad(int arr[]) {
    printf("%zu\n", sizeof(arr));   // arr decayed to int* the moment it crossed the function
}                                     // boundary (Section 6.1) — this prints 8, not the array's length

// GOOD — pass the length explicitly; there is no other reliable way in C
void print_len_good(int *arr, size_t len) {
    printf("%zu\n", len * sizeof(*arr));
}
```

Clang and GCC flag the bad version directly: `-Wsizeof-array-argument`
(again, part of `-Wall`) — *"sizeof on array function parameter will
return size of 'int *' instead of 'int[]'."* Another entry for
checklist rule 11's file.

### 46.9 Comparing floating-point numbers with `==`

```c
// BAD
int equal_bad(double a, double b) {
    return a == b;    // 0.1 + 0.2 == 0.3 is FALSE — binary floating point can't represent
}                        // most decimal fractions exactly, so rounding error accumulates

// GOOD
#include <math.h>
int equal_good(double a, double b) {
    return fabs(a - b) < 1e-9;    // compare within a tolerance appropriate to your value range
}
```

### 46.10 Macros without full parenthesization, or with side-effecting arguments

```c
// BAD
#define SQUARE(x) x * x
int bad = SQUARE(1 + 2);        // expands to 1 + 2 * 1 + 2 = 5, not 9

// BAD (still, even parenthesized)
#define SQUARE2(x) ((x) * (x))
int i = 2;
int also_bad = SQUARE2(i++);    // expands to ((i++) * (i++)) — UB: i modified twice with no
                                  // sequence point between the two modifications (Section 8.4, 28)

// GOOD — a real function; the compiler still inlines it, and it evaluates i++ exactly once
static inline int square(int x) { return x * x; }
```

### 46.11 Double free and use-after-free from an un-nulled pointer

```c
// BAD
free(conn->socket_buf);
// ... fifty lines later, in an error path that re-checks conn->socket_buf ...
if (conn->socket_buf) free(conn->socket_buf);   // conn->socket_buf still holds the OLD address —
                                                   // "if not NULL" doesn't protect against double free

// GOOD
free(conn->socket_buf);
conn->socket_buf = NULL;    // now every later `if (conn->socket_buf)` check is actually meaningful
```

### 46.12 Feeding `malloc` an unchecked multiplication

```c
// BAD
int *matrix = malloc(rows * cols * sizeof(int));   // rows*cols can overflow size_t for
                                                      // attacker-influenced dimensions, wrapping
                                                      // to a tiny allocation that later writes overflow

// GOOD
int *matrix = calloc((size_t)rows * (size_t)cols, sizeof(int));  // calloc checks the
                                                                    // multiplication for overflow internally
```

### 46.13 Mutable global state that quietly breaks under threads or re-entrancy

```c
// BAD — works fine single-threaded, corrupts data the moment two threads call it,
// with no compiler warning to point at the problem
static char scratch[256];
char *format_message_bad(const char *name) {
    snprintf(scratch, sizeof scratch, "hello, %s!", name);
    return scratch;    // every caller gets a pointer to the SAME shared buffer
}

// GOOD — caller owns distinct storage; no shared mutable state to race on
void format_message_good(const char *name, char *out, size_t out_cap) {
    snprintf(out, out_cap, "hello, %s!", name);
}
```

This is precisely the bug class Section 20.1 demonstrated with a shared
counter, generalized: any file-scope or `static` mutable object that a
function reads or writes is a potential data race the moment that
function is called from more than one thread, or re-enters itself
(a signal handler calling back into code already running).

### 46.14 Skipping the warnings instead of fixing them

```bash
# BAD — silences the exact class of bug this whole section demonstrates
clang -w -o app app.c

# GOOD — the one flag change that would have caught 46.4, 46.8, and more, for free
clang -std=gnu17 -Wall -Wextra -Wpedantic -Werror -o app app.c
```

If there's one takeaway from this entire catalog, it's this last one:
sections 46.4 and 46.8 above are not hypothetical dangers you have to
remember to look for by eye — they're warnings your compiler already
emits by default under `-Wall`. The single highest-leverage habit in
this whole guide is never disabling that.

---

# Part IX — Hardware & driver programming

> Everything through Part VIII runs on top of an OS that already manages
> memory, files, and processes for you. This Part is what changes when
> your C code has to talk directly to a physical chip — a sensor, an
> LED, a USB device — either from a normal process, or from inside the
> kernel itself. This is the layer most job postings mean by "systems"
> or "embedded" experience, and where a senior engineer is expected to
> be able to go from a chip's datasheet to working code.
>
> **A note on testing.** Every code sample elsewhere in this guide was
> compiled and run before being written down. The `/dev/mem`-mmap and
> serial/termios examples below are ordinary POSIX C and were verified
> the same way. The Linux-kernel-specific examples (kernel modules,
> `ioctl` calls against `<linux/i2c-dev.h>`/`<linux/spi/spidev.h>`)
> cannot be compiled on the machine this guide was written on — a
> kernel module builds against your *running* kernel's exact header
> tree (`/lib/modules/$(uname -r)/build`), not a generic C toolchain,
> and `linux/*.h` headers don't exist outside Linux. What's shown is
> the standard, extremely stable API every Linux driver tutorial and
> textbook has used for well over a decade; treat it as accurate
> boilerplate to adapt, and build/test it on a real Linux box, a
> Raspberry Pi, or a QEMU VM before trusting it on real hardware.

## 47. Two worlds: user-space hardware I/O vs. writing a kernel driver

The phrase "write a driver" gets used for two very different jobs, and
conflating them is the single most common confusion senior candidates
run into in interviews for embedded/systems roles.

| | User-space hardware I/O | Kernel driver |
|---|---|---|
| Runs as | A normal process | Code loaded into the kernel itself |
| Talks to hardware via | A `/dev` node the kernel already provides (`/dev/i2c-1`, `/dev/spidev0.0`, `/dev/ttyUSB0`, `/dev/gpiochip0`) | Directly: memory-mapped registers, interrupts, DMA |
| Crash blast radius | That one process dies; the OS is fine | Can panic the entire machine — there's no memory protection between your bug and the kernel |
| When you need it | The overwhelming majority of real "hardware" work: reading a sensor, driving a motor controller, talking to a USB device, flashing firmware | The kernel doesn't already expose the hardware as a `/dev` node, or you need interrupt-tight integration with a kernel subsystem (a new bus type, a new filesystem, a new network device) |
| Debugging | gdb/lldb, exactly like any other program in this guide | `printk`+`dmesg`, kernel oops messages, `ftrace`, sometimes a JTAG debugger or QEMU |

**Real-world example.** A team building a robotics product needs to read
an IMU (inertial measurement unit) over I2C and drive motor PWM signals.
Neither needs a kernel driver — the SoC vendor's Linux distribution
already exposes I2C and PWM as `/dev` nodes and sysfs files; the actual
"driver" work is a user-space C program doing `ioctl()` calls against
those nodes (Section 48). The team only reaches for actual kernel code
when they add a custom FPGA-based sensor board with no existing driver
and genuinely new interrupt-driven behavior the kernel needs to know
about (Sections 50–53). This 90/10 split — mostly user-space I/O, rarely
a real kernel module — holds for almost every "hardware" job description
you'll see; this Part covers both halves so you can tell, for a given
task, which one you actually need.

---

## 48. Talking to hardware from user space: GPIO, I2C, SPI, and serial

Linux exposes most hardware buses as ordinary files. This is why a huge
amount of "driver" work is really just Part II's file I/O (Section 12)
plus `ioctl()`, a system call that's `read`/`write`'s escape hatch for
"do something device-specific that doesn't fit read/write."

### 48.1 GPIO: turning a pin into a `/dev` handle

Modern Linux (kernel 4.8+) exposes GPIO through a character device per
chip (`/dev/gpiochip0`) rather than the older, now-deprecated
`/sys/class/gpio` sysfs interface. The standard way to drive it from C
is `libgpiod`, which wraps the `ioctl` protocol for you:

```c
// requires: apt install libgpiod-dev ; link with -lgpiod
#include <gpiod.h>
#include <stdio.h>

int main(void) {
    struct gpiod_chip *chip = gpiod_chip_open("/dev/gpiochip0");
    if (!chip) { perror("gpiod_chip_open"); return 1; }

    struct gpiod_line *led = gpiod_chip_get_line(chip, 17);  // BCM pin 17
    if (!led) { perror("gpiod_chip_get_line"); gpiod_chip_close(chip); return 1; }

    if (gpiod_line_request_output(led, "blink-example", 0) < 0) {
        perror("gpiod_line_request_output");
        gpiod_chip_close(chip);
        return 1;
    }

    gpiod_line_set_value(led, 1);   // turn the LED on
    gpiod_line_release(led);
    gpiod_chip_close(chip);
    return 0;
}
```

The 0/1 you pass to `gpiod_line_set_value` ends up, several layers down,
as exactly the kind of register write Section 49 shows you doing by
hand — `libgpiod` exists so that in the 90% case (Section 47's table)
you never have to compute a register offset yourself.

### 48.2 I2C: addressed devices on a shared two-wire bus

I2C devices share one bus, each answering to its own 7-bit address. From
user space, you open the bus, tell the kernel which address you want to
talk to, then `read`/`write` like a normal file:

```c
#include <stdio.h>
#include <stdint.h>
#include <fcntl.h>
#include <unistd.h>
#include <linux/i2c-dev.h>   // Linux-only header — see this Part's opening note
#include <sys/ioctl.h>

#define LM75_ADDR 0x48   // a common temperature-sensor address

int read_temperature(const char *bus_path) {
    int fd = open(bus_path, O_RDWR);
    if (fd < 0) { perror("open"); return -1; }

    if (ioctl(fd, I2C_SLAVE, LM75_ADDR) < 0) {   // "all reads/writes on this fd target 0x48"
        perror("ioctl(I2C_SLAVE)");
        close(fd);
        return -1;
    }

    uint8_t reg = 0x00;                 // the LM75's temperature register
    if (write(fd, &reg, 1) != 1) { perror("write"); close(fd); return -1; }

    uint8_t raw[2];
    if (read(fd, raw, 2) != 2) { perror("read"); close(fd); return -1; }
    close(fd);

    int16_t temp_raw = (int16_t)((raw[0] << 8) | raw[1]) >> 7; // datasheet-specific bit layout
    return temp_raw / 2;                // degrees Celsius, per this sensor's datasheet
}

int main(void) {
    int c = read_temperature("/dev/i2c-1");
    if (c != -1) printf("%d C\n", c);
    return 0;
}
```

Many simple Linux userspace I2C tools follow this shape: open the bus
file, `ioctl(I2C_SLAVE, address)` to select the chip, then `write()` the
register you want followed by `read()` of however many bytes the
datasheet says that register holds. Kernel drivers and SMBus helpers wrap
the same idea in different APIs. The bit-shifting to turn two raw bytes
into a signed temperature (`temp_raw >> 7`, `/ 2`) is datasheet-specific
— Section 8.2's union-based byte inspection and Section 34's `sds`-style
"read the bytes, interpret per the format" mindset are exactly the skills
this line exercises.

### 48.3 SPI: full-duplex, chip-select, and one `ioctl` per transfer

SPI has no addressing — instead, a separate chip-select line picks which
device is listening, and every transfer is full-duplex (you always both
send and receive the same number of bytes, even if you only care about
one direction):

```c
#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include <linux/spi/spidev.h>   // Linux-only header
#include <sys/ioctl.h>

int spi_transfer(int fd, const uint8_t *tx, uint8_t *rx, size_t len) {
    struct spi_ioc_transfer xfer = {
        .tx_buf = (unsigned long)tx,
        .rx_buf = (unsigned long)rx,
        .len    = (uint32_t)len,
        .speed_hz = 1000000,   // 1 MHz — must not exceed the chip's datasheet max
        .bits_per_word = 8,
    };
    return ioctl(fd, SPI_IOC_MESSAGE(1), &xfer);   // one call = one full-duplex transfer
}

int main(void) {
    int fd = open("/dev/spidev0.0", O_RDWR);
    if (fd < 0) { perror("open"); return 1; }

    uint8_t tx[2] = {0x9F, 0x00};   // e.g. a flash chip's "read JEDEC ID" command
    uint8_t rx[2] = {0};
    if (spi_transfer(fd, tx, rx, sizeof tx) < 0) { perror("ioctl"); close(fd); return 1; }

    printf("received: %02x %02x\n", rx[0], rx[1]);
    close(fd);
    return 0;
}
```

### 48.4 Serial/UART: extending Section 25's `termios` beyond a terminal

Section 25 used `termios` to turn your *terminal* into raw, byte-at-a-time
mode. The exact same API configures a *real* serial port talking to
external hardware (a GPS module, an Arduino, an RS-232 sensor) — the
only difference is you now also set a baud rate and enable the receiver:

```c
#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>
#include <termios.h>

int open_serial_port(const char *path) {
    int fd = open(path, O_RDWR | O_NOCTTY | O_SYNC);
    if (fd < 0) { perror("open"); return -1; }

    struct termios tty;
    if (tcgetattr(fd, &tty) != 0) { perror("tcgetattr"); close(fd); return -1; }

    cfsetispeed(&tty, B9600);
    cfsetospeed(&tty, B9600);

    tty.c_cflag |= (CLOCAL | CREAD);   // ignore modem control lines, enable the receiver —
    tty.c_cflag &= ~PARENB;             // without CREAD the port simply never delivers bytes
    tty.c_cflag &= ~CSTOPB;
    tty.c_cflag &= ~CSIZE;
    tty.c_cflag |= CS8;                  // 8N1: 8 data bits, no parity, 1 stop bit

    tty.c_lflag &= ~(ICANON | ECHO | ECHOE | ISIG);  // raw mode, same as Section 25
    tty.c_iflag &= ~(IXON | IXOFF | IXANY);
    tty.c_oflag &= ~OPOST;

    tty.c_cc[VMIN]  = 1;    // block until at least 1 byte arrives
    tty.c_cc[VTIME] = 0;

    if (tcsetattr(fd, TCSANOW, &tty) != 0) { perror("tcsetattr"); close(fd); return -1; }
    return fd;
}
```

**War story.** Forgetting `CREAD` (or leaving `CLOCAL` unset on a port
whose carrier-detect line isn't wired up) is *the* classic "my serial
code compiles fine, connects fine, and then just never receives a
byte" bug — nothing errors, `read()` just blocks forever, because the
port's receiver was never actually enabled at the hardware level. It's
the serial-programming equivalent of Section 46.4's silent
compile-but-wrong bugs: syntactically correct, semantically incomplete.

---

## 49. Memory-mapped registers: `volatile`, `/dev/mem`, and memory barriers

### 49.1 What `volatile` actually promises — and doesn't

`volatile` tells the compiler: *"this memory can change for reasons you
can't see, and every access in the source must become a real load or
store — never cache it in a register, never eliminate it as
'redundant,' never reorder it relative to other volatile accesses."*
That's the core compiler-side guarantee; it is not a complete hardware
ordering model. It does **not** make an access atomic, does **not**
insert a CPU/device memory barrier, and does **not** synchronize with
another thread — for those you still need Section 38's atomics, OS or
architecture barrier primitives, or Section 52's kernel locks. Confusing
`volatile` with thread-safety is one of the most persistent myths in C;
`volatile` is mainly for memory whose value can change outside the
compiler's model of your program — a hardware register being the
textbook example.

```c
// Without volatile, the compiler is free to read `*status` ONCE, cache it
// in a register, and spin forever checking the cached copy — even though
// the hardware is actively changing the real memory underneath it.
while ((*status & READY_BIT) == 0) { }        // BUG: may never see the bit flip

while ((*(volatile uint32_t *)status & READY_BIT) == 0) { }   // correct: reloads every time
```

### 49.2 Mapping a hardware register block from user space

On Linux, `/dev/mem` gives a (root-only) process access to physical
memory — including a peripheral's registers — through the same `mmap`
you'd use for a file (Section 39's allocator used anonymous `mmap`; this
is the same call, pointed at physical hardware addresses instead):

```c
#include <stdio.h>
#include <stdint.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/mman.h>

#define BCM2835_PERI_BASE 0x3F000000              // Raspberry Pi 3's peripheral base address —
#define GPIO_BASE         (BCM2835_PERI_BASE + 0x200000)  // straight out of the SoC datasheet
#define BLOCK_SIZE        4096

static volatile uint32_t *gpio_map(void) {
    int fd = open("/dev/mem", O_RDWR | O_SYNC);
    if (fd < 0) { perror("open /dev/mem"); return NULL; }

    void *map = mmap(NULL, BLOCK_SIZE, PROT_READ | PROT_WRITE, MAP_SHARED, fd, GPIO_BASE);
    close(fd);   // safe once mapped — the mapping itself keeps the pages accessible

    if (map == MAP_FAILED) { perror("mmap"); return NULL; }
    return (volatile uint32_t *)map;
}

#define GPFSEL1 1   // word offset: function-select register for GPIOs 10-19
#define GPSET0  7   // word offset: writing a 1 here drives that GPIO high
#define GPCLR0  10  // word offset: writing a 1 here drives that GPIO low

void gpio_set_output(volatile uint32_t *gpio, int pin) {
    int reg = pin / 10, shift = (pin % 10) * 3;
    gpio[reg] &= ~(7u << shift);   // clear this pin's 3-bit function-select field
    gpio[reg] |= (1u << shift);     // 001 = output, per the BCM2835 datasheet's register table
}

void gpio_write(volatile uint32_t *gpio, int pin, int value) {
    if (value) gpio[GPSET0 + pin / 32] = 1u << (pin % 32);
    else       gpio[GPCLR0 + pin / 32] = 1u << (pin % 32);
}

int main(void) {
    volatile uint32_t *gpio = gpio_map();   // must run as root: sudo ./a.out
    if (!gpio) return 1;
    gpio_set_output(gpio, 17);
    gpio_write(gpio, 17, 1);   // turn an LED on
    return 0;
}
```

Every field here — `0x3F000000`, the `GPFSEL`/`GPSET`/`GPCLR` offsets,
the 3-bits-per-pin layout — comes directly from the BCM2835/BCM2711
peripheral datasheet, not from any general C knowledge. **This is the
actual job**: reading a register map out of a datasheet and turning each
row into a `#define` and a masked read-modify-write, exactly like
Section 17's bit-manipulation primitives, aimed at real silicon instead
of an in-memory flags variable. (In a real kernel driver you'd use
`ioremap()` instead of `/dev/mem` + `mmap` — same idea, Section 51.)

### 49.3 Memory barriers: when `volatile` alone isn't enough

`volatile` stops the *compiler* eliding accesses and constrains ordering
between volatile accesses in the abstract C program, but it is not a
portable device-memory ordering primitive. A modern CPU, interconnect, or
posted-write path can still delay when a write becomes visible to a
device. When the *order* of two register writes matters (very common:
"set the configuration register, *then* the enable register, in that
exact order, or the device ignores the config"), use the barrier primitive
provided by your OS, compiler, or architecture:

```c
#include <stdatomic.h>   // useful for C atomics; real MMIO often needs OS/arch barriers

gpio[CONFIG_REG] = desired_config;
atomic_thread_fence(memory_order_release);   // teaching sketch; not a universal MMIO barrier
gpio[ENABLE_REG] = 1;
```

Inside the Linux kernel you'd use the kernel's own `wmb()`/`rmb()`/`mb()`
or MMIO helpers such as `writel()`/`readl()` instead of `<stdatomic.h>`,
because they are defined for the kernel's CPU and device-memory model.
The intuition is similar to Section 38's acquire/release pairing, but
the C standard alone does not define how device registers observe memory.

**War story.** A board bring-up engineer wrote register-configuration
code that worked perfectly in a debug build and intermittently failed in
the optimized release build. The missing piece was not "make everything
volatile"; it was the device's required ordering between "set config" and
"set enable." Adding the platform's write barrier/MMIO helper made that
ordering explicit — the kind of small-looking change that can take days
to discover on real hardware.

---

## 50. Your first Linux kernel module: hello world, the Makefile, insmod/dmesg

A kernel module is code the kernel loads into its own address space at
runtime — no `main()`, no libc, no user-space memory protection between
your bug and a full system crash. Everything here needs a Linux machine
(a real one, a VM, or Raspberry Pi) with kernel headers installed for
your *exact running kernel* (`sudo apt install linux-headers-$(uname -r)`
on Debian/Ubuntu).

```c
// hello.c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>

static int __init hello_init(void) {
    printk(KERN_INFO "hello: module loaded\n");   // printk, not printf — there is no stdout
    return 0;                                        // 0 = success, matching every syscall convention
}

static void __exit hello_exit(void) {
    printk(KERN_INFO "hello: module unloaded\n");
}

module_init(hello_init);   // register these as the load/unload entry points
module_exit(hello_exit);

MODULE_LICENSE("GPL");                 // non-GPL modules can't call GPL-only kernel symbols
MODULE_AUTHOR("you");
MODULE_DESCRIPTION("A minimal Linux kernel module");
```

```makefile
# Makefile — a "Kbuild" makefile, not the CMake/Make from Section 41
obj-m += hello.o

all:
	$(MAKE) -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	$(MAKE) -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
```

```bash
make
sudo insmod hello.ko
dmesg | tail -1        # hello: module loaded
lsmod | grep hello      # confirms it's resident
sudo rmmod hello
dmesg | tail -1        # hello: module unloaded
```

**Why the build looks nothing like Section 41's CMake setup:** a kernel
module isn't linked against glibc or any user-space ABI at all — it's
compiled against, and only guaranteed to load into, the *exact* kernel
tree at `/lib/modules/$(uname -r)/build`. This is the sharp edge behind
Section 40's ABI-stability discussion: user-space libraries go to great
lengths (opaque pointers, symbol versioning) to stay binary-compatible
across versions; the Linux kernel's internal APIs deliberately make
**no such promise** between kernel versions — a module built for 6.6
is not expected to load into 6.8 without a rebuild. That's a conscious
trade-off (it lets kernel internals evolve freely) that every driver
maintainer has to live with.

---

## 51. Character device drivers: `file_operations`, `cdev`, and `copy_to/from_user`

A character device driver is what turns `module_init` from Section 50
into something a user-space program can `open()`/`read()`/`write()` as
`/dev/something` — the kernel-side implementation of every `/dev` node
Section 48's examples opened from the outside.

```c
// chardev.c — a device that remembers the last message written to it
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/fs.h>
#include <linux/cdev.h>
#include <linux/uaccess.h>   // copy_to_user / copy_from_user
#include <linux/device.h>
#include <linux/version.h>   // LINUX_VERSION_CODE / KERNEL_VERSION, for the class_create() shim below

#define DEVICE_NAME "mychardev"
#define BUF_SIZE    128

static char    message[BUF_SIZE];
static size_t  message_len;
static dev_t   dev_num;
static struct cdev my_cdev;
static struct class *my_class;

static int dev_open(struct inode *inode, struct file *file) {
    (void)inode; (void)file;
    return 0;    // nothing to set up for this simple device
}

static ssize_t dev_read(struct file *file, char __user *buf, size_t len, loff_t *offset) {
    (void)file;
    if (*offset >= message_len) return 0;              // EOF
    size_t remaining = message_len - (size_t)*offset;
    size_t to_copy = len < remaining ? len : remaining;

    // copy_to_user validates the user pointer AND does the copy — a raw memcpy would
    // let a malicious/buggy caller pass a kernel address and read arbitrary kernel memory
    if (copy_to_user(buf, message + *offset, to_copy)) return -EFAULT;

    *offset += (loff_t)to_copy;
    return (ssize_t)to_copy;
}

static ssize_t dev_write(struct file *file, const char __user *buf, size_t len, loff_t *offset) {
    (void)file; (void)offset;
    size_t to_copy = len < BUF_SIZE - 1 ? len : BUF_SIZE - 1;   // never overflow `message`

    if (copy_from_user(message, buf, to_copy)) return -EFAULT;
    message[to_copy] = '\0';
    message_len = to_copy;
    return (ssize_t)to_copy;
}

static const struct file_operations fops = {
    .owner = THIS_MODULE,
    .open  = dev_open,
    .read  = dev_read,
    .write = dev_write,
};

static int __init chardev_init(void) {
    if (alloc_chrdev_region(&dev_num, 0, 1, DEVICE_NAME) < 0) return -1;

    cdev_init(&my_cdev, &fops);
    if (cdev_add(&my_cdev, dev_num, 1) < 0) {
        unregister_chrdev_region(dev_num, 1);
        return -1;
    }

#if LINUX_VERSION_CODE < KERNEL_VERSION(6, 4, 0)
    my_class = class_create(THIS_MODULE, DEVICE_NAME);   // two-arg form: kernels before 6.4
#else
    my_class = class_create(DEVICE_NAME);                 // 6.4+ dropped the owner argument
#endif
    device_create(my_class, NULL, dev_num, NULL, DEVICE_NAME);  // makes /dev/mychardev appear

    printk(KERN_INFO "chardev: /dev/%s ready\n", DEVICE_NAME);
    return 0;
}

static void __exit chardev_exit(void) {
    device_destroy(my_class, dev_num);
    class_destroy(my_class);
    cdev_del(&my_cdev);
    unregister_chrdev_region(dev_num, 1);
}

module_init(chardev_init);
module_exit(chardev_exit);
MODULE_LICENSE("GPL");
```

```bash
make && sudo insmod chardev.ko
echo "hello kernel" > /dev/mychardev
cat /dev/mychardev          # hello kernel
sudo rmmod chardev
```

**Why `copy_to_user`/`copy_from_user` instead of `memcpy` (Section 10's
usual tool):** kernel and user space are separate address spaces with
separate page tables. A user-space pointer is only meaningful in the
context of *that process's* mappings — `memcpy`ing through it directly
from kernel code would either fault or, worse, silently read/write
whatever the raw address happens to resolve to in kernel space.
`copy_to_user`/`copy_from_user` validate that the user address is
actually mapped and accessible to that process *before* touching it,
returning an error instead of corrupting memory — this single pair of
functions is the kernel-space equivalent of Section 29's entire "trust
nothing that crosses a boundary" chapter, with the boundary being the
user/kernel split instead of the network.

---

## 52. Interrupts and kernel concurrency: top/bottom halves, spinlocks vs. mutexes

### 52.1 Registering an interrupt handler

```c
#include <linux/interrupt.h>
#include <linux/gpio.h>

static irqreturn_t button_isr(int irq, void *dev_id) {
    (void)irq; (void)dev_id;
    printk(KERN_INFO "button: pressed\n");
    return IRQ_HANDLED;    // tells the kernel this handler claimed the interrupt
}

static int irq_number;

static int __init button_init(void) {
    irq_number = gpio_to_irq(17);   // the IRQ line wired to GPIO 17
    return request_irq(irq_number, button_isr, IRQF_TRIGGER_RISING, "button_handler", NULL);
}

static void __exit button_exit(void) {
    free_irq(irq_number, NULL);
}
```

### 52.2 The rule that changes everything: an interrupt handler cannot sleep

Code running in **interrupt context** (an ISR, or the top half of one)
runs with normal scheduling paused on that CPU — it cannot call anything
that might block: no `mutex_lock`, no `kmalloc(..., GFP_KERNEL)`
(use `GFP_ATOMIC` instead), no blocking I/O. Calling a sleeping function
from interrupt context doesn't error gracefully — it produces a kernel
BUG (`"scheduling while atomic"` or `"sleeping function called from
invalid context"`) and can wedge or crash the whole machine. This is why
real interrupt handling is split in two:

- **Top half** (the ISR itself, like `button_isr` above): do the
  absolute minimum — acknowledge the hardware, grab the data that would
  otherwise be lost, and return, as fast as possible.
- **Bottom half** (a tasklet or a workqueue): the actual work, deferred
  to run later in a context that's allowed to sleep.

```c
#include <linux/workqueue.h>

static struct work_struct button_work;

static void button_work_fn(struct work_struct *work) {
    (void)work;
    // safe here to sleep, take a mutex, allocate with GFP_KERNEL, etc. —
    // this runs in process context, not interrupt context
    printk(KERN_INFO "button: doing the slow part\n");
}

static irqreturn_t button_isr(int irq, void *dev_id) {
    (void)irq; (void)dev_id;
    schedule_work(&button_work);   // defer the real work — return from the ISR immediately
    return IRQ_HANDLED;
}

// in the module's init function: INIT_WORK(&button_work, button_work_fn);
```

### 52.3 Spinlocks vs. mutexes: the rule that governs which lock to use

Section 20's `pthread_mutex_t` assumed something the kernel can't always
assume: that the thread holding a lock is free to be put to sleep while
another thread waits. That's false the moment an interrupt handler needs
the same data a normal kernel thread does — the ISR *cannot* sleep, so
it cannot wait on a mutex.

| | Can sleep while waiting? | Use when |
|---|---|---|
| `mutex_lock()` | Yes | Data only ever touched from process context (never from an ISR) |
| `spin_lock()` / `spin_lock_irqsave()` | No — busy-waits | Data shared with an interrupt handler; `_irqsave` variant also disables local interrupts, since the ISR could otherwise preempt the very code holding the lock, on the same CPU |
| `atomic_t` / `atomic_inc()` etc. | N/A — single indivisible op | A simple counter or flag, no larger critical section needed |

```c
#include <linux/spinlock.h>

static DEFINE_SPINLOCK(button_lock);
static int press_count;

static irqreturn_t button_isr(int irq, void *dev_id) {
    (void)irq; (void)dev_id;
    unsigned long flags;
    spin_lock_irqsave(&button_lock, flags);   // disable interrupts on this CPU + acquire the lock
    press_count++;
    spin_unlock_irqrestore(&button_lock, flags);
    return IRQ_HANDLED;
}
```

**War story.** "Sleeping in atomic context" — calling a mutex or a
blocking allocation from inside a spinlock-held region or an interrupt
handler — is one of the most common kernel bugs new driver authors
write, and it's exactly analogous to Section 46.13's "mutable global
state that quietly breaks under threads": it often works fine in
testing (the sleep path is rarely hit) and then wedges a production
machine the first time the unlucky timing actually occurs.

---

## 53. Debugging drivers: `printk`/dmesg, an oops backtrace, and QEMU+gdb

### 53.1 `printk` log levels and `dmesg`

```c
printk(KERN_EMERG   "system is unusable\n");   // level 0 — highest
printk(KERN_ALERT   "action must be taken immediately\n");
printk(KERN_CRIT    "critical condition\n");
printk(KERN_ERR     "error condition\n");
printk(KERN_WARNING "warning condition\n");
printk(KERN_NOTICE  "normal but significant\n");
printk(KERN_INFO    "informational\n");         // what Sections 50-52 used throughout
printk(KERN_DEBUG   "debug-level message\n");  // level 7 — lowest
```

```bash
dmesg | tail -20            # see recent kernel log messages, including every printk above
dmesg -w                    # follow, like `tail -f`, while you insmod/trigger your driver
cat /proc/sys/kernel/printk # which levels also go to the console, not just the ring buffer
```

### 53.2 Reading a kernel oops

A bug in a kernel module (a NULL dereference, an out-of-bounds access —
Section 28's UB list, now with no memory-protection safety net at all)
produces an **oops**: a backtrace printed to `dmesg`, naming the faulting
function and an instruction-pointer address, then usually killing only
the offending process (a full **panic** — Section 30's hardening context
in user space has no real equivalent here — halts everything).

```
BUG: kernel NULL pointer dereference, address: 0000000000000008
...
RIP: 0010:chardev_read+0x12/0x40 [chardev]
Call Trace:
 vfs_read+0x9e/0x1a0
 ksys_read+0x67/0xe0
 do_syscall_64+0x38/0x90
```

`chardev_read+0x12/0x40 [chardev]` reads exactly like Section 14.4's
`bt` output in gdb, just for kernel code: the faulting function, an
offset into it, and which module it came from. `addr2line` (or the
kernel source tree's `scripts/faddr2line`) turns that offset back into
a source file and line number — the same "map an address back to a
source location" workflow as any other crash in this guide, just with
`dmesg` standing in for gdb's backtrace.

### 53.3 Iterating without risking real hardware: QEMU + gdb

You do not need to `insmod` untested driver code onto real hardware to
debug it — booting a kernel under QEMU with gdb attached gives you a
disposable "machine" you can crash as many times as you need:

```bash
qemu-system-x86_64 -kernel bzImage -initrd initrd.img \
    -append "console=ttyS0 nokaslr" -nographic -s -S
    # -s: open a gdb server on :1234    -S: pause the CPU before boot

# in another terminal:
gdb vmlinux
(gdb) target remote :1234
(gdb) break chardev_init
(gdb) continue
```

This is the same gdb from Section 14.4, attached over the network
instead of to a local process — the entire "set a breakpoint, continue,
inspect a variable" workflow you already know carries over unchanged;
only the target (a whole virtual machine's kernel, paused before it
even boots) is new.

---

## 54. Capstone: one sensor, driven two ways

Take the LM75-style I2C temperature sensor from Section 48.2 and build
it both ways from Section 47's table, so the trade-off is concrete
rather than theoretical.

**Version A — user space (what you'd actually ship for this task):**
exactly Section 48.2's `read_temperature`, wrapped in a small CLI
(Section 24) that prints a reading every second:

```c
#include <stdio.h>
#include <unistd.h>

int read_temperature(const char *bus_path);   // from Section 48.2

int main(void) {
    while (1) {
        int c = read_temperature("/dev/i2c-1");
        if (c != -1) printf("%d C\n", c);
        sleep(1);
    }
}
```

No kernel code, no reboot risk, ships as an ordinary binary — this is
correctly the default answer for "read a sensor" per Section 47.

**Version B — a minimal kernel driver**, for when the sensor's readings
need to feed a kernel subsystem directly (for example, the kernel's
`hwmon` framework, so `/sys/class/hwmon/.../temp1_input` and every
existing temperature-monitoring tool works automatically, with no
custom user-space program at all):

```c
// lm75_mini.c — sketch: an i2c_driver that reads the same register Section 48.2 did,
// but from kernel context, matched to the hardware via the i2c core instead of a
// hand-opened /dev/i2c-1 node
#include <linux/module.h>
#include <linux/i2c.h>

static int lm75_read_temp(struct i2c_client *client) {
    s32 raw = i2c_smbus_read_word_data(client, 0x00);   // same register as Section 48.2
    if (raw < 0) return raw;
    s16 be = (s16)((raw << 8) | (raw >> 8));              // the chip returns big-endian
    return (be >> 7) / 2;                                   // identical math to Section 48.2
}

static int lm75_probe(struct i2c_client *client) {
    int temp = lm75_read_temp(client);
    dev_info(&client->dev, "lm75_mini: %d C\n", temp);
    return 0;
}

static const struct i2c_device_id lm75_id[] = { { "lm75_mini", 0 }, { } };
MODULE_DEVICE_TABLE(i2c, lm75_id);

static struct i2c_driver lm75_driver = {
    .driver = { .name = "lm75_mini" },
    .probe  = lm75_probe,
    .id_table = lm75_id,
};
module_i2c_driver(lm75_driver);   // registers the driver; the i2c core calls .probe
                                    // automatically once a matching device is found
MODULE_LICENSE("GPL");
```

The line worth sitting with is `.probe = lm75_probe`: Version A decides
*when* to talk to the sensor (your `while` loop calls the shots).
Version B instead *registers interest* — `module_i2c_driver` hands
control to the kernel's i2c core, which calls `lm75_probe` on your
behalf whenever a device matching `lm75_id` (typically declared in the
board's Device Tree, outside the scope of this guide) actually shows
up. That inversion — you stop calling the hardware and start being
*called* by a framework when the hardware appears — is the core idea
behind every Linux driver model (i2c, SPI, USB, platform, PCI): each
framework differs in its match rules and callback set, but every one of
them replaces your `main()`-style control flow with a `probe`/`remove`
pair the kernel drives.

---

# Part X — Signals, non-local jumps, and memory-mapped files

> Three topics that didn't fit anywhere else in this guide because
> they're not "beginner," "advanced," or "hardware" specifically — they're
> general-purpose systems C that shows up in interviews, code reviews,
> and real production incidents regardless of what you're building.
> Unlike Part IX's kernel-space code, every example here is portable
> POSIX C, compiled and run on the same machine this whole guide was
> written on.

## 55. Signals: handling asynchronous events safely

### 55.1 What a signal actually is

A signal is the kernel interrupting your process's normal control flow
to say "something happened," asynchronously — a user pressed Ctrl+C
(`SIGINT`), a child process exited (`SIGCHLD`), your process wrote to a
closed pipe (`SIGPIPE`), or it dereferenced bad memory (`SIGSEGV`). By
default, most signals terminate the process (some also dump core); you
can install a handler to run your own code instead.

### 55.2 `sigaction`, not `signal` — and why

```c
#include <signal.h>
#include <stdio.h>
#include <unistd.h>

static volatile sig_atomic_t shutdown_requested = 0;

void handle_sigint(int sig) {
    (void)sig;
    shutdown_requested = 1;   // only touch a sig_atomic_t here — nothing else is safe (55.3)
}

int main(void) {
    struct sigaction sa = {0};
    sa.sa_handler = handle_sigint;
    sigemptyset(&sa.sa_mask);   // don't block any additional signals while this handler runs
    sa.sa_flags = 0;

    if (sigaction(SIGINT, &sa, NULL) != 0) { perror("sigaction"); return 1; }

    while (!shutdown_requested) {
        printf("working...\n");
        sleep(1);
    }
    printf("shutting down cleanly\n");
    return 0;
}
```

The older `signal()` function still exists but its behavior around
whether a handler stays installed after firing, and whether interrupted
syscalls auto-restart, historically differed between BSD and System V —
`sigaction` pins down every one of those choices explicitly
(`sa_flags`), which is why every serious codebase uses it instead. This
mirrors §2.2's advice to pin `-std=` explicitly rather than rely on
"whatever the default happens to be": `signal()`'s portability gaps are
exactly that kind of implicit, compiler/platform-dependent default.

### 55.3 Async-signal-safety: the rule that makes most of Part II illegal

A handler can run *at any point* — including in the middle of `malloc`
updating its internal free list, or `printf` mid-write to a
partially-flushed buffer. If the handler itself then calls `malloc` or
`printf`, it can deadlock against the very state the interrupted call
was in the middle of modifying. POSIX defines a short list of
**async-signal-safe** functions safe to call from a handler (`write`,
`_exit`, `sig*` functions, and a few others) — `malloc`, `free`,
`printf`, and most of the C standard library are **not** on that list.

```c
// BAD — printf is not async-signal-safe; can deadlock or corrupt output
void handle_sigint_bad(int sig) { (void)sig; printf("caught!\n"); }

// GOOD — write() to a raw fd is async-signal-safe; sig_atomic_t needs no library call at all
void handle_sigint_good(int sig) {
    (void)sig;
    const char msg[] = "caught!\n";
    write(STDOUT_FILENO, msg, sizeof msg - 1);
}
```

The pattern from 55.2 — set a `volatile sig_atomic_t` flag in the
handler, check it in your normal control flow — sidesteps the whole
problem: the handler does the absolute minimum, and all real work
happens back in ordinary, non-signal context. This is precisely Part
IX §52.2's "top half does the minimum, defer the real work" rule,
applied to signals instead of hardware interrupts — the same
constraint, the same fix, a different trigger.

### 55.4 Reaping children with `SIGCHLD`

```c
#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/wait.h>

void reap_children(int sig) {
    (void)sig;
    int saved_errno = errno;          // waitpid can clobber errno; the handler must restore it
    pid_t pid;
    while ((pid = waitpid(-1, NULL, WNOHANG)) > 0) {
        // a real handler would just record `pid` somewhere async-signal-safe;
        // it must NOT log/printf here (55.3)
    }
    (void)pid;
    errno = saved_errno;
}

int main(void) {
    struct sigaction sa = {0};
    sa.sa_handler = reap_children;
    sigemptyset(&sa.sa_mask);
    sa.sa_flags = SA_RESTART;   // let interrupted syscalls (e.g. a blocking read) auto-retry
    sigaction(SIGCHLD, &sa, NULL);

    for (int i = 0; i < 3; i++) {
        if (fork() == 0) { _exit(0); }   // each child exits immediately
    }
    sleep(1);   // give SIGCHLD handlers time to reap all three before main exits
    printf("parent done\n");
    return 0;
}
```

Without a `SIGCHLD` handler (or an explicit `waitpid` call), an exited
child becomes a **zombie** — gone from execution but still occupying a
process-table slot until its parent collects its exit status. Section
26's mini-shell called `waitpid` synchronously after every `fork`, which
works for a shell that blocks on each foreground command; a server that
spawns background workers needs exactly this asynchronous
`SIGCHLD`-driven reaping instead, since it can't afford to block waiting
for any one child.

### 55.5 War story, and a real technique to know (not to casually run)

Installing a handler for `SIGSEGV`/`SIGBUS` paired with `sigsetjmp`/
`siglongjmp` (Section 56.3) — turning a hardware fault into a
recoverable error instead of a crash — is a genuine production
technique: it's how some embedded-scripting-language sandboxes and
guard-page-based memory checkers turn an out-of-bounds access into a
catchable error rather than taking the whole process down. It's also
exactly the kind of code that behaves differently across environments in
ways that matter: on this guide's own development machine, deliberately
triggering a `SIGSEGV` inside a sandboxed shell caused the process to
hang rather than reach the handler — almost certainly the OS's own
crash-reporting/debugging hooks intercepting the fault first, ahead of
the process's own handler. The lesson isn't "don't do this" — it's
**test signal-based fault recovery on a plain terminal on real hardware,
never inside a CI runner or sandboxed environment**, because the exact
thing you're trying to observe (a raw hardware trap reaching your
handler) is precisely what crash-reporting tooling is also trying to
intercept first.

---

## 56. Non-local jumps: `setjmp`/`longjmp` and why they're not exceptions

### 56.1 The mechanism

`setjmp` saves the current stack context into a `jmp_buf`; `longjmp`
restores it, making execution jump back to exactly that `setjmp` call —
*as if it were returning again*, but with the value you pass to
`longjmp` instead of `0`. It's the only way in C to unwind multiple
stack frames at once without every intermediate function explicitly
checking and returning an error code.

```c
#include <setjmp.h>
#include <stdio.h>

static jmp_buf recovery_point;

void deeply_nested(int depth, int fail_at) {
    if (depth == fail_at) {
        longjmp(recovery_point, 1);   // jump straight back to setjmp, skipping every return
    }
    printf("entering depth %d\n", depth);
    deeply_nested(depth + 1, fail_at);
    printf("leaving depth %d\n", depth);   // never runs for frames at or past fail_at
}

int main(void) {
    if (setjmp(recovery_point) == 0) {
        deeply_nested(0, 3);
        printf("this line never runs\n");
    } else {
        printf("recovered: error occurred at depth 3\n");
    }
    return 0;
}
// entering depth 0 / entering depth 1 / entering depth 2 / recovered: error occurred at depth 3
// — notice NO "leaving depth" lines print at all; longjmp skipped every one of those returns.
```

### 56.2 The one rule that isn't optional: mark shared locals `volatile`

```c
#include <setjmp.h>
#include <stdio.h>

static jmp_buf recovery_point;

void risky(volatile int *counter) {
    (*counter)++;
    longjmp(recovery_point, 1);
}

int main(void) {
    volatile int counter = 0;   // MUST be volatile: written before longjmp, read after it —
                                  // without volatile, the compiler may keep counter in a
                                  // register that the stack-rewind doesn't restore, so the
                                  // increment could appear to vanish after the jump
    if (setjmp(recovery_point) == 0) {
        risky(&counter);
    } else {
        printf("counter after longjmp: %d\n", counter);   // reliably 1, because of volatile
    }
    return 0;
}
```

This is the *other* legitimate use of `volatile` alongside Part IX
§49.1's hardware registers: any local variable that's modified between
`setjmp` and a `longjmp` that returns to it, and then read afterward,
needs `volatile` — otherwise the C standard explicitly leaves its value
after the jump indeterminate.

### 56.3 What it doesn't do: no cleanup, no destructors, no stack unwinding of resources

```c
void process(void) {
    FILE *f = fopen("data.txt", "r");
    parse(f);        // if parse() longjmps out on a syntax error...
    fclose(f);        // ...this line never runs — the FILE* leaks, exactly like an
}                       // error-path leak (§44.1), just triggered by a jump instead of a return
```

C has no equivalent of C++ destructors or `defer` — a `longjmp` that
skips over a scope skips *everything* in it, including any pending
`free`/`fclose`/`pthread_mutex_unlock`. Anything a `longjmp` might jump
over needs its own `goto cleanup`-style handling (§4.2) reachable
*before* the jump, or must be tracked and released at the recovery
point itself. This is precisely why real interpreters that use
`setjmp`/`longjmp` for error recovery (many do, including reference
implementations of scripting languages) pair it with an explicit
resource-tracking stack that the recovery handler walks and releases —
`longjmp` itself will not do it for you.

### 56.4 `sigsetjmp`/`siglongjmp`: the signal-safe variant

Plain `setjmp`/`longjmp` don't reliably save/restore the process's
*signal mask* across the jump. If you're recovering from inside a signal
handler (Section 55.5's pattern), use `sigsetjmp(buf, 1)` /
`siglongjmp` instead — the `1` tells it to also save the current signal
mask, so `siglongjmp` correctly restores which signals were
blocked/unblocked at the time, instead of potentially leaving a signal
masked forever after the jump.

---

## 57. Memory-mapped files: `mmap` for real files, `msync`, and shared mappings

You've now seen `mmap` twice: Section 39 used it *anonymously* to grab
raw pages from the OS for a custom allocator; Part IX §49.2 used it on
`/dev/mem` to reach hardware registers. Its most common real-world use
is neither of those — mapping an actual file into memory, so reading it
is just reading memory, with no explicit `read()` loop at all.

### 57.1 Reading a file via `mmap` instead of Section 12.2's `fread` loop

```c
#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/mman.h>
#include <sys/stat.h>

int mmap_read_file(const char *path, const char **out_data, size_t *out_len) {
    *out_data = NULL;
    *out_len = 0;

    int fd = open(path, O_RDONLY);
    if (fd < 0) { perror("open"); return 0; }

    struct stat st;
    if (fstat(fd, &st) != 0) { perror("fstat"); close(fd); return 0; }
    if (st.st_size == 0) { close(fd); return 1; }   // empty file: success, no mapping

    void *map = mmap(NULL, (size_t)st.st_size, PROT_READ, MAP_PRIVATE, fd, 0);
    close(fd);   // safe to close once mapped — the mapping itself keeps the pages available
    if (map == MAP_FAILED) { perror("mmap"); return 0; }

    *out_data = map;
    *out_len = (size_t)st.st_size;
    return 1;     // *out_data is NOT malloc'd — release it with munmap(), never free()
}

int main(void) {
    size_t len = 0;
    const char *data = NULL;
    if (!mmap_read_file("/etc/hosts", &data, &len)) return 1;
    if (!data) { puts("empty file"); return 0; }

    printf("mapped %zu bytes, first line: ", len);
    for (size_t i = 0; i < len && data[i] != '\n'; i++) putchar(data[i]);
    putchar('\n');

    munmap((void *)data, len);
    return 0;
}
```

There's no `read()` call at all — the kernel demand-pages the file's
contents in as you touch them, and the OS's page cache backs the memory
directly, so re-reading the same region a second time is often cheap.
For a multi-gigabyte file you only need to scan once, this avoids both
Section 12.2's explicit read loop *and* the extra copy into a
user-space buffer that a `read()`-into-your-own-buffer approach performs.

### 57.2 Writing through a mapping, and why `msync` exists

```c
#include <stdio.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/mman.h>

int main(void) {
    const char *path = "/tmp/mmap_demo.txt";
    int fd = open(path, O_RDWR | O_CREAT | O_TRUNC, 0644);
    if (fd < 0) { perror("open"); return 1; }

    size_t len = 64;
    if (ftruncate(fd, (off_t)len) != 0) {   // must size the file before mapping it writable —
        perror("ftruncate"); close(fd); return 1;   // mmap doesn't grow the file for you
    }

    char *map = mmap(NULL, len, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0);
    close(fd);
    if (map == MAP_FAILED) { perror("mmap"); return 1; }

    strcpy(map, "written straight through the mapping, no write() call");
    msync(map, len, MS_SYNC);   // force this write to disk now, instead of waiting for
                                  // the kernel's own writeback timing to get around to it
    munmap(map, len);

    FILE *f = fopen(path, "r");   // prove it landed on disk, via an ordinary read this time
    char buf[64] = {0};
    fread(buf, 1, sizeof buf - 1, f);
    fclose(f);
    printf("read back: %s\n", buf);
    return 0;
}
```

`MAP_SHARED` (instead of `MAP_PRIVATE`) means writes go back to the
underlying file and are visible to any *other* process that maps the
same file `MAP_SHARED` — this is a legitimate, fast IPC mechanism
between unrelated processes with no pipe or socket involved at all,
each one just reading/writing the same mapped region.

### 57.3 The trap: accessing past what you mapped is `SIGBUS`, not a normal error

```c
char *map = mmap(NULL, 100, PROT_READ, MAP_PRIVATE, fd, 0);
char last_byte = map[99];   // fine, if the file is at least 100 bytes
char oops      = map[500];  // if the file is shorter than 500 bytes: SIGBUS, not a clean error —
                              // there's no bounds check here any more than there was in §29.1
```

Unlike a `read()` past EOF (which just returns 0 bytes, cleanly), an
out-of-bounds access on a mapped region is a hardware page fault the
kernel turns into `SIGBUS` — Section 55's signal-handling material is
exactly what you'd reach for to catch this defensively in, say, a
library that maps files of a size it doesn't fully control.

**Real-world example.** This is how SQLite's optional mmap I/O mode
works (Section 33), how `git` reads multi-gigabyte pack files without
loading them wholesale into a buffer, and — closing a loop back to
Section 40 — precisely the mechanism the dynamic linker (`ld.so`) uses
to map a shared library's code and data segments into your process at
program startup, before `main()` ever runs.

---

**That's the full arc** — from the toolchain in Part I to reading real
allocator, kernel-module, and mmap internals here in Part X. What follows is
reference material, not more reading: the C17→C23 feature diff, the UB
cheat sheet, compiler flags, a glossary, and further reading — check back
against these while you're actually writing C, rather than reading them
front to back.

---

# Appendices

## A. C17 → C23: what's new, feature by feature

```c
// nullptr — a real, type-safe null pointer constant (vs. NULL, which is
// typically (void*)0 or just 0 — ambiguous in variadic/overload-like contexts)
int *p = nullptr;

// bool/true/false are now KEYWORDS — no more #include <stdbool.h>
bool ok = true;

// constexpr — a compile-time constant with real type-checking (unlike #define)
constexpr int MAX_USERS = 100;
int users[MAX_USERS];

// typeof / typeof_unqual — write generic macros without _Generic boilerplate
#define SWAP(a, b) do { typeof(a) tmp = (a); (a) = (b); (b) = tmp; } while (0)

// [[attributes]] — standardized, compiler-checked hints
[[nodiscard]] int must_check(void);
[[deprecated("use new_func instead")]] void old_func(void);
[[maybe_unused]] static int debug_only_var;

// Binary literals and digit separators — readability wins, zero semantic change
unsigned int mask = 0b1010'1100;

// #embed — include a binary file's bytes as an initializer list, no external tool needed
// const unsigned char logo[] = { #embed "logo.png" };

// _BitInt(N) — exact-width integers of ANY bit width, not just 8/16/32/64
_BitInt(24) rgb_pixel;   // exactly 24 bits, no padding to 32

// Removed: implicit int ("x;" meaning "int x;") and K&R-style function
// definitions are gone — every function must have a real prototype now.
```

## B. Undefined-behavior and pitfalls cheat sheet

- Signed overflow is UB; unsigned wraps predictably — Section 3.3, 28.
- `arr[i]` is `*(arr + i)`; there is no bounds check, ever — Section 6.1.
- Arrays decay to pointers on function-call; `sizeof` on a decayed pointer
  gives the pointer's size, not the array's — Section 6.1.
- Never return a pointer to a local (stack) variable — Section 5.1.
- `strcpy`/`strcat`/`sprintf`/`gets` have no bound; prefer
  `snprintf`/`strlcpy`/`fgets` — Section 7.2, 29.
- `malloc`'s return value is not guaranteed non-NULL; always check —
  Section 10.1.
- Match every `malloc`/`strdup`/`_new` with exactly one `free`/`_free`;
  never twice, never after — Section 10.2, 10.3.
- `free(NULL)` is always safe; make your own `_free` functions mirror
  that — Section 10.3.
- Data races are UB, not just "unpredictable" — the compiler may assume
  they don't happen — Section 20.1, 28.
- Format-string/argument mismatches in `printf`/`scanf` are UB —
  Section 3.2, 22.1.
- `va_arg`'s type must match what the caller actually passed; there is no
  runtime check — Section 22.1.

## C. Compiler-flags cheat sheet

```bash
# Everyday development
clang -std=c17 -Wall -Wextra -Wpedantic -Werror -g -O0 \
      -fsanitize=address,undefined -o app app.c

# Thread-safety checking
clang -std=c17 -g -fsanitize=thread -o app app.c

# Fuzzing
clang -std=c17 -g -fsanitize=fuzzer,address -o fuzz_target fuzz.c

# Release / hardened build
clang -std=c17 -O2 -D_FORTIFY_SOURCE=2 -fstack-protector-strong \
      -fPIE -pie -Wl,-z,relro,-z,now -o app app.c

# Static analysis
clang-tidy app.c -- -std=c17
scan-build clang -std=c17 -c app.c
cppcheck --enable=all --std=c17 app.c
```

## D. Glossary

- **Translation unit** — one `.c` file plus everything its `#include`s
  pull in, after preprocessing; the compiler's actual unit of work.
- **ABI (Application Binary Interface)** — the compiled-code-level
  contract (calling convention, struct layout, symbol names) that lets
  separately compiled pieces link and call each other correctly.
- **UB (Undefined Behavior)** — behavior the standard places zero
  requirements on; the compiler may assume it never occurs.
- **Amortized O(1)** — an operation that's occasionally expensive but
  averages out to constant time over many calls (Section 16.1's doubling
  dynamic array).
- **Intrusive container** — a container whose bookkeeping fields are
  embedded inside the payload struct itself, recovered via
  `container_of`, instead of the payload being stored inside generic
  container nodes (Section 19).
- **Opaque pointer/handle** — a pointer to a struct whose definition is
  hidden from the caller, used to build stable, ABI-safe public APIs
  (Section 40.1).
- **False sharing** — a performance bug where independent data on the
  same CPU cache line causes unnecessary cross-core cache invalidation
  (Section 42.3).

## E. Further reading and real codebases worth reading

- *The C Programming Language* (K&R, 2nd ed.) — dated in places (pre-C99)
  but still the sharpest, shortest tour of the language's core ideas.
- *Modern C* (Jens Gustedt, free online) — written around C11/C17,
  the best bridge from "K&R C" to the standard this guide targets.
- SQLite source (`sqlite.org`) — Section 33; read `os_unix.c` for the VFS
  layer, `vdbe.c` for the opcode interpreter.
- Redis source (`github.com/redis/redis`) — Section 34; start with
  `src/sds.c` and `src/ae.c`, both short and self-contained.
- curl source (`github.com/curl/curl`) — Section 35; `lib/easy.c` and
  `lib/multi.c`.
- cJSON source (`github.com/DaveGamble/cJSON`) — Section 36; small enough
  to read start to finish in one sitting after this guide's Part II.
- The Linux kernel's `include/linux/list.h` — Section 19's
  `container_of`, in its original, real-world form.
- `man 7 signal`, `man 2 mmap`, `man 3 pthread_create` — the primary
  sources this guide's socket/threading/memory sections summarize;
  the man pages are the actual specification of the syscalls you're
  calling.
- **[`wiki/DSA/real-life-ds-algo-guide.md`](../DSA/real-life-ds-algo-guide.md)**
  — this guide's data-structures-and-algorithms companion. It covers
  every classic DS/algo topic (arrays through Dijkstra and topological
  sort) ELI5-first, with a Go implementation *and* a C implementation
  side by side for each one — the C versions are the beginner on-ramp to
  Sections 17–23's from-scratch data structures and Section 22's
  Dijkstra/graph material above, in the same toy-box-analogy style as
  this guide's "Real-world example" boxes.
