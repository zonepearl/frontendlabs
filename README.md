<div align="center">

# Frontend Labs

**Long-form engineering field guides, read like books.**
From the kernel to LLMs: operating systems, Docker & Kubernetes, networking, security,
Go, Rust, Python, C and AI, with runnable labs in every guide.

[**frontendlabs.xyz**](https://frontendlabs.xyz) · [Guides](#the-guides) · [What's new](#whats-new) · [Run it locally](#run-it-locally) · [Contributing](#contributing)

[![deploy](https://github.com/zonepearl/frontendlabs/actions/workflows/pages.yml/badge.svg)](https://github.com/zonepearl/frontendlabs/actions/workflows/pages.yml)
[![e2e](https://github.com/zonepearl/frontendlabs/actions/workflows/ci.yml/badge.svg)](https://github.com/zonepearl/frontendlabs/actions/workflows/ci.yml)
[![content: CC BY-NC-SA 4.0](https://img.shields.io/badge/content-CC%20BY--NC--SA%204.0-blue)](LICENSE-CONTENT)
[![code: MIT](https://img.shields.io/badge/code-MIT-green)](LICENSE)

</div>

---

## The guides

| | Guide | Topic | Level |
|---|---|---|---|
| 🤖 | [AI from Zero to LLMs](https://frontendlabs.xyz/ai/) | AI | Beginner → Expert |
| 📐 | [Real-Life Mathematics](https://frontendlabs.xyz/maths/) | AI | Beginner → Expert |
| 🖥️ | [Operating Systems, Linux & Containers](https://frontendlabs.xyz/linux/) | Systems | Beginner → Expert · step 1 |
| ☸️ | [Docker & Kubernetes — The Real-Life Field Guide](https://frontendlabs.xyz/kubernetes/) | Systems | Beginner → Production · after step 1 |
| 📈 | [Scale, Load & Performance Testing](https://frontendlabs.xyz/scale-perf/) | Systems | In depth |
| 🧭 | [The OSI Model, One Click at a Time](https://frontendlabs.xyz/osi/) | Networking | Overview · step 2a |
| 🌐 | [Networking from Zero (TCP/IP)](https://frontendlabs.xyz/tcp-ip/) | Networking | Beginner → Expert · step 2b |
| 🔗 | [The HTTPS Request Lifecycle](https://frontendlabs.xyz/https/) | Networking | Intermediate → Expert · step 4 |
| 📡 | [TCP/IP — The Complete Field Guide (reference)](https://frontendlabs.xyz/tcp-ip-reference/) | Networking | Reference |
| 🔒 | [Security from Zero](https://frontendlabs.xyz/security-from-zero/) | Security | Beginner → Advanced · step 3a |
| 🕵️ | [Security Engineering in Depth](https://frontendlabs.xyz/security-in-depth/) | Security | Advanced · step 3b |
| 🐹 | [Go — The Complete Field Guide](https://frontendlabs.xyz/go/) | Programming | Beginner → Expert · alongside |
| 🗓️ | [Go 120-Day Engineering Plan](https://frontendlabs.xyz/go-plan/) | Programming | Day by day |
| 🦀 | [Rust — The Complete Field Guide](https://frontendlabs.xyz/rust/) | Programming | Beginner → Expert · alongside |
| 🐍 | [Python — The Complete Field Guide](https://frontendlabs.xyz/python/) | Programming | Beginner → Expert · alongside |
| 🛠️ | [C — The Complete Field Guide](https://frontendlabs.xyz/c/) | Programming | Beginner → Expert |
| 🟨 | [JavaScript — The Complete Field Guide](https://frontendlabs.xyz/javascript/) | Programming | Beginner → Expert |
| 🧩 | [Data Structures & Algorithms](https://frontendlabs.xyz/dsa/) | Programming | Beginner → Advanced |

The systems course reads in order: **1 OS → 2 Networking → 3 Security → 4 HTTPS**, with Go, Rust and Python alongside.
**Docker & Kubernetes** picks up after step 1 and takes containers to production.
Each guide's Markdown is in [`docs/`](docs/); the site is the official, always-current edition.

## What's new

- **Real-Life Mathematics, in depth** ([`/maths/`](https://frontendlabs.xyz/maths/)).
  The linear algebra, calculus and probability & statistics chapters now explain
  every formula symbol by symbol, with a worked example under each one, the
  theorems behind them (spectral, Taylor, the Fundamental Theorem, CLT, Hoeffding,
  Eckart–Young) and how AI and networking use them: attention's `√d_k`, Adam,
  backprop by hand, PageRank, Markov loss models, TCP's RTO and throughput,
  Kelly's fair-sharing. New: integral calculus and "How to read a formula".
- **Docker & Kubernetes — The Real-Life Field Guide** ([`/kubernetes/`](https://frontendlabs.xyz/kubernetes/)).
  Beginner to production, built around one cloud system ("Orbit") assembled end to end:
  production Dockerfiles and supply-chain security, Kubernetes core and production workloads,
  probes and self-healing, autoscaling (HPA, KEDA, Karpenter), OpenTelemetry with
  Prometheus, Loki, Tempo and SLO burn-rate alerts, GitOps with Argo CD, canary releases,
  backup/DR and chaos. Part 0 is a scripted Linux lab VM on macOS. The 90-day DevOps
  plan ships alongside it at [`/kubernetes/plan/`](https://frontendlabs.xyz/kubernetes/plan/).
  It sits in **Systems**, after the OS & Linux guide.
- **AI from Zero to LLMs, Parts 20–21** (chapters 74–81): agent harnesses, long-running
  agent loops, agent evals (pass^k), deploying agents to production, agent interop
  (AGENTS.md, Skills, MCP, A2A), an overnight PR-opening agent, System One decision
  models (Jev), and JEPA world models. The home-page AI spotlight links straight to them.
- **Python — The Complete Field Guide** ([`/python/`](https://frontendlabs.xyz/python/)).
  The language and its internals, the series' OS, networking and security labs in Python,
  data science and AI/ML from NumPy to PyTorch and LLM apps, and a walkthrough of
  `requests` down to the socket. It runs alongside every step of the course.
- **Rust — The Complete Field Guide** ([`/rust/`](https://frontendlabs.xyz/rust/)), alongside
  Go across the systems series, together with the new home page: an AI spotlight,
  the engineering-path ribbon, and topic sections (Systems, Networking, Security, Programming).

## What the site gives readers

- **One page per chapter**, plus a guide overview with full contents and a
  "read as one page" edition for find-in-page and printing.
- **Search** across every guide and chapter (`/`), **continue reading** where you left off,
  and **paper / light / night** themes with adjustable text size and width.
- **Readable on any screen**, with copy buttons on code, keyboard shortcuts
  (`t` contents, `n`/`p` chapters, `d` theme), and print-to-PDF.
- **Search-engine ready**: a canonical URL, description and `<h1>` per chapter,
  JSON-LD (`Course`, `TechArticle`, `BreadcrumbList`), and a sitemap of every indexable page.

## Run it locally

Requires **Go** (version in [`go.mod`](go.mod)). Node is needed only to rebuild the CSS or run the tests.

```bash
make serve     # build a preview and serve http://localhost:8080; pages reload as guides change
make build     # production build into dist/
make sync      # copy the latest guides from the wiki into docs/
make check     # build, then verify every link and #anchor
make test      # check + the full end-to-end suite (what every commit must pass)
```

`serve` watches the wiki and `docs/` and re-renders only what changed, usually within a second.

## How it works

```text
wiki/*.md ──sync──▶ docs/*.md ──build──▶ dist/  ──GitHub Pages──▶ frontendlabs.xyz
                                          ├── <guide>/              overview + contents
                                          ├── <guide>/<chapter>/    one page per chapter (indexed)
                                          └── <guide>/all/          the whole book (noindex)
```

A small Go generator renders the Markdown with [goldmark](https://github.com/yuin/goldmark),
splits each guide into chapter pages, and resolves every link and `#anchor`
(including links between guides) to the page that holds it. Old one-page links
such as `/rust/#6-ownership…` redirect to the right chapter. There is no framework
and no client-side rendering: the output is static HTML that works on any host.

## Adding a guide

Add an entry to [`guides.json`](guides.json), then `make sync test`:

```json
{
  "slug": "dsa",
  "title": "Data Structures & Algorithms",
  "icon": "🧩",
  "level": "Beginner → Advanced",
  "summary": "Plain-language explanations and Go implementations for every core topic.",
  "source": "DSA/real-life-ds-algo-guide.md",
  "originals": ["DSA/real-life-ds-algo-guide.html"]
}
```

List its slug in one of the home page `sections`, and in `path` if it is part of the systems course.
The build fails if a guide is missing from every section, and the test suite covers new guides automatically.

## Quality gate

Every commit and every deploy must pass an end-to-end suite
([Playwright](https://playwright.dev), Chromium, desktop and phone):

| Area | Guards |
|---|---|
| Home | one card per guide; sections, course path and spotlight follow `guides.json`; search; themes |
| Layout | nothing scrolls sideways on phones; grids, columns and sticky bars where the design expects them |
| Style | current asset hashes; stylesheet applied; WCAG contrast in every theme |
| Content | every page loads cleanly; anchors resolve; Markdown fully rendered; metadata on every chapter page; contents snapshots |
| Visual | screenshot comparisons of the home page, a guide overview and a chapter |

```bash
npm ci && npx playwright install chromium   # once
make hooks                                  # once per clone: run `make test` before every commit
npx playwright test --update-snapshots      # accept an intended design or content change, then review the diff
```

## Contributing

`main` is protected: changes arrive through pull requests whose `e2e` check passed,
with no force pushes and no bypass. Workflows use pinned actions and read-only tokens,
and secret scanning and Dependabot are enabled. See [SECURITY.md](SECURITY.md) to report a vulnerability.

```bash
git switch -c my-change
make test && git commit -am "Describe the change" && git push -u origin my-change
# open a pull request; it merges once the e2e check is green
```

<details>
<summary><b>Project layout</b></summary>

```text
guides.json             guides, home page sections, course path, spotlight
docs/                   guide Markdown (+ original HTML editions), mirroring the wiki
cmd/guides/             CLI: sync | build | serve | watch
internal/site/          generator: markdown, chapter split, links, SEO, home, sitemap, live reload
web/templates/          home, book (overview / chapter / one page), partials
web/static/             reader.js, home.js, app.css (Tailwind output), brand assets
web/styles/input.css    Tailwind v4 source: theme tokens and book typography
tests/e2e/              Playwright specs and snapshots
scripts/                link checker, GitHub hardening script
.github/                workflows, rulesets, Dependabot, CODEOWNERS
```

`site.home_layout` in `guides.json` switches the home page between `"spotlight"` (default) and `"classic"`.
The wiki's own HTML books are published under `/originals/` with their links rewritten to the site.

</details>

## License

- **Guides** (`docs/` and the pages built from it): [CC BY-NC-SA 4.0](LICENSE-CONTENT).
  Share and adapt them for non-commercial use, with credit to frontendlabs.xyz.
- **Code** (`cmd/`, `internal/`, `web/`, `scripts/`) and **code samples inside the guides**: [MIT](LICENSE).
