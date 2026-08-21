import type { Metadata } from "next";
import { Fira_Code } from "next/font/google";
import "./globals.css";
import { synclink, type PlatformSettings } from "@/lib/api";

const fira = Fira_Code({
  variable: "--font-fira",
  subsets: ["latin"],
  weight: ["400", "500", "600"],
});

async function loadSettings(): Promise<PlatformSettings | null> {
  try {
    return await synclink.publicSettings();
  } catch {
    return null;
  }
}

export async function generateMetadata(): Promise<Metadata> {
  const s = await loadSettings();
  const title = s?.metaTitle || s?.siteName || "SyncLink";
  const description = s?.metaDescription || s?.tagline || "One page. Every link.";
  return {
    title,
    description,
    icons: s?.favicon ? { icon: s.favicon } : undefined,
    openGraph: {
      title,
      description,
      images: s?.ogImage ? [{ url: s.ogImage }] : undefined,
    },
    other: s?.themeColor ? { "theme-color": s.themeColor } : undefined,
  };
}

export default async function RootLayout({ children }: LayoutProps<"/">) {
  const s = await loadSettings();
  const theme = s?.themeColor || "#111111";
  return (
    <html lang="en" className={`${fira.variable} h-full antialiased`}>
      <head>
        <meta name="theme-color" content={theme} />
      </head>
      <body className="page-enter min-h-full bg-[#faf9f7] font-sans text-neutral-950">{children}</body>
    </html>
  );
}
