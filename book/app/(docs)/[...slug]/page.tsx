import {
  DocsBody,
  DocsPage,
  DocsTitle,
} from "fumadocs-ui/layouts/docs/page";
import defaultMdxComponents, { createRelativeLink } from "fumadocs-ui/mdx";
import type { MDXComponents } from "mdx/types";
import { notFound } from "next/navigation";
import { source } from "@/lib/source";

interface PageProps {
  params: Promise<{ slug: string[] }>;
}

export default async function Page(props: PageProps) {
  const { slug } = await props.params;
  const page = source.getPage(slug);

  if (!page) notFound();

  const MDX = page.data.body;
  const components: MDXComponents = {
    ...defaultMdxComponents,
    a: createRelativeLink(source, page),
  };

  return (
    <DocsPage toc={page.data.toc}>
      <DocsTitle>{page.data.title}</DocsTitle>
      <DocsBody>
        <MDX components={components} />
      </DocsBody>
    </DocsPage>
  );
}

export function generateStaticParams() {
  return source.generateParams();
}

export const dynamicParams = false;
