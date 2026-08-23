export function PageCover({ url, kind = "image", dark = false }: { url?: string | null; kind?: string; dark?: boolean }) {
  if (!url) return null;
  const frame = dark ? "border-white/15 bg-black/20" : "border-neutral-200 bg-white";
  if (kind === "video") {
    const yt = url.match(/(?:youtu\.be\/|youtube\.com\/(?:watch\?v=|embed\/))([A-Za-z0-9_-]{6,})/);
    const src = yt ? `https://www.youtube.com/embed/${yt[1]}` : url;
    return (
      <div className={`mb-8 aspect-video w-full overflow-hidden rounded-3xl border ${frame}`}>
        <iframe title="Cover" src={src} className="h-full w-full" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowFullScreen />
      </div>
    );
  }
  return (
    // eslint-disable-next-line @next/next/no-img-element
    <img src={url} alt="" className={`mb-8 h-40 w-full rounded-3xl border object-cover ${frame}`} />
  );
}
