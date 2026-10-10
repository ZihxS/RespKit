import { RootProvider } from "fumadocs-ui/provider/next";
import type { Metadata } from "next";
import type { ReactNode } from "react";
import { brandDescription, brandName, brandTitle } from "@/lib/branding";
import "./global.css";

const siteURL = "https://zihxs.github.io/RespKit";

export const metadata: Metadata = {
  metadataBase: new URL(siteURL),
  title: {
    default: brandTitle,
    template: `%s | ${brandName}`,
  },
  description: brandDescription,
  openGraph: {
    title: brandTitle,
    description: brandDescription,
    url: siteURL,
    siteName: brandName,
    type: "website",
  },
};

export default function RootLayout({ children }: { children: ReactNode }) {
  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  return (
    <html lang="en" suppressHydrationWarning>
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
