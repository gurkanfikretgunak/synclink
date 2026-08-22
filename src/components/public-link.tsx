"use client";

import type { CSSProperties } from "react";
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
  return (
    <a
      href={link.url}
      target="_blank"
      rel="noreferrer"
      className={className}
      style={style}
      onClick={() => {
        void synclink.recordClick(slug, link.id);
      }}
    >
      <span>{link.title}</span>
      {typeof link.clicks === "number" ? (
        <span className="ml-2 text-[10px] tracking-[0.16em] opacity-50">{link.clicks}</span>
      ) : null}
    </a>
  );
}
