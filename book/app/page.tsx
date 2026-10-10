import Link from "next/link";

import { Brand } from "@/components/brand";

const sections = [
  {
    title: "Get started with RespKit",
    description: "Install the package and send your first response from a Go handler.",
    href: "/quickstart",
  },
  {
    title: "Framework integrations",
    description: "Examples for net/http, Gin, Echo, Fiber, Chi, mux, and Beego.",
    href: "/frameworks",
  },
  {
    title: "Error mapping",
    description: "Use stable error codes and safe messages in API responses.",
    href: "/errors",
  },
];

export default function Home() {
  return (
    <main className="home-hero flex min-h-screen flex-col">
      <header className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 py-6">
        <Link className="inline-flex items-center" href="/">
          <Brand />
        </Link>
        <Link
          className="rounded-lg border border-fd-border px-4 py-2 text-sm font-medium transition-colors hover:bg-fd-muted"
          href="/docs"
        >
          Documentation
        </Link>
      </header>

      <section className="mx-auto flex w-full max-w-6xl flex-1 flex-col justify-center px-6 pb-20 pt-10">
        <div className="max-w-3xl">
          <p className="mb-5 inline-flex rounded-full border border-fd-border bg-fd-card px-3 py-1 text-sm text-fd-muted-foreground">
            HTTP response library for Go
          </p>
          <h1 className="text-5xl font-semibold tracking-tight sm:text-6xl">
            Consistent HTTP responses from one package.
          </h1>
          <p className="mt-6 max-w-2xl text-lg leading-8 text-fd-muted-foreground">
            Format responses, map errors, and return localized messages in Go
            applications.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link
              className="rounded-lg bg-fd-primary px-5 py-3 text-sm font-semibold text-white transition-opacity hover:opacity-90"
              href="/docs"
            >
              Read the docs
            </Link>
            <a
              className="rounded-lg border border-fd-border bg-fd-card px-5 py-3 text-sm font-semibold transition-colors hover:bg-fd-muted"
              href="https://github.com/ZihxS/RespKit"
              rel="noreferrer"
              target="_blank"
            >
              View on GitHub
            </a>
          </div>
        </div>

        <div className="mt-16 grid gap-4 md:grid-cols-3">
          {sections.map((section) => (
            <Link
              className="home-card rounded-xl border border-fd-border bg-fd-card p-5"
              href={section.href}
              key={section.title}
            >
              <h2 className="font-semibold">{section.title}</h2>
              <p className="mt-2 text-sm leading-6 text-fd-muted-foreground">
                {section.description}
              </p>
            </Link>
          ))}
        </div>
      </section>
    </main>
  );
}
