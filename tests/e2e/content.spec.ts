// Content sanity (desktop only). Each document is a landing page (/<doc>/), one
// page per chapter (/<doc>/<chapter>/), and the whole book on one page
// (/<doc>/all/, noindex). Checks: everything loads cleanly, links and anchors
// resolve, Markdown rendered fully, every chapter page carries the metadata
// search engines need, and the reader's navigation follows the course. Tables
// of contents are snapshotted, so a push that drops or renames chapters fails
// until the change is accepted with --update-snapshots.
import { test, expect, type APIRequestContext } from "@playwright/test";
import { docs, guide, layout, manifest, watchErrors } from "./site";

const BASE = "https://frontendlabs.xyz/";

/** The one-page edition when the document has chapters, else the landing page itself. */
async function fullPage(request: APIRequestContext, url: string) {
  return (await request.get(`${url}all/`)).ok() ? `${url}all/` : url;
}

let sitemap: Set<string>;
test.beforeAll(async ({ request }) => {
  const xml = await (await request.get("sitemap.xml")).text();
  sitemap = new Set([...xml.matchAll(/<loc>([^<]+)<\/loc>/g)].map((m) => m[1]));
});

for (const d of docs) {
  test.describe(d.url, () => {
    test("landing page loads cleanly with intact structure", async ({ page }) => {
      const check = watchErrors(page);
      const res = await page.goto(d.url);
      expect(res!.status()).toBe(200);
      await expect(page).toHaveTitle(/\S.* · frontendlabs$/);
      await expect(page.locator("body")).toHaveAttribute("data-guide", d.guide.slug);
      await expect(page.locator("h1.book-title")).toHaveCount(1);
      const dupes = await page.evaluate(() => {
        const all = [...document.querySelectorAll("[id]")].map((e) => e.id);
        return all.filter((id, i) => all.indexOf(id) !== i);
      });
      expect(dupes, "duplicate element ids").toEqual([]);
      await page.waitForLoadState("networkidle");
      check();
    });

    test("one-page edition: anchors resolve, Markdown rendered, copy buttons", async ({ page, request }) => {
      const url = await fullPage(request, d.url);
      await page.goto(url);
      if (url.endsWith("all/")) {
        await expect(page.locator('meta[name="robots"]')).toHaveAttribute("content", /noindex/);
      }
      const missing = await page.evaluate(() =>
        [...document.querySelectorAll<HTMLAnchorElement>("#tocList a.toc-link")]
          .filter((a) => a.hash && a.pathname === location.pathname)
          .map((a) => decodeURIComponent(a.hash.slice(1)))
          .filter((id) => !document.getElementById(id)));
      expect(missing, "contents entries without a target").toEqual([]);

      const leftovers = await page.evaluate(() => {
        const bad: string[] = [];
        for (const el of document.querySelectorAll("#content :is(p, li, td, th, blockquote)")) {
          const isPara = el.tagName === "P";
          const clone = el.cloneNode(true) as HTMLElement;
          clone.querySelectorAll("code, pre, kbd").forEach((c) => c.remove());
          const t = clone.textContent ?? "";
          if (/\]\((?:\.{0,2}\/|https?:\/\/|#)[^)\s]*\)/.test(t)) bad.push(`unrendered link: ${t.slice(0, 80)}`);
          // "# of /24s" in a table is real text; a paragraph starting "## Title" is a broken heading
          if (isPara && /^#{1,6} [A-Z]/.test(t.trim())) bad.push(`unrendered heading: ${t.slice(0, 80)}`);
          if (t.includes("```")) bad.push(`unrendered code fence: ${t.slice(0, 80)}`);
          if (t.includes("Read this on frontendlabs.xyz")) bad.push("GitHub-only notice shown on the site");
        }
        return bad;
      });
      expect(leftovers).toEqual([]);

      const pres = await page.locator("#content pre").count();
      if (pres > 0) await expect(page.locator("#content button.copy")).toHaveCount(pres);
    });

    test("table of contents matches the snapshot", async ({ page }) => {
      await page.goto(d.url);
      // titles only: the reading time (<em>) shifts with any wording edit and is not content
      const toc = await page.locator("#tocList a.toc-link > span").allTextContents();
      expect(toc.length, "contents list is not empty").toBeGreaterThan(0);
      expect(toc.map((t) => t.replace(/\s+/g, " ").trim()).join("\n") + "\n")
        .toMatchSnapshot(`toc-${d.url.replace(/\/$/, "").replace(/\//g, "__")}.txt`);
    });

    test("search metadata: landing page and every chapter page", async ({ page, request }) => {
      await page.goto(d.url);
      expect(await page.locator('meta[name="robots"]').count(), "landing pages are indexable").toBe(0);
      await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", BASE + d.url);
      const landingLD = JSON.parse((await page.locator('script[type="application/ld+json"]').textContent())!);
      const types = landingLD["@graph"].map((n: { "@type": string }) => n["@type"]);
      expect(types).toContain(d.main ? "Course" : "TechArticle");
      expect(types).toContain("BreadcrumbList");
      expect(sitemap.has(BASE + d.url), "landing page in sitemap").toBe(true);

      const chapters = await page.locator("#contents li a").evaluateAll((as) =>
        (as as HTMLAnchorElement[]).filter((a) => !a.hash).map((a) => a.pathname.slice(1)));
      const titles = new Set<string>();
      for (const url of chapters) {
        const res = await request.get(url);
        expect(res.status(), url).toBe(200);
        const html = await res.text();
        const title = html.match(/<title>([^<]*)<\/title>/)![1];
        expect(title, `${url} title`).toMatch(/\S.* · frontendlabs$/);
        expect(titles.has(title), `${url} title is unique`).toBe(false);
        titles.add(title);
        expect(html, `${url} canonical`).toContain(`<link rel="canonical" href="${BASE}${url}">`);
        expect(html.match(/<meta name="description" content="([^"]*)"/)![1].length, `${url} description`).toBeGreaterThan(20);
        expect(html, `${url} has its chapter as h1`).toMatch(/<h1 id="[^"]+" class="chapter-title">/);
        const ld = JSON.parse(html.match(/<script type="application\/ld\+json">([^<]*)<\/script>/)![1]);
        const article = ld["@graph"].find((n: { "@type": string }) => n["@type"] === "TechArticle");
        expect(article?.url, `${url} TechArticle`).toBe(BASE + url);
        expect(ld["@graph"].some((n: { "@type": string }) => n["@type"] === "BreadcrumbList")).toBe(true);
        // a chapter is either indexable and in the sitemap, or noindex (too short) and left out
        const noindex = /<meta name="robots" content="noindex/.test(html);
        expect(sitemap.has(BASE + url), `${url}: sitemap and robots agree`).toBe(!noindex);
      }
    });
  });
}

test("home page describes the site as a WebSite", async ({ page }) => {
  await page.goto("./");
  const ld = JSON.parse((await page.locator('script[type="application/ld+json"]').textContent())!);
  expect(ld["@graph"].map((n: { "@type": string }) => n["@type"])).toContain("WebSite");
});

test("search index lists every guide", async ({ request }) => {
  const idx: { t: string; u: string }[] = await (await request.get("search.json")).json();
  for (const g of manifest.guides) {
    expect(idx.some((e) => e.u === `${g.slug}/`), `${g.slug} in search.json`).toBe(true);
  }
});

test("old single-page links redirect to the chapter page", async ({ page }) => {
  await page.goto("rust/#6-ownership-moves-copies-and-drop");
  await expect(page).toHaveURL(/\/rust\/6-ownership-moves-copies-and-drop\/#6-ownership-moves-copies-and-drop$/);
  await expect(page.locator("h1.chapter-title")).toHaveText(/Ownership/);
});

test("chapter pages chain together, also by keyboard", async ({ page }) => {
  await page.goto("rust/6-ownership-moves-copies-and-drop/");
  await expect(page.locator('a[rel="prev"]')).toContainText("5. Functions");
  await page.locator('a[rel="next"]').click();
  await expect(page.locator("h1.chapter-title")).toHaveText(/^7\. Borrowing/);
  await page.keyboard.press("n");
  await expect(page.locator("h1.chapter-title")).toHaveText(/^8\. Slices/);
});

test.describe("reader navigation", () => {
  test.skip(layout !== "spotlight", "previous/next rules below are the spotlight layout's");

  for (const [i, slug] of (manifest.path ?? []).entries()) {
    test(`${slug}: previous/next guide follow the course path`, async ({ page }) => {
      await page.goto(`${slug}/`);
      const prev = manifest.path![i - 1];
      const next = manifest.path![i + 1];
      const nav = page.locator("a.guide-nav");
      if (prev) await expect(nav.filter({ hasText: "Previous guide" })).toContainText(guide(prev).title);
      else await expect(nav.filter({ hasText: "Previous guide" })).toHaveCount(0);
      if (next) await expect(nav.filter({ hasText: "Next guide" })).toContainText(guide(next).title);
      else await expect(nav.filter({ hasText: "Next guide" })).toHaveCount(0);
    });
  }

  test("spotlight deep links land on their headings", async ({ page }) => {
    const sp = manifest.spotlight!;
    await page.goto("./");
    const hrefs = await page.locator(`a[href^="${sp.guide}/"][href*="#"]`).evaluateAll((as) => as.map((a) => a.getAttribute("href")!));
    expect(hrefs).toHaveLength(sp.links.length);
    for (const [i, l] of sp.links.entries()) {
      expect(hrefs[i], l.title).toMatch(new RegExp(`#${l.anchor}$`));
      await page.goto(hrefs[i]);
      const target = page.locator(`[id="${l.anchor}"]`);
      await expect(target, l.title).toHaveCount(1);
      await expect(target).toBeInViewport();
    }
  });
});

test("reader controls: theme, contents filter, keyboard", async ({ page }) => {
  await page.goto("rust/");
  const theme = () => page.locator("html").getAttribute("data-theme");
  const before = await theme();
  await page.locator("#themeBtn").click();
  expect(await theme()).not.toBe(before);
  await page.keyboard.press("d");
  const third = await theme();
  expect(new Set([before, third]).size).toBe(2);

  const all = await page.locator("#tocList a.toc-link:visible").count();
  await page.locator("#tocFilter").fill("lifetimes");
  await expect.poll(() => page.locator("#tocList a.toc-link:visible").count()).toBeLessThan(all);
  await expect(page.locator("#tocList a.toc-link:visible").first()).toContainText(/lifetime/i);
});
