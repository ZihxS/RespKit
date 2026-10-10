import { existsSync, readdirSync, readFileSync } from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const root = path.dirname(fileURLToPath(import.meta.url));
const contentDir = path.join(root, "..", "content", "docs");
const publicDir = path.join(root, "..", "public");

const files = readdirSync(contentDir, { withFileTypes: true })
  .filter((entry) => entry.isFile() && entry.name.endsWith(".mdx"))
  .map((entry) => entry.name);
const routes = new Set(files.map((file) => `/${file.replace(/\.mdx$/, "")}`));
const missing = [];

for (const filename of files) {
  const content = readFileSync(path.join(contentDir, filename), "utf8");

  for (const match of content.matchAll(/\[[^\]]*\]\(([^)\s]+)\)/g)) {
    const href = match[1];
    if (!href.startsWith("/") || href.startsWith("//")) continue;

    const route = href.split(/[?#]/, 1)[0].replace(/\/$/, "") || "/";
    const publicPath = path.join(publicDir, route.replace(/^\//, ""));
    if (!routes.has(route) && !existsSync(publicPath)) {
      missing.push(`${filename}: ${href}`);
    }
  }
}

if (missing.length > 0) {
  console.error(`Unknown internal documentation links:\n${missing.join("\n")}`);
  process.exit(1);
}

console.log(`Checked internal links in ${files.length} documentation pages.`);
