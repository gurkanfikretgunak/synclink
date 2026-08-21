import { SiteNav } from "@/components/site-nav";
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
      <SiteNav settings={settings} />
      <article className="mx-auto w-full max-w-3xl space-y-6 px-6 pb-20">
        <p className="text-xs tracking-[0.28em] text-neutral-400">ABOUT</p>
        <h1 className="text-4xl font-medium tracking-tight">{settings.tagline || "One page. Every link."}</h1>
        <p className="whitespace-pre-wrap text-neutral-600">{settings.about}</p>
        {settings.supportEmail ? <p className="text-sm text-neutral-500">{settings.supportEmail}</p> : null}
      </article>
    </main>
  );
}
