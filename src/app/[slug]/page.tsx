import { notFound } from "next/navigation";
import { LockedPublicPage } from "@/components/locked-public-page";
import { PublicPageBody } from "@/components/public-page-body";
import { synclink } from "@/lib/api";

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
