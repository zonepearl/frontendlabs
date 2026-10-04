// E2E suite for the generated site. It builds the site from docs/ into
// .e2e-dist/, serves it as plain static files (like GitHub Pages), and tests
// it in Chromium at desktop and phone sizes.
//
//   npx playwright test                    everything (what the pre-commit hook runs)
//   npx playwright test --grep-invert @visual   skip screenshot comparisons (what CI runs)
//   npx playwright test --update-snapshots      accept intended design/content changes
import { defineConfig, devices } from "@playwright/test";

const PORT = 4173;

export default defineConfig({
  testDir: "tests/e2e",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? [["list"], ["github"]] : [["list"]],
  use: {
    baseURL: `http://127.0.0.1:${PORT}/`,
    trace: "retain-on-failure",
  },
  // Text snapshots (content: tables of contents) are the same on every OS, so
  // their names carry no platform and CI can compare against them.
  snapshotPathTemplate: "tests/e2e/__snapshots__/{testFileName}/{arg}{ext}",
  expect: {
    // Screenshots depend on the OS's fonts, so each platform keeps its own baselines.
    toHaveScreenshot: {
      pathTemplate: "tests/e2e/__screenshots__/{platform}/{projectName}/{arg}{ext}",
      maxDiffPixelRatio: 0.01,
      animations: "disabled",
    },
  },
  projects: [
    {
      name: "desktop",
      use: { ...devices["Desktop Chrome"], viewport: { width: 1360, height: 900 } },
    },
    {
      // Phone layout: only the specs that are about layout and the home page.
      name: "mobile",
      use: { ...devices["Pixel 7"] },
      testMatch: /(home|layout|visual)\.spec\.ts/,
    },
  ],
  webServer: {
    command: `go run ./cmd/guides build -out .e2e-dist && python3 -m http.server ${PORT} --bind 127.0.0.1 -d .e2e-dist`,
    url: `http://127.0.0.1:${PORT}/`,
    reuseExistingServer: false,
    timeout: 120_000,
    stdout: "ignore",
  },
});
