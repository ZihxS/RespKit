import type { BaseLayoutProps } from "fumadocs-ui/layouts/shared";

import { Brand } from "@/components/brand";

export function baseOptions(): BaseLayoutProps {
  return {
    githubUrl: "https://github.com/ZihxS/RespKit",
    nav: {
      title: <Brand />,
    },
  };
}
