import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { LockedPublicPage } from "@/components/locked-public-page";
import { PublicPageBody } from "@/components/public-page-body";
import { synclink } from "@/lib/api";

export async function generateMetadata({ params }: PageProps<"/[slug]">): Promise<Metadata> {
  const { slug } = await params;
  try {
    const page = await synclink.getPublicPage(slug);
    const title = page.displayName.trim() || slug;
    const description = page.bio.trim() || undefined;
    const image = page.avatarUrl || page.coverUrl || undefined;
    const images = image ? [{ url: image }] : undefined;
    return {
      title,
      description,
      openGraph: { title, description, type: "profile", images },
      twitter: { card: "summary", title, description, images: image ? [image] : undefined },
    };
  } catch (err) {
    const message = err instanceof Error ? err.message : "";
    if (message === "locked") {
      return { title: slug, robots: { index: false, follow: false } };
    }
    return {};
  }
}

export default async function PublicPage({ params }: PageProps<"/[slug]">) {
  const { slug } = await params;
  try {
    const page = await synclink.getPublicPage(slug);
    return <PublicPageBody page={page} />;
  } catch (err) {
    const message = err instanceof Error ? err.message : "";
    if (message === "locked") {
      return <LockedPublicPage slug={slug} />;
    }
    notFound();
  }
}
