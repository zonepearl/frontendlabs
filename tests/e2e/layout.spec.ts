// Layout sanity measured from element boxes, so it holds on any OS and font:
// no sideways scrolling, grids and columns where the design expects them,
// sticky bars that stay put. Runs on desktop (1360px) and phone (Pixel 7).
import { test, expect, type Page } from "@playwright/test";
import { docs, layout, manifest } from "./site";

const box = async (page: Page, selector: string) => (await page.locator(selector).first().boundingBox())!;
const isPhone = (page: Page) => page.viewportSize()!.width < 640;

async function noSidewaysScroll(page: Page) {
  const { scroll, client } = await page.evaluate(() => ({
    scroll: document.documentElement.scrollWidth,
    client: document.documentElement.clientWidth,
  }));
  expect(scroll, `page is ${scroll - client}px wider than the viewport`).toBeLessThanOrEqual(client + 1);
}

test("home page never scrolls sideways", async ({ page }) => {
  await page.goto("./");
  await noSidewaysScroll(page);
});

// Each guide's one-page edition holds every chapter, so one load per guide covers all its
// content; long code lines must scroll inside their block, not widen the page.
for (const d of docs.filter((d) => d.main)) {
  test(`${d.url} never scrolls sideways`, async ({ page, request }) => {
    const all = (await request.get(`${d.url}all/`)).ok();
    await page.goto(all ? `${d.url}all/` : d.url);
    await noSidewaysScroll(page);
    if (all) {
      await page.goto(d.url); // and the landing page with its contents list
      await noSidewaysScroll(page);
    }
  });
}

test("a chapter page fits the screen and shows its breadcrumb", async ({ page }) => {
  await page.goto("rust/6-ownership-moves-copies-and-drop/");
  await noSidewaysScroll(page);
  await expect(page.locator('nav[aria-label="Breadcrumb"]')).toBeVisible();
  await expect(page.locator("h1.chapter-title")).toBeVisible();
});

test("header and section chips stay visible while scrolling", async ({ page }) => {
  await page.goto("./");
  await page.mouse.wheel(0, 2500);
  await page.waitForTimeout(200);
  expect((await box(page, "header")).y).toBeLessThanOrEqual(1);
  if (layout === "spotlight") {
    const nav = await box(page, 'nav[aria-label="Sections"]');
    expect(nav.y, "section chips stick below the header").toBeLessThan(120);
  }
});

test.describe("spotlight geometry", () => {
  test.skip(layout !== "spotlight", "home_layout is not spotlight");
  const ai = () => manifest.sections!.find((s) => s.guides.includes(manifest.spotlight!.guide))!.id;

  test("featured card and side column: side by side on desktop, stacked on phones", async ({ page }) => {
    await page.goto("./");
    const featured = await box(page, `#${ai()} [data-guide="${manifest.spotlight!.guide}"]`);
    const extra = await box(page, `#${ai()} a[data-guide]`);
    if (isPhone(page)) {
      expect(extra.y, "side column goes below the featured card").toBeGreaterThan(featured.y + featured.height - 1);
    } else {
      expect(extra.x, "side column is to the right").toBeGreaterThan(featured.x + featured.width - 1);
      expect(Math.abs(extra.y - featured.y), "tops are aligned").toBeLessThan(2);
      expect(featured.width, "featured card is the wider one").toBeGreaterThan(extra.width);
    }
  });

  test("the course path is one even row on desktop", async ({ page }) => {
    test.skip(isPhone(page), "the path wraps by design on phones");
    await page.goto("./");
    const boxes = await page.locator("#path ol > li").evaluateAll((els) =>
      els.map((e) => { const r = e.getBoundingClientRect(); return { y: r.top, w: r.width }; }));
    expect(new Set(boxes.map((b) => Math.round(b.y))).size, "all steps on one row").toBe(1);
    const widths = boxes.map((b) => b.w);
    expect(Math.max(...widths) - Math.min(...widths), "steps are equal width").toBeLessThan(2);
  });

  test("topic cards: three columns on desktop, one on phones", async ({ page }) => {
    await page.goto("./");
    const sec = manifest.sections!.find((s) => s.guides.length >= 3 && s.id !== ai())!;
    const xs = await page.locator(`section#${sec.id} a[data-guide]`).evaluateAll((els) =>
      els.slice(0, 3).map((e) => Math.round(e.getBoundingClientRect().left)));
    expect(new Set(xs).size).toBe(isPhone(page) ? 1 : 3);
  });
});

test("a guide page shows the book layout: contents beside the text on desktop", async ({ page }) => {
  await page.goto("rust/");
  const content = await box(page, "#content");
  expect(content.width).toBeGreaterThan(isPhone(page) ? 300 : 500);
  if (!isPhone(page)) {
    await expect(page.locator("#toc")).toBeVisible();
    const toc = await box(page, "#toc");
    expect(toc.x + toc.width, "contents sit to the left of the text").toBeLessThanOrEqual(content.x + 1);
  }
});
