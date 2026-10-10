import { RootProvider } from "fumadocs-ui/provider/next";
import type { Metadata } from "next";
import type { ReactNode } from "react";
import "./global.css";

const siteURL = "https://zihxs.github.io/RespKit";

export const metadata: Metadata = {
  metadataBase: new URL(siteURL),
  title: {
    default: "RespKit Documentation",
    template: "%s | RespKit",
  },
  description: "Dokumentasi RespKit, pustaka response HTTP untuk Go.",
  openGraph: {
    title: "RespKit Documentation",
    description: "Dokumentasi RespKit, pustaka response HTTP untuk Go.",
    url: siteURL,
    siteName: "RespKit",
    type: "website",
  },
};

export default function RootLayout({ children }: { children: ReactNode }) {
  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  return (
    <html lang="id" suppressHydrationWarning>
      <body className="flex min-h-screen flex-col font-sans">
        <RootProvider
          search={{
            options: {
              api: `${basePath}/api/search`,
              type: "static",
            },
          }}
        >
          {children}
        </RootProvider>
      </body>
    </html>
  );
}
