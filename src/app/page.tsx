import Image from "next/image";
import Link from "next/link";
import { buttonVariants } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { SiteNav } from "@/components/site-nav";
import { synclink, type PlatformSettings } from "@/lib/api";

const stations = [
  { n: "01", title: "Compose", body: "Name, bio, a short list. Nothing extra." },
  { n: "02", title: "Publish", body: "Your page lives at /your-slug. Fast, public, text-first." },
  { n: "03", title: "Share", body: "One URL. Every link you want people to keep." },
];

async function loadSettings(): Promise<Partial<PlatformSettings>> {
  try {
    return await synclink.publicSettings();
  } catch {
    return {};
  }
}

function HeroImage({ src, alt }: { src: string; alt: string }) {
  if (src.startsWith("http://") || src.startsWith("https://")) {
    return (
      // remote admin-driven art — next/image needs a host allowlist
      // eslint-disable-next-line @next/next/no-img-element
      <img src={src} alt={alt} className="h-auto w-full object-cover" />
    );
  }
  return (
    <Image src={src} alt={alt} width={1600} height={1200} className="h-auto w-full object-cover" priority />
  );
}

export default async function Home() {
  const settings = await loadSettings();
  const title = settings.heroTitle || "One page.\nEvery link.";
  const subtitle = settings.heroSubtitle || settings.tagline || "A quieter public page. White space, type, and a few stills. Edit from the dashboard.";
  const cta = settings.heroCta || "Create your page";
  const ctaHref = settings.heroCtaHref || "/dashboard";
  const hero = settings.heroImage || "/stations/hero.png";
  const demo = settings.demoSlug || "gurkan";
  const lines = title.split("\n");

  return (
    <main className="page-enter min-h-full bg-[#faf9f7]">
      <SiteNav settings={settings} />

      <section className="mx-auto grid w-full max-w-5xl items-center gap-12 px-6 pb-16 pt-4 md:grid-cols-2 md:pb-24">
        <div className="space-y-6">
          <p className="text-xs uppercase tracking-[0.28em] text-neutral-500">Station 01 — Arrival</p>
          <h1 className="text-5xl font-medium tracking-tight text-neutral-950 md:text-6xl">
            {lines.map((line) => (
              <span key={line} className="block">{line}</span>
            ))}
          </h1>
          <p className="max-w-sm text-base leading-relaxed text-neutral-600">{subtitle}</p>
          <div className="flex flex-wrap gap-3">
            <Link href={ctaHref} className={buttonVariants({ size: "lg" })}>{cta}</Link>
            <Link href={`/${demo}`} className={buttonVariants({ variant: "ghost", size: "lg" })}>
              See a live page
            </Link>
          </div>
        </div>
        <div className="img-hover overflow-hidden rounded-3xl border border-neutral-200/80 bg-white shadow-[0_20px_60px_-32px_rgba(0,0,0,0.35)]">
          <HeroImage src={hero} alt="SyncLink visual station" />
        </div>
      </section>

      <section className="mx-auto grid w-full max-w-5xl gap-5 px-6 pb-24 md:grid-cols-3">
        {stations.map((station, index) => (
          <Card key={station.n} className="border-neutral-200/80 bg-white/80 shadow-none transition duration-300 hover:-translate-y-1">
            <CardHeader>
              <p className="text-xs tracking-[0.24em] text-neutral-400">{station.n}</p>
              <CardTitle className="text-xl font-medium">{station.title}</CardTitle>
              <CardDescription className="text-neutral-600">{station.body}</CardDescription>
            </CardHeader>
            {index === 1 ? (
              <CardContent>
                <div className="img-hover overflow-hidden rounded-xl bg-[#f4f1ec]">
                  <Image src="/stations/orbit.png" alt="" width={800} height={800} className="mx-auto h-40 w-40 object-contain" />
                </div>
              </CardContent>
            ) : null}
          </Card>
        ))}
      </section>
    </main>
  );
}
