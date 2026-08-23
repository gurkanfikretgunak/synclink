"use client";

import type { CSSProperties } from "react";
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

  return (
    <a
      href={link.url}
      target="_blank"
      rel="noreferrer"
      className={className}
      style={style}
      onClick={() => {
        void synclink.recordClick(slug, link.id).then((res) => {
          if (typeof res.clicks === "number") setClicks(res.clicks);
          setLast(shortWhen(new Date().toISOString()));
        });
      }}
    >
      <span>{link.title}</span>
      <span className="ml-2 text-[10px] tracking-[0.16em] opacity-50">
        {clicks}
        {last ? ` · ${last}` : ""}
      </span>
    </a>
  );
}
