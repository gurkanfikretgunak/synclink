import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { LockedPublicPage } from "@/components/locked-public-page";
import { PublicPageBody } from "@/components/public-page-body";
import { synclink, type PublicPage } from "@/lib/api";

function ProfileJsonLd({ page }: { page: PublicPage }) {
  const sameAs = (page.socials || []).map((item) => item.url).filter(Boolean);
  const image = page.avatarUrl || page.coverUrl || undefined;
  const person: Record<string, unknown> = {
    "@type": "Person",
    name: page.displayName.trim() || page.slug,
    alternateName: page.slug,
    identifier: page.slug,
  };
  if (page.bio.trim()) person.description = page.bio.trim();
  if (image) person.image = image;
  if (sameAs.length) person.sameAs = sameAs;
  const data: Record<string, unknown> = {
    "@context": "https://schema.org",
    "@type": "ProfilePage",
    mainEntity: person,
  };
  if (page.publishedAt) data.dateCreated = page.publishedAt;
  return <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(data) }} />;
}

export async function generateMetadata({ params }: PageProps<"/[slug]">): Promise<Metadata> {
  const { slug } = await params;
  try {
    const page = await synclink.getPublicPage(slug);
    const title = page.displayName.trim() || slug;
    const description = page.bio.trim() || undefined;
    const image = page.avatarUrl || page.coverUrl || undefined;
    const images = image ? [{ url: image }] : undefined;
    const icon = [page.avatarUrl, page.coverUrl]
      .map((url) => (typeof url === "string" ? url.trim() : ""))
      .find((url) => url.startsWith("https://") || url.startsWith("http://"));
    const themeColor = page.accentColor?.trim() || undefined;
    return {
      title,
      description,
      applicationName: title,
      ...(themeColor ? { themeColor } : {}),
      ...(icon ? { icons: { icon, apple: icon } } : {}),
      appleWebApp: { capable: true, title, statusBarStyle: "default" },
      alternates: { canonical: `/${slug}` },
      openGraph: { title, description, type: "profile", images, url: `/${slug}` },
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
    return (
      <>
        <ProfileJsonLd page={page} />
        <PublicPageBody page={page} />
      </>
    );
  } catch (err) {
    const message = err instanceof Error ? err.message : "";
    if (message === "locked") {
      return <LockedPublicPage slug={slug} />;
    }
    notFound();
  }
}
