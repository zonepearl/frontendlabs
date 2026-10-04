// Visual regression (@visual): pixel comparisons against approved screenshots.
// Fonts differ between macOS and Linux, so baselines are kept per platform
// and CI skips this file (--grep-invert @visual); the pre-commit hook runs it.
// After an intended design change: npx playwright test visual --update-snapshots
import { test, expect, type Page } from "@playwright/test";

async function settle(page: Page) {
  await page.evaluate(() => document.fonts.ready);
  await page.waitForLoadState("networkidle");
}

// Parts that legitimately change between builds or with local reading history.
const volatile = (page: Page) => [page.locator("footer"), page.locator("#continue"), page.locator(".resume-badge")];

test.describe("@visual", () => {
  test("home page", async ({ page }) => {
    await page.goto("./");
    await settle(page);
    await expect(page).toHaveScreenshot("home.png", { fullPage: true, mask: volatile(page) });
  });

  test("home page, night theme", async ({ page }) => {
    await page.addInitScript(() => localStorage.setItem("fl:prefs", JSON.stringify({ theme: "night" })));
    await page.goto("./");
    await settle(page);
    await expect(page).toHaveScreenshot("home-night.png", { mask: volatile(page) });
  });

  test("guide page: title and contents", async ({ page }) => {
    await page.goto("rust/");
    await settle(page);
    await expect(page).toHaveScreenshot("reader-rust.png", { mask: volatile(page) });
  });
});
