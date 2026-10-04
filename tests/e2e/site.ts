// Shared test data, read from guides.json so new guides are tested automatically.
import { readFileSync } from "node:fs";
import { createHash } from "node:crypto";
import { expect, type Page } from "@playwright/test";

export interface Guide {
  slug: string;
  title: string;
  short?: string;
  icon: string;
  category: string;
  step?: string;
  source: string;
  pages?: { slug: string; source: string }[];
}

export interface Manifest {
  site: { title: string; home_layout?: string };
  categories: { id: string; title: string }[];
  guides: Guide[];
  sections?: { id: string; title: string; guides: string[] }[];
  path?: string[];
  spotlight?: { guide: string; eyebrow: string; links: { title: string; anchor: string }[]; coming: string[] };
}

export const manifest: Manifest = JSON.parse(readFileSync("guides.json", "utf8"));
export const layout = manifest.site.home_layout === "spotlight" ? "spotlight" : "classic";
export const guide = (slug: string) => manifest.guides.find((g) => g.slug === slug)!;

/** Every rendered document: each guide's main page plus its sub-pages. */
export const docs = manifest.guides.flatMap((g) => [
  { guide: g, url: `${g.slug}/`, main: true },
  ...(g.pages ?? []).map((p) => ({ guide: g, url: `${g.slug}/${p.slug}/`, main: false })),
]);

/** The ?v= content hash the build should put on an asset (see assetVersion in build.go). */
export function assetHash(file: string): string {
  return createHash("sha256").update(readFileSync(`web/static/${file}`)).digest("hex").slice(0, 10);
}

/** Collects console errors, page errors, and failed same-origin requests while a test runs. */
export function watchErrors(page: Page): () => void {
  const problems: string[] = [];
  page.on("console", (m) => { if (m.type() === "error") problems.push(`console: ${m.text()}`); });
  page.on("pageerror", (e) => problems.push(`page error: ${e.message}`));
  page.on("requestfailed", (r) => { if (r.url().startsWith("http://127.0.0.1")) problems.push(`failed: ${r.url()}`); });
  page.on("response", (r) => {
    if (r.url().startsWith("http://127.0.0.1") && r.status() >= 400) problems.push(`HTTP ${r.status()}: ${r.url()}`);
  });
  return () => expect(problems, "browser errors").toEqual([]);
}

/** WCAG contrast ratio between two CSS rgb()/rgba() colours. */
export function contrast(a: string, b: string): number {
  const lum = (c: string) => {
    const [r, g, bl] = c.match(/[\d.]+/g)!.slice(0, 3).map(Number).map((v) => {
      const s = v / 255;
      return s <= 0.03928 ? s / 12.92 : ((s + 0.055) / 1.055) ** 2.4;
    });
    return 0.2126 * r + 0.7152 * g + 0.0722 * bl;
  };
  const [hi, lo] = [lum(a), lum(b)].sort((x, y) => y - x);
  return (hi + 0.05) / (lo + 0.05);
}
