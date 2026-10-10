import { createMDX } from "fumadocs-mdx/next";

const basePath = "/RespKit";
const isGitHubPages = process.env.RESPKIT_DOCS_GITHUB_PAGES === "true";

/** @type {import('next').NextConfig} */
const config = {
  assetPrefix: isGitHubPages ? "https://zihxs.github.io/RespKit" : undefined,
  basePath,
  env: {
    NEXT_PUBLIC_BASE_PATH: basePath,
  },
  images: { unoptimized: true },
  output: "export",
  reactStrictMode: true,
  trailingSlash: true,
};

export default createMDX()(config);
