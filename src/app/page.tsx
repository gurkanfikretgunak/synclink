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

const stations = [
  {
    n: "01",
    title: "Compose",
    body: "Name, bio, a short list. Nothing extra.",
  },
  {
    n: "02",
    title: "Publish",
    body: "Your page lives at /your-slug. Fast, public, text-first.",
  },
  {
    n: "03",
    title: "Share",
    body: "One URL. Every link you want people to keep.",
  },
];

export default function Home() {
  return (
    <main className="min-h-full bg-[#faf9f7]">
      <header className="mx-auto flex w-full max-w-5xl items-center justify-between px-6 py-6">
        <p className="text-sm font-medium tracking-[0.22em]">SYNCLINK</p>
        <div className="flex gap-2">
          <Link href="/about" className={buttonVariants({ variant: "ghost" })}>About</Link>
          <Link href="/admin" className={buttonVariants({ variant: "ghost" })}>Admin</Link>
          <Link href="/dashboard" className={buttonVariants({ variant: "outline" })}>Dashboard</Link>
        </div>
      </header>

      <section className="mx-auto grid w-full max-w-5xl items-center gap-12 px-6 pb-16 pt-4 md:grid-cols-2 md:pb-24">
        <div className="space-y-6">
          <p className="text-xs uppercase tracking-[0.28em] text-neutral-500">
            Station 01 — Arrival
          </p>
          <h1 className="text-5xl font-medium tracking-tight text-neutral-950 md:text-6xl">
            One page.
            <br />
            Every link.
          </h1>
          <p className="max-w-sm text-base leading-relaxed text-neutral-600">
            A quieter public page. White space, type, and a few stills. Edit from
            the dashboard.
          </p>
          <div className="flex flex-wrap gap-3">
            <Link href="/dashboard" className={buttonVariants({ size: "lg" })}>
              Create your page
            </Link>
            <Link href="/gurkan" className={buttonVariants({ variant: "ghost", size: "lg" })}>
              See a live page
            </Link>
          </div>
        </div>
        <div className="img-hover overflow-hidden rounded-3xl border border-neutral-200/80 bg-white shadow-[0_20px_60px_-32px_rgba(0,0,0,0.35)]">
          <Image
            src="/stations/hero.png"
            alt="SyncLink visual station"
            width={1600}
            height={1200}
            className="h-auto w-full object-cover"
            priority
          />
        </div>
      </section>

      <section className="mx-auto grid w-full max-w-5xl gap-5 px-6 pb-24 md:grid-cols-3">
        {stations.map((station, index) => (
          <Card key={station.n} className="border-neutral-200/80 bg-white/80 shadow-none">
            <CardHeader>
              <p className="text-xs tracking-[0.24em] text-neutral-400">{station.n}</p>
              <CardTitle className="text-xl font-medium">{station.title}</CardTitle>
              <CardDescription className="text-neutral-600">{station.body}</CardDescription>
            </CardHeader>
            {index === 1 ? (
              <CardContent>
                <div className="img-hover overflow-hidden rounded-xl bg-[#f4f1ec]">
                  <Image
                    src="/stations/orbit.png"
                    alt=""
                    width={800}
                    height={800}
                    className="mx-auto h-40 w-40 object-contain"
                  />
                </div>
              </CardContent>
            ) : null}
          </Card>
        ))}
      </section>
    </main>
  );
}
