#!/usr/bin/env python3
"""Production link check for the generated site.

Scans every HTML file in dist/ (guide pages, original editions, 404) and reads
URLs from real tags only (href, src, srcset, content of og:url/og:image). Fails on:
  - local/dev references: localhost, 127.0.0.1, file:, filesystem paths, the live-reload endpoint
  - protocol-relative URLs (//host/...)
  - relative or root-relative links whose file or #anchor doesn't exist in dist/
  - absolute links to the site's own domain that don't exist in dist/
  - sitemap.xml entries without a page, or canonical URLs that don't match the page
Warns (doesn't fail) on insecure http:// links to other sites.

Usage: scripts/checklinks.py [dist] [base_url]
"""
import re
import sys
from pathlib import Path
from urllib.parse import unquote, urlparse

dist = Path(sys.argv[1] if len(sys.argv) > 1 else "dist").resolve()
base = (sys.argv[2] if len(sys.argv) > 2 else "https://frontendlabs.xyz").rstrip("/")
host = urlparse(base).netloc

TAG = re.compile(r"<([a-zA-Z][a-zA-Z0-9-]*)(\s[^<>]*?)?/?>")
ATTR = re.compile(r'([a-zA-Z:-]+)="([^"]*)"')
LOCAL = re.compile(r"^(https?://(localhost|127\.0\.0\.1|0\.0\.0\.0|\[::1\])([:/]|$)|file:|/Users/|/private/|/tmp/|/home/|[A-Za-z]:\\)")

errors, warnings = [], []
id_cache = {}


def ids(p: Path):
    if p not in id_cache:
        id_cache[p] = set(re.findall(r'\sid="([^"]+)"', p.read_text(errors="ignore")))
    return id_cache[p]


def resolve(page: Path, url: str):
    """Return the dist file a URL points to, or None if it's outside the site."""
    if url.startswith(base + "/") or url == base:
        rel = url[len(base):].lstrip("/")
        t = dist / unquote(rel.split("#")[0])
    elif url.startswith("/"):
        t = dist / unquote(url.split("#")[0].lstrip("/"))
    else:
        f = url.split("#")[0]
        t = (page.parent / unquote(f)).resolve() if f else page
    if t.is_dir():
        t = t / "index.html"
    return t


def check(page: Path, url: str, where: str):
    rel = str(page.relative_to(dist))
    if not url or url.startswith(("data:", "mailto:", "javascript:void")):
        return
    if LOCAL.match(url):
        errors.append((rel, url, "local/dev reference")); return
    if url.startswith("//"):
        errors.append((rel, url, "protocol-relative URL")); return
    p = urlparse(url)
    if p.scheme in ("http", "https") and p.netloc != host:
        if p.scheme == "http":
            warnings.append((rel, url, "insecure http:// link"))
        return
    if p.scheme and p.scheme not in ("http", "https"):
        return
    if url.startswith("/") and rel != "404.html":
        warnings.append((rel, url, "root-relative (works only at the domain root)"))
    t = resolve(page, url)
    if not t.exists():
        errors.append((rel, url, f"missing file ({where})")); return
    frag = unquote(url.partition("#")[2])
    if frag and t.suffix == ".html" and frag not in ids(t) and "'+" not in frag:
        errors.append((rel, url, "missing anchor"))


pages = sorted(dist.rglob("*.html"))
for page in pages:
    html = page.read_text(errors="ignore")
    if "/__live" in html:
        errors.append((str(page.relative_to(dist)), "/__live", "dev live-reload client in a production page"))
    for tag, attrs in TAG.findall(html):
        a = dict(ATTR.findall(attrs or ""))
        for name in ("href", "src"):
            if name in a:
                check(page, a[name], f"<{tag} {name}>")
        if "srcset" in a:
            for part in a["srcset"].split(","):
                check(page, part.strip().split(" ")[0], f"<{tag} srcset>")
        if tag == "meta" and a.get("property") in ("og:url", "og:image"):
            check(page, a.get("content", ""), a["property"])
        if tag == "link" and a.get("rel") == "canonical":
            expect = base + "/" + str(page.parent.relative_to(dist)).replace(".", "") + "/"
            expect = expect.replace("//", "/").replace("https:/", "https://").replace("http:/", "http://")
            if page.name == "index.html" and a.get("href") != expect.replace(base + "//", base + "/"):
                if not (page.parent == dist and a.get("href") == base + "/"):
                    errors.append((str(page.relative_to(dist)), a.get("href", ""), f"canonical should be {expect}"))

sitemap = dist / "sitemap.xml"
if sitemap.exists():
    for loc in re.findall(r"<loc>([^<]+)</loc>", sitemap.read_text()):
        if not resolve(dist / "index.html", loc).exists():
            errors.append(("sitemap.xml", loc, "no page for sitemap URL"))
else:
    errors.append(("sitemap.xml", "", "missing"))

links = sum(1 for _ in [0])
print(f"{len(pages)} HTML files checked; {len(errors)} errors, {len(warnings)} warnings")
for e in errors[:50]:
    print("  ERROR  ", *e)
seen = set()
for w in warnings:
    if (w[2], w[1]) in seen:
        continue
    seen.add((w[2], w[1]))
    if len(seen) <= 15:
        print("  warn   ", *w)
sys.exit(1 if errors else 0)
