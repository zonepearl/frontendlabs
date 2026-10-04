// Home page: every guide is reachable, sections and the course path follow
// guides.json, and search and themes work. Runs on desktop and phone.
import { test, expect } from "@playwright/test";
import { manifest, layout, guide, watchErrors } from "./site";

test.beforeEach(async ({ page }) => {
  await page.goto("./");
});

test("loads without browser errors", async ({ page }) => {
  const check = watchErrors(page);
  await page.reload();
  await page.waitForLoadState("networkidle");
  check();
});

test("every guide has exactly one card that links to it", async ({ page }) => {
  for (const g of manifest.guides) {
    const card = page.locator(`[data-guide="${g.slug}"]`);
    await expect(card, `card for ${g.slug}`).toHaveCount(1);
    await expect(card).toContainText(g.title);
    // the featured card is a container; its links point at the guide
    const href = (await card.getAttribute("href")) ?? (await card.locator(`a[href="${g.slug}/"]`).first().getAttribute("href"));
    expect(href, `link for ${g.slug}`).toBe(`${g.slug}/`);
  }
});

test("headline stats match the manifest", async ({ page }) => {
  const stats = page.locator("main dl").first();
  await expect(stats).toContainText(`${manifest.guides.length} guides`);
  const chapters = Number((await stats.locator("dt").nth(1).textContent())!.trim());
  expect(chapters).toBeGreaterThan(500);
});

test("search finds guides and chapters", async ({ page }) => {
  await page.locator("#search").fill("ownership");
  const results = page.locator("#results");
  await expect(results).toBeVisible();
  await expect(results).toContainText("Rust");
  await page.locator("#search").fill("transformer");
  await expect(results).toContainText("AI");
});

test("theme button cycles paper, light, night and restyles the page", async ({ page }) => {
  const html = page.locator("html");
  const bg = () => page.evaluate(() => getComputedStyle(document.body).backgroundColor);
  const seen = new Set<string>();
  for (let i = 0; i < 3; i++) {
    seen.add(`${await html.getAttribute("data-theme")}=${await bg()}`);
    await page.locator("#theme").click();
  }
  expect([...seen].map((s) => s.split("=")[0]).sort()).toEqual(["light", "night", "paper"]);
  expect(new Set([...seen].map((s) => s.split("=")[1])).size, "each theme has its own background").toBe(3);
});

test.describe("spotlight layout", () => {
  test.skip(layout !== "spotlight", "home_layout is not spotlight");
  const sp = manifest.spotlight!;
  const spotSection = manifest.sections!.find((s) => s.guides.includes(sp.guide))!;

  test("the spotlight comes first and features its guide", async ({ page }) => {
    const firstSection = page.locator("main > section[id]").first();
    await expect(firstSection).toHaveAttribute("id", spotSection.id);
    await expect(firstSection).toContainText(sp.eyebrow);
    await expect(firstSection.locator(`[data-guide="${sp.guide}"]`)).toContainText("Featured");
    await expect(firstSection.locator(`[data-guide="${sp.guide}"] h2`)).toHaveText(guide(sp.guide).title);
  });

  test("deep links point into the featured guide's chapter pages, in order", async ({ page }) => {
    const links = page.locator(`#${spotSection.id} a[href^="${sp.guide}/"][href*="#"]`);
    await expect(links).toHaveText(sp.links.map((l) => l.title));
    for (const [i, l] of sp.links.entries()) {
      // the chapter page that holds the anchor: <guide>/<chapter>/#<anchor>
      await expect(links.nth(i)).toHaveAttribute("href", new RegExp(`^${sp.guide}/[^/#]+/#${l.anchor}$`));
    }
  });

  test("the other guides of the AI section and the coming-next line are shown", async ({ page }) => {
    const box = page.locator(`#${spotSection.id}`);
    for (const slug of spotSection.guides.filter((s) => s !== sp.guide)) {
      await expect(box.locator(`a[data-guide="${slug}"]`)).toBeVisible();
    }
    for (const c of sp.coming) await expect(box).toContainText(c);
  });

  test("the engineering path ribbon follows manifest.path", async ({ page }) => {
    const steps = page.locator("#path ol a");
    await expect(steps).toHaveCount(manifest.path!.length);
    for (const [i, slug] of manifest.path!.entries()) {
      const g = guide(slug);
      await expect(steps.nth(i)).toHaveAttribute("href", `${slug}/`);
      await expect(steps.nth(i)).toContainText(g.short ?? g.title);
      await expect(steps.nth(i).locator("span").first()).toHaveText(g.step!);
    }
    const alongside = manifest.guides.filter((g) => g.step === "Alongside");
    await expect(page.locator("#path p a")).toHaveCount(alongside.length);
    for (const g of alongside) await expect(page.locator(`#path p a[href="${g.slug}/"]`)).toBeVisible();
  });

  test("topic sections and their cards follow manifest.sections", async ({ page }) => {
    const topics = manifest.sections!.filter((s) => s.id !== spotSection.id);
    const ids = await page.locator("main > section[id]").evaluateAll((els) => els.map((e) => e.id));
    expect(ids.filter((id) => topics.some((t) => t.id === id))).toEqual(topics.map((t) => t.id));
    for (const t of topics) {
      const sec = page.locator(`section#${t.id}`);
      await expect(sec.locator("h2")).toHaveText(t.title);
      const slugs = await sec.locator("a[data-guide]").evaluateAll((els) => els.map((e) => e.getAttribute("data-guide")));
      expect(slugs, `cards in ${t.id}`).toEqual(t.guides);
    }
  });

  test("section chips jump to sections that exist", async ({ page }) => {
    const chips = page.locator('nav[aria-label="Sections"] a');
    expect(await chips.count()).toBeGreaterThanOrEqual(manifest.sections!.length);
    for (const href of await chips.evaluateAll((els) => els.map((e) => e.getAttribute("href")))) {
      await expect(page.locator(href!), `target of ${href}`).toHaveCount(1);
    }
    await expect(page.locator("#series"), "old #series links still land somewhere").toHaveCount(1);
  });
});

test.describe("classic layout", () => {
  test.skip(layout !== "classic", "home_layout is not classic");

  test("one section per category, cards in manifest order", async ({ page }) => {
    for (const c of manifest.categories) {
      const sec = page.locator(`section#${c.id}`);
      await expect(sec.locator("h2")).toHaveText(c.title);
      const slugs = await sec.locator("a[data-guide]").evaluateAll((els) => els.map((e) => e.getAttribute("data-guide")));
      expect(slugs).toEqual(manifest.guides.filter((g) => g.category === c.id).map((g) => g.slug));
    }
  });
});
