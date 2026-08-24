import type { MetadataRoute } from "next";
import { synclink } from "@/lib/api";

function siteBase(): string {
  const raw =
    process.env.NEXT_PUBLIC_SITE_URL ||
    process.env.VERCEL_PROJECT_PRODUCTION_URL ||
    "https://synclink-mocha.vercel.app";
  const trimmed = raw.trim().replace(/\/+$/, "");
  if (/^https?:\/\//i.test(trimmed)) return trimmed;
  return `https://${trimmed}`;
}

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const base = siteBase();
  const lastModified = new Date();
  const entries: MetadataRoute.Sitemap = [
    { url: `${base}/`, lastModified },
    { url: `${base}/about`, lastModified },
  ];

  try {
    const settings = await synclink.publicSettings();
    const slug = settings.demoSlug?.trim();
    if (slug) {
      entries.push({ url: `${base}/${slug}`, lastModified });
    }
  } catch {
    // still emit / and /about if settings fail
  }

  return entries;
}
