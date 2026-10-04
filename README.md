# frontendlabs.xyz — live field-guide reader

A book-style web app for the wiki's real-life field guides. The Markdown in
[`docs/`](docs/) is the source. A small Go program renders it into a static site
in `dist/`, and regenerates the affected page **within about a second** whenever
a guide changes, either in `docs/` or in the original wiki folders.

```bash
make serve        # sync, build a preview into .preview/, serve on http://localhost:8080, regenerate on change
make build        # production build into dist/  -> deploy dist/ to https://frontendlabs.xyz
make sync         # copy the latest guides from the wiki into docs/
make dev          # serve + rebuild on template/CSS edits (needs Node for Tailwind)
make check        # build, then verify every link and #anchor in dist/ for production
```

Requirements: Go 1.22+ only. Node is needed **only** to change the CSS
(`npm install && npm run css`). The compiled stylesheet is committed in
`web/static/app.css`.

## Why this stack

| Option | Verdict |
|---|---|
| htmx + Tailwind | Needs a server rendering HTML fragments anyway. Partial swaps add nothing to reading long books, and break deep links and the browser's find-in-page. |
| Vanilla JS + Tailwind (render Markdown in the browser) | Parses up to 1.2 MB of Markdown on every visit, needs a web server even locally (`fetch` is blocked on `file://`), and is invisible to search engines. |
| **Go static-site generator + Tailwind (chosen)** | Renders once, serves plain HTML. One binary does sync, build, serve, and watch. Uses the same heading anchors as the wiki's builders, so the guides' cross-links keep working. |
| Astro Starlight / VitePress / Docusaurus / mdBook | Good tools, but they expect front matter or one file per chapter. These guides are single 10k–20k-line books, and the book reader (themes, resume, chapter navigation) would have to be rebuilt anyway. |

The result works from any static host or straight from disk (`dist/index.html`).

## What you get

- **Home page:** an icon card for every guide, grouped by category, with the
  systems series shown in reading order (OS → networking → security → HTTPS).
  Each card shows chapters, reading time, and runnable Go programs, plus
  search across every guide and chapter (`/`) and a "continue reading" strip.
- **Book reader** for every guide: a title page, a contents sidebar with parts,
  chapters and reading times, a chapter filter, scroll tracking, a progress bar,
  paper/light/night themes, font size and column width, resume where you left
  off, previous/next chapter and guide, coloured cards for "In one sentence" /
  "Build it in Go" / "Common mistakes" / "Check yourself" sections, copy buttons
  on code, checklists that remember their state, keyboard shortcuts
  (`t d w + - n p /`), and print-to-PDF.
- **Links between guides** (`../v2-https/real-life-guide-v1.md#chapter-25-...`)
  are rewritten to the app's URLs, and links written for GitHub's anchor style
  also resolve.
- **Original editions:** the wiki's own HTML books are published under
  `/originals/`, and each guide's Markdown under `/docs/`, both linked from the
  reader's `⋯` menu.

## How regeneration works

```
wiki/<folder>/*.md ──sync──▶ docs/<folder>/*.md ──render──▶ dist/<slug>/index.html ──SSE──▶ open browser tabs reload
     (originals)        (mirror, same paths)        (only the changed guide)        (only tabs showing that page)
```

- `serve` and `watch` poll file fingerprints (modification time and size)
  every 700 ms. This is standard-library Go, works the same on macOS and Linux,
  and needs no file-notification dependency.
- A changed **wiki original** is copied into `docs/` (written to a temporary
  file and renamed, so nothing ever reads half a file). A changed file in
  **`docs/`** re-renders just that document plus the home page and search index.
  In testing, a guide edit was live in the browser within 0.3 s, and the
  re-render took about 8 ms.
- `serve` pushes the changed URLs to open tabs over Server-Sent Events. A tab
  reloads only if it shows that page, and keeps its scroll position.
- `watch` does the same without a web server, for a host that only regenerates
  `dist/` for a separate static server.
- A rebuild error is logged and the last good page keeps being served.

**`docs/` vs the wiki:** `sync` copies wiki → `docs/` whenever a wiki original
changes. Edit the wiki originals, or run with `-sync=false` if you want to
edit `docs/` directly and never have it overwritten.


## Adding or changing a guide

Everything about presentation is in [`guides.json`](guides.json):

```json
{"slug": "dsa", "title": "Data Structures & Algorithms", "icon": "🧩",
 "category": "foundations", "step": "", "level": "Beginner → Advanced",
 "summary": "…", "source": "DSA/real-life-ds-algo-guide.md",
 "originals": ["DSA/real-life-ds-algo-guide.html"],
 "pages": [{"slug": "week1", "source": "…/week1.md"}]}
```

