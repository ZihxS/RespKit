import Link from "next/link";

const sections = [
  {
    title: "Mulai memakai RespKit",
    description: "Pasang package dan kirim response pertama dari handler Go.",
    href: "/quickstart",
  },
  {
    title: "Integrasi framework",
    description: "Contoh untuk net/http, Gin, Echo, Fiber, Chi, mux, dan Beego.",
    href: "/frameworks",
  },
  {
    title: "Pemetaan error",
    description: "Gunakan kode error yang stabil dan aman untuk client API.",
    href: "/errors",
  },
];

export default function Home() {
  return (
    <main className="home-hero flex min-h-screen flex-col">
      <header className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 py-6">
        <Link className="text-lg font-semibold tracking-tight" href="/">
          RespKit
        </Link>
        <Link
          className="rounded-lg border border-fd-border px-4 py-2 text-sm font-medium transition-colors hover:bg-fd-muted"
          href="/docs"
        >
          Dokumentasi
        </Link>
      </header>

      <section className="mx-auto flex w-full max-w-6xl flex-1 flex-col justify-center px-6 pb-20 pt-10">
        <div className="max-w-3xl">
          <p className="mb-5 inline-flex rounded-full border border-fd-border bg-fd-card px-3 py-1 text-sm text-fd-muted-foreground">
            Pustaka response HTTP untuk Go
          </p>
          <h1 className="text-5xl font-semibold tracking-tight sm:text-6xl">
            Response API yang konsisten, dari satu package.
          </h1>
          <p className="mt-6 max-w-2xl text-lg leading-8 text-fd-muted-foreground">
            Format response, pemetaan error, dan terjemahan yang aman untuk
            aplikasi Go Anda.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link
              className="rounded-lg bg-fd-primary px-5 py-3 text-sm font-semibold text-white transition-opacity hover:opacity-90"
              href="/docs"
            >
              Buka dokumentasi
            </Link>
            <a
              className="rounded-lg border border-fd-border bg-fd-card px-5 py-3 text-sm font-semibold transition-colors hover:bg-fd-muted"
              href="https://github.com/ZihxS/RespKit"
              rel="noreferrer"
              target="_blank"
            >
              Lihat di GitHub
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

      <footer className="mx-auto w-full max-w-6xl px-6 py-6 text-sm text-fd-muted-foreground">
        RespKit · Dokumentasi berbahasa Indonesia
      </footer>
    </main>
  );
}
