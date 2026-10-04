// Content sanity for every rendered document (desktop only): it loads cleanly,
// its contents list is intact and every entry has a target, Markdown rendered
// fully, and the reader's navigation follows the course. Tables of contents
// are snapshotted, so a push that drops or renames chapters fails until the
// change is accepted with --update-snapshots.
import { test, expect } from "@playwright/test";
import { docs, guide, layout, manifest, watchErrors } from "./site";

for (const d of docs) {
  test.describe(d.url, () => {
    test("loads cleanly with intact structure", async ({ page }) => {
      const check = watchErrors(page);
      const res = await page.goto(d.url);
      expect(res!.status()).toBe(200);
      // the browser title comes from the document's own heading, branded for the site
      await expect(page).toHaveTitle(/\S.* · frontendlabs$/);
      await expect(page.locator(`[data-guide="${d.guide.slug}"]`).first()).toBeAttached();
      expect(await page.locator("#content :is(h1, h2)").count(), "has headings").toBeGreaterThan(0);

      const ids = await page.evaluate(() => {
        const all = [...document.querySelectorAll("[id]")].map((e) => e.id);
        return all.filter((id, i) => all.indexOf(id) !== i);
      });
      expect(ids, "duplicate element ids").toEqual([]);
      await page.waitForLoadState("networkidle");
      check();
    });

    test("every contents entry points at a heading on the page", async ({ page }) => {
      await page.goto(d.url);
      const missing = await page.evaluate(() =>
        [...document.querySelectorAll<HTMLAnchorElement>("a.toc-link")]
          .filter((a) => a.hash && a.pathname === location.pathname)   // entries may also link to sub-pages
          .map((a) => decodeURIComponent(a.hash.slice(1)))
          .filter((id) => !document.getElementById(id)));
      expect(missing).toEqual([]);
    });

    test("Markdown rendered completely", async ({ page }) => {
      await page.goto(d.url);
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
        }
        return bad;
      });
      expect(leftovers).toEqual([]);
    });

    test("code blocks get copy buttons", async ({ page }) => {
      await page.goto(d.url);
      const pres = await page.locator("#content pre").count();
      test.skip(pres === 0, "no code on this page");
      await expect(page.locator("#content button.copy")).toHaveCount(pres);
    });

    test("table of contents matches the snapshot", async ({ page }) => {
      await page.goto(d.url);
      // titles only: the reading time (<em>) shifts with any wording edit and is not content
      const toc = await page.locator("a.toc-link > span").allTextContents();
      expect(toc.length, "contents list is not empty").toBeGreaterThan(0);
      expect(toc.map((t) => t.replace(/\s+/g, " ").trim()).join("\n") + "\n")
        .toMatchSnapshot(`toc-${d.url.replace(/\/$/, "").replace(/\//g, "__")}.txt`);
    });
  });
}

test("search index lists every guide", async ({ request }) => {
  const idx: { t: string; u: string }[] = await (await request.get("search.json")).json();
  for (const g of manifest.guides) {
    expect(idx.some((e) => e.u === `${g.slug}/`), `${g.slug} in search.json`).toBe(true);
  }
});

test.describe("reader navigation", () => {
  test.skip(layout !== "spotlight", "previous/next rules below are the spotlight layout's");

  for (const [i, slug] of (manifest.path ?? []).entries()) {
    test(`${slug}: previous/next follow the course path`, async ({ page }) => {
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
    for (const l of sp.links) {
      await page.goto(`${sp.guide}/#${l.anchor}`);
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

  const all = await page.locator("a.toc-link:visible").count();
  await page.locator("#tocFilter").fill("lifetimes");
  await expect.poll(() => page.locator("a.toc-link:visible").count()).toBeLessThan(all);
  await expect(page.locator("a.toc-link:visible").first()).toContainText(/lifetime/i);
});
