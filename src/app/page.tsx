import Image from "next/image";
import Link from "next/link";
import { buttonVariants } from "@/components/ui/button";
import { SiteNav } from "@/components/site-nav";
import { synclink, type PlatformSettings } from "@/lib/api";

const features = [
  { label: "Compose", body: "Name, bio, a short list." },
  { label: "Publish", body: "Live at /your-slug." },
  { label: "Share", body: "One URL. QR, vCard, copy." },
];

async function loadSettings(): Promise<Partial<PlatformSettings>> {
  try {
    return await synclink.publicSettings();
  } catch {
    return {};
  }
}

function monogram(title: string) {
  const words = title.replace(/\n/g, " ").trim().split(/\s+/).filter(Boolean);
  if (words.length === 0) return "S";
  if (words.length === 1) return words[0].slice(0, 2).toUpperCase();
  return (words[0][0] + words[1][0]).toUpperCase();
}

function DeviceFrame({
  src,
  alt,
  title,
}: {
  src: string;
  alt: string;
  title: string;
}) {
  const remote = src.startsWith("http://") || src.startsWith("https://");
  return (
    <div className="mx-auto w-full max-w-sm">
      <div className="rounded-[2rem] border border-neutral-300 bg-neutral-900 p-1.5 shadow-[0_20px_60px_-32px_rgba(0,0,0,0.35)]">
        <div className="overflow-hidden rounded-[1.55rem] bg-[#faf9f7]">
          {src ? (
            remote ? (
              // remote admin-driven art — next/image needs a host allowlist
              // eslint-disable-next-line @next/next/no-img-element
              <img src={src} alt={alt} className="aspect-[4/5] h-auto w-full object-cover" />
            ) : (
              <Image
                src={src}
                alt={alt}
                width={900}
                height={1125}
                className="aspect-[4/5] h-auto w-full object-cover"
                priority
              />
            )
          ) : (
            <div className="flex aspect-[4/5] items-center justify-center">
              <span className="text-5xl font-medium tracking-tight text-neutral-800">{monogram(title)}</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default async function Home() {
  const settings = await loadSettings();
  const title = settings.heroTitle || "One page.\nEvery link.";
  const subtitle =
    settings.heroSubtitle ||
    settings.tagline ||
    "A quieter public page. White space, type, and a few stills. Edit from the dashboard.";
  const cta = settings.heroCta || "Create your page";
  const ctaHref = settings.heroCtaHref || "/dashboard";
  const hero = settings.heroImage || "";
  const demo = settings.demoSlug || "gurkan";
  const lines = title.split("\n");

  return (
    <main className="page-enter min-h-full bg-[#faf9f7]">
      <SiteNav settings={settings} />

      <section className="mx-auto grid w-full max-w-5xl items-center gap-12 px-6 pb-16 pt-8 md:grid-cols-2 md:pb-20 md:pt-10">
        <div className="space-y-6">
          <h1 className="text-5xl font-medium tracking-tight text-neutral-950 md:text-6xl">
            {lines.map((line) => (
              <span key={line} className="block">{line}</span>
            ))}
          </h1>
          <p className="max-w-sm text-base leading-relaxed text-neutral-600">{subtitle}</p>
          <div className="flex flex-wrap gap-3">
            <Link href={ctaHref} className={buttonVariants({ size: "lg" })}>{cta}</Link>
            <Link href={`/${demo}`} className={buttonVariants({ variant: "ghost", size: "lg" })}>
              See /{demo}
            </Link>
          </div>
          <Link
            href={`/${demo}`}
            className="inline-flex font-mono text-sm text-neutral-500 hover:text-neutral-800"
          >
            /{demo}
          </Link>
        </div>
        <DeviceFrame src={hero} alt="" title={title} />
      </section>

      <section className="mx-auto flex w-full max-w-5xl flex-col gap-8 px-6 pb-16 md:flex-row md:gap-12">
        {features.map((item) => (
          <div key={item.label} className="flex-1">
            <p className="text-[10px] uppercase tracking-[0.28em] text-neutral-400">{item.label}</p>
            <p className="mt-2 text-sm leading-relaxed text-neutral-700">{item.body}</p>
          </div>
        ))}
      </section>

      <footer className="mx-auto w-full max-w-5xl border-t border-neutral-200/80 px-6 py-6">
        <p className="text-sm text-neutral-500">
          <Link href="/dashboard" className="hover:text-neutral-800">Edit from the dashboard</Link>
          <span className="mx-2 text-neutral-300">·</span>
          <Link href="/about" className="hover:text-neutral-800">About</Link>
        </p>
      </footer>
    </main>
  );
}
