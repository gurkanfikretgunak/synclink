import Link from "next/link";
import { buttonVariants } from "@/components/ui/button";
import { synclink } from "@/lib/api";

export default async function AboutPage() {
  let settings = {
    siteName: "SyncLink",
    tagline: "One page. Every link.",
    about: "SyncLink is a quiet public page for the links you want people to keep.",
    supportEmail: "",
  };
  try {
    settings = { ...settings, ...(await synclink.publicSettings()) };
  } catch {
    /* local/demo fallback */
  }

  return (
    <main className="min-h-full bg-[#faf9f7]">
      <header className="mx-auto flex w-full max-w-3xl items-center justify-between px-6 py-6">
        <p className="text-xs tracking-[0.28em] text-neutral-500">{(settings.siteName || "SYNCLINK").toUpperCase()}</p>
        <Link href="/" className={buttonVariants({ variant: "ghost", size: "sm" })}>Home</Link>
      </header>
      <article className="mx-auto w-full max-w-3xl space-y-6 px-6 pb-20">
        <p className="text-xs tracking-[0.28em] text-neutral-400">ABOUT</p>
        <h1 className="text-4xl font-medium tracking-tight">{settings.tagline || "One page. Every link."}</h1>
        <p className="whitespace-pre-wrap text-neutral-600">{settings.about}</p>
        {settings.supportEmail ? <p className="text-sm text-neutral-500">{settings.supportEmail}</p> : null}
      </article>
    </main>
  );
}
