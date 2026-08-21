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
      <main className="mx-auto flex min-h-full w-full max-w-md flex-col justify-center gap-3 px-6 py-16">
        <p className="text-sm text-neutral-500">SyncLink</p>
        <h1 className="text-2xl font-medium">{slug}</h1>
        <p className="text-neutral-600">{error}</p>
      </main>
    );
  }

  return (
    <main className="mx-auto flex min-h-full w-full max-w-md flex-col gap-8 px-6 py-16">
      <header className="space-y-2">
        <p className="text-sm text-neutral-500">SyncLink</p>
        <h1 className="text-3xl font-medium tracking-tight">{page.displayName}</h1>
        {page.bio ? <p className="text-neutral-600">{page.bio}</p> : null}
      </header>
      <ul className="space-y-3">
        {page.links.map((link) => (
          <li key={link.id}>
            <a href={link.url} className="block border border-neutral-200 px-4 py-3 text-center hover:bg-neutral-50" rel="noreferrer" target="_blank">
              {link.title}
            </a>
          </li>
        ))}
      </ul>
    </main>
  );
}
