import Image from "next/image";
import { notFound } from "next/navigation";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { synclink } from "@/lib/api";

const bgClass: Record<string, string> = {
  cream: "bg-[#faf9f7] text-neutral-950",
  white: "bg-white text-neutral-950",
  dark: "bg-[#111111] text-white",
  motion: "bg-motion text-neutral-950",
};

const shapeClass: Record<string, string> = {
  circle: "rounded-full",
  rounded: "rounded-2xl",
  square: "rounded-none",
};

const motionClass: Record<string, string> = {
  none: "transition-colors",
  subtle: "transition-transform duration-300 hover:-translate-y-0.5",
  lively: "link-lively transition-transform duration-300 hover:-translate-y-1 hover:scale-[1.02]",
};

export default async function PublicPage({ params }: PageProps<"/[slug]">) {
  const { slug } = await params;
  let page;
  try {
    page = await synclink.getPublicPage(slug);
  } catch {
    notFound();
  }

  const tone = bgClass[page.background] || bgClass.cream;
  const shape = shapeClass[page.avatarShape] || shapeClass.circle;
  const motion = motionClass[page.motion] || motionClass.subtle;
  const accent = page.accentColor || "#111111";
  const dark = page.background === "dark";

  return (
    <main className={`relative min-h-full overflow-hidden ${tone}`}>
      <div className="relative mx-auto flex min-h-full w-full max-w-md flex-col items-center px-6 py-16">
        <div className="img-hover mb-8 overflow-hidden rounded-3xl">
          <Image src="/stations/orbit.png" alt="" width={720} height={480} className="h-28 w-full object-cover opacity-90" />
        </div>
        <Avatar className={`size-20 border ${shape} ${dark ? "border-white/15" : "border-neutral-200"}`}>
          {page.avatarUrl ? <AvatarImage src={page.avatarUrl} alt={page.displayName} /> : null}
          <AvatarFallback className={dark ? "bg-white/10 text-white" : ""}>
            {page.displayName.slice(0, 1).toUpperCase()}
          </AvatarFallback>
        </Avatar>
        <h1 className="mt-5 text-2xl font-medium tracking-tight">{page.displayName}</h1>
        {page.bio ? <p className={`mt-2 text-center text-sm leading-6 ${dark ? "text-white/70" : "text-neutral-600"}`}>{page.bio}</p> : null}
        <ul className="mt-10 w-full space-y-3">
          {page.links.map((link: { id: string; url: string; title: string }) => (
            <li key={link.id}>
              <a
                href={link.url}
                target="_blank"
                rel="noreferrer"
                className={`block rounded-2xl border px-4 py-4 text-center text-sm ${motion} ${dark ? "border-white/15 bg-white/5" : "border-neutral-200 bg-white"}`}
                style={{ boxShadow: `0 0 0 1px ${accent}14` }}
              >
                {link.title}
              </a>
            </li>
          ))}
        </ul>
      </div>
    </main>
  );
}
