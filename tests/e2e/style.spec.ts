// Style sanity: the right stylesheet is loaded (cache-busted by content hash),
// the design tokens apply, and every theme keeps text readable.
import { test, expect } from "@playwright/test";
import { assetHash, contrast, layout, manifest } from "./site";

const THEME_BG: Record<string, string> = {           // --bg in web/styles/input.css
  paper: "rgb(247, 242, 231)",
  light: "rgb(253, 252, 250)",
  night: "rgb(21, 23, 28)",
};

test("assets are versioned with their current content hash", async ({ page, request }) => {
  await page.goto("./");
  const css = await page.locator('link[rel="stylesheet"][href*="app.css"]').getAttribute("href");
  expect(css, "stale app.css hash: was the site rebuilt after a CSS change?").toContain(`app.css?v=${assetHash("app.css")}`);
  const js = await page.locator('script[src*="home.js"]').getAttribute("src");
  expect(js).toContain(`home.js?v=${assetHash("home.js")}`);
  expect((await request.get(css!)).status()).toBe(200);

  await page.goto("rust/");
  const reader = await page.locator('script[src*="reader.js"]').getAttribute("src");
  expect(reader).toContain(`reader.js?v=${assetHash("reader.js")}`);
});

test("the compiled stylesheet is applied, not just linked", async ({ page }) => {
  await page.goto("./");
  const s = await page.evaluate(() => {
    const card = document.querySelector("a[data-guide]")!;
    const cs = getComputedStyle(card);
    return { radius: cs.borderTopLeftRadius, border: cs.borderTopStyle, font: getComputedStyle(document.body).fontFamily };
  });
  expect(s.radius, "cards are rounded-2xl").toBe("16px");
  expect(s.border).toBe("solid");
  expect(s.font).toContain("ui-sans-serif");
});

for (const theme of Object.keys(THEME_BG)) {
  test(`${theme} theme: tokens apply and text stays readable`, async ({ page }) => {
    await page.addInitScript((t) => {
      // home.js restores the theme from saved preferences; start each run from this theme.
      document.documentElement.dataset.theme = t;
      localStorage.setItem("fl:prefs", JSON.stringify({ theme: t }));
    }, theme);
    await page.goto("./");
    await page.evaluate((t) => { document.documentElement.dataset.theme = t; }, theme);
    const c = await page.evaluate(() => {
      const h1 = document.querySelector("h1")!;
      const muted = document.querySelector("main p")!;
      return {
        bg: getComputedStyle(document.body).backgroundColor,
        h1: getComputedStyle(h1).color,
        muted: getComputedStyle(muted).color,
      };
    });
    expect(c.bg).toBe(THEME_BG[theme]);
    expect(contrast(c.h1, c.bg), "headline contrast (WCAG AA large: 3)").toBeGreaterThanOrEqual(4.5);
    expect(contrast(c.muted, c.bg), "body text contrast (WCAG AA: 4.5)").toBeGreaterThanOrEqual(4.5);
  });
}

test.describe("spotlight styling", () => {
  test.skip(layout !== "spotlight", "home_layout is not spotlight");

  test("spotlight panel has its gradient, accent border and large radius", async ({ page }) => {
    await page.goto("./");
    const id = manifest.sections!.find((s) => s.guides.includes(manifest.spotlight!.guide))!.id;
    const s = await page.locator(`#${id} > div`).first().evaluate((el) => {
      const cs = getComputedStyle(el);
      return { image: cs.backgroundImage, radius: cs.borderTopLeftRadius, border: cs.borderTopColor };
    });
    expect(s.image).toContain("linear-gradient");
    expect(s.radius).toBe("24px");
    expect(s.border).not.toBe("rgba(0, 0, 0, 0)");
  });

  test("Featured badge and Start reading use the accent colour with readable text", async ({ page }) => {
    await page.goto("./");
    for (const name of ["Featured", "Start reading"]) {
      const el = page.getByText(name, { exact: false }).first();
      const { bg, fg } = await el.evaluate((e) => ({ bg: getComputedStyle(e).backgroundColor, fg: getComputedStyle(e).color }));
      expect(bg, `${name} background`).not.toBe("rgba(0, 0, 0, 0)");
      expect(contrast(fg, bg), `${name} contrast`).toBeGreaterThanOrEqual(4.5);
    }
  });
});
