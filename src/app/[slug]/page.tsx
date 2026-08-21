import Image from "next/image";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { buttonVariants } from "@/components/ui/button";
import { synclink } from "@/lib/api";

type Props = { params: Promise<{ slug: string }> };

export default async function PublicPage({ params }: Props) {
  const { slug } = await params;
  let page = null;
  let error = "";
  try {
    page = await synclink.getPublicPage(slug);
  } catch (err) {
    error = err instanceof Error ? err.message : "Page not found";
  }

  if (!page) {
    return (
      <main className="flex min-h-full flex-col items-center justify-center bg-[#faf9f7] px-6 py-16">
        <p className="mb-6 text-xs tracking-[0.28em] text-neutral-400">SYNCLINK</p>
        <h1 className="text-3xl font-medium tracking-tight">{slug}</h1>
        <p className="mt-2 text-neutral-600">{error}</p>
      </main>
    );
  }

  const initial = (page.displayName || page.slug).slice(0, 1).toUpperCase();

  return (
    <main className="min-h-full bg-[#faf9f7]">
      <div className="mx-auto flex w-full max-w-md flex-col items-center gap-8 px-6 py-16">
        <p className="text-xs tracking-[0.28em] text-neutral-400">SYNCLINK</p>
        <Avatar className="size-20 border border-neutral-200 bg-white">
          {page.avatarUrl ? <AvatarImage src={page.avatarUrl} alt="" /> : null}
          <AvatarFallback className="bg-white text-lg">{initial}</AvatarFallback>
        </Avatar>
        <header className="space-y-2 text-center">
          <h1 className="text-3xl font-medium tracking-tight">{page.displayName}</h1>
          {page.bio ? <p className="text-neutral-600">{page.bio}</p> : null}
        </header>
        <ul className="w-full space-y-3">
          {page.links.map((link) => (
            <li key={link.id}>
              <a
                href={link.url}
                className={buttonVariants({
                  variant: "outline",
                  className:
                    "h-12 w-full justify-center rounded-2xl border-neutral-200 bg-white text-base font-medium hover:bg-white",
                })}
                rel="noreferrer"
                target="_blank"
              >
                {link.title}
              </a>
            </li>
          ))}
        </ul>
        <Image
          src="/stations/orbit.png"
          alt=""
          width={80}
          height={80}
          className="mt-6 opacity-70"
        />
      </div>
    </main>
  );
}