`source` and `originals` are paths relative to the wiki root, and the same paths
under `docs/`. `pages` adds extra documents as sub-pages of a guide (used for the
Go plan's 19 documents). Add an entry, run
`make sync build`, and the guide appears on the home page.

## Testing: the commit gate

Every commit and every deploy runs an end-to-end suite (Playwright, Chromium)
against a fresh build of the site, served as static files like GitHub Pages.

```bash
npm ci && npx playwright install chromium   # once
make hooks        # once per clone: install the pre-commit hook
make test         # links + the full suite (what the hook runs, ~20 s)
make test-ci      # the same without screenshot comparisons (what CI runs)
npx playwright test --update-snapshots      # accept an intended design or content change
```

| Spec | What it guards |
|---|---|
| `home.spec.ts` | every guide has exactly one card; sections, the course path ribbon, the spotlight and its deep links follow `guides.json`; search; themes |
| `layout.spec.ts` | no page scrolls sideways on a phone; spotlight columns, one-row path, 3-column cards on desktop and 1 on phones; sticky bars |
| `style.spec.ts` | assets carry their current `?v=` content hash; the stylesheet is applied; every theme's colours and WCAG text contrast |
| `content.spec.ts` | all documents load without browser errors or duplicate ids; every contents entry has a target; Markdown rendered fully; copy buttons; previous/next follow the course; **table-of-contents snapshots** |
| `visual.spec.ts` (`@visual`) | screenshot comparisons of the home page (paper and night) and a guide page, desktop and phone |

Tests read `guides.json`, so a new guide is covered automatically. Snapshots
live in `tests/e2e/__snapshots__/` (contents lists, the same on every OS) and
`tests/e2e/__screenshots__/<platform>/` (screenshots, per OS because fonts
differ). CI runs on Linux and has no screenshot baselines, so it skips
`@visual`; the pre-commit hook runs everything. Review the diff of any
updated snapshot before committing it: that is the change readers will see.

## Home page layouts

`site.home_layout` in `guides.json` picks the home page. Change it and rebuild
(`make build`, or let `make serve` pick it up) to switch:

| Value | Template | What it shows |
|---|---|---|
| `"spotlight"` | `web/templates/home-spotlight.html` | AI spotlight first (the `spotlight` block: featured guide, deep links, "coming next"), the course as a path ribbon, then topic `sections` (Systems, Networking, Security, Programming) |
| `"classic"` (or unset) | `web/templates/home.html` | the original layout: one section per entry in `categories` |

In the spotlight layout, a guide's "Next guide" button follows `path` (the
course order) when the guide is on it, and its section order otherwise, so
moving HTTPS into the Networking section does not change the course. The
build fails if a guide is in no section, so nothing can silently disappear
from the home page. The classic layout ignores `sections`, `path` and
`spotlight`, and builds exactly what it did before.

## Layout

```
guides.json            manifest: guides, icons, categories, series steps
docs/                  the source Markdown (+ original HTML), mirroring the wiki's paths
cmd/guides/            CLI: sync | build | serve | watch
internal/site/
  manifest.go          guides.json, and doc -> URL mapping
  sync.go              wiki -> docs copy, change fingerprints
  markdown.go          goldmark (GFM), wiki-compatible heading IDs, link rewriting
  book.go              title page, chapters, section cards, contents, reading time, GitHub anchor aliases
  build.go             pages, home, search index, assets, sitemap/robots/404/CNAME; incremental rebuild
  originals.go         publishes the original HTML editions with their links fixed
  serve.go             watcher, live-reload (SSE), static server
web/templates/         home.html, book.html, partials.html  (html/template)
web/static/            reader.js, home.js, app.css (Tailwind output), brand/ (logo sizes from FrontEndLabs.png)
web/styles/input.css   Tailwind v4 source: theme tokens and book typography
scripts/checklinks.py  verifies every link and anchor in dist/
```

## Original editions

The wiki's own HTML books are published under `/originals/`. Their links were
written for the wiki's folder layout (sibling `.md` files, folders,
GitHub-style anchors), so the build rewrites every `<a href>` that wouldn't
resolve under `/originals/` to point at the app's rendered page and anchor
instead. Example markup inside code blocks is never touched.

## License

- **Guides** (everything in `docs/`, and the pages built from it):
  [CC BY-NC-SA 4.0](LICENSE-CONTENT). You may share and adapt them for
  non-commercial use, with credit to frontendlabs.xyz and under the same license.
- **Code** (the generator: `cmd/`, `internal/`, `web/`, `scripts/`): [MIT](LICENSE).

Code samples inside the guides may also be used under the MIT license, so you
can copy them into your own projects, commercial or not.
