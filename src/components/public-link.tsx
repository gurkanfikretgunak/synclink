"use client";

import type { CSSProperties } from "react";
import { useState } from "react";
import { synclink, type PublicLink } from "@/lib/api";

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
        });
      }}
    >
      <span>{link.title}</span>
      <span className="ml-2 text-[10px] tracking-[0.16em] opacity-50">{clicks}</span>
    </a>
  );
}
