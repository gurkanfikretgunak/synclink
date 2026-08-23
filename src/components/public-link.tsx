"use client";

import type { CSSProperties, MouseEvent } from "react";
import { useState } from "react";
import { synclink, type PublicLink } from "@/lib/api";

function shortWhen(value?: string | null) {
  if (!value) return "";
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return "";
  return d.toLocaleString("tr-TR", { day: "2-digit", month: "short", hour: "2-digit", minute: "2-digit" });
}

export function PublicLinkButton({
  slug,
  link,
  className,
  style,
}: {
  slug: string;
  link: PublicLink;
  className?: string;
  style?: CSSProperties;
}) {
  const [clicks, setClicks] = useState(link.clicks ?? 0);
  const [last, setLast] = useState(shortWhen(link.lastClickedAt));

  function onClick(event: MouseEvent<HTMLAnchorElement>) {
    if (link.sensitive && !window.confirm("This link is marked sensitive / 18+.")) {
      event.preventDefault();
      return;
    }
    void synclink.recordClick(slug, link.id).then((res) => {
      if (typeof res.clicks === "number") setClicks(res.clicks);
      setLast(shortWhen(new Date().toISOString()));
    });
  }

  return (
    <a href={link.url} target="_blank" rel="noreferrer" className={className} style={style} onClick={onClick}>
      {link.thumbnailUrl ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={link.thumbnailUrl} alt="" className="mx-auto mb-2 h-10 w-10 rounded-lg object-cover" />
      ) : null}
      <span>{link.featured ? "★ " : ""}{link.title}</span>
      <span className="ml-2 text-[10px] tracking-[0.16em] opacity-50">
        {clicks}
        {last ? ` · ${last}` : ""}
      </span>
      {link.embedUrl ? (
        <span className="mt-3 block overflow-hidden rounded-xl">
          <iframe title={link.title} src={link.embedUrl} className="aspect-video h-auto w-full" />
        </span>
      ) : null}
    </a>
  );
}
