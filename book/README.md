# RespKit documentation site

The published documentation is available at [zihxs.github.io/RespKit](https://zihxs.github.io/RespKit/). Its pages are in `content/docs/` and use MDX.

## Run locally

```bash
npm ci
npm run dev
```

## Build the static site

```bash
npm run build
```

The generated site is written to `out/`. GitHub Pages builds and publishes it from the `main` branch workflow.

## Branding

The project title is **RespKit — Consistent HTTP Responses for Go**. The site title, tagline, and description are defined in `lib/branding.ts`. Keep the README heading aligned with that title when changing the branding.

The navigation logo uses `assets/respkit.png` from the repository root. Update that file to change the logo in both the README and the documentation site.
