"use client";

import { useState } from "react";

export function PublicHandle({ slug, dark = false }: { slug: string; dark?: boolean }) {
  const [copied, setCopied] = useState(false);

  if (!slug) return null;

  async function copy() {
    const url = typeof window !== "undefined" ? `${window.location.origin}/${slug}` : `/${slug}`;
    try {
      await navigator.clipboard.writeText(url);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 1500);
    } catch {
      // ignore clipboard errors
    }
  }

  return (
    <button
      type="button"
      aria-label="Copy public URL"
      onClick={() => void copy()}
      className={`mt-1 cursor-pointer border-0 bg-transparent p-0 font-mono text-[10px] tracking-[0.16em] ${dark ? "text-white/40" : "text-neutral-400"}`}
    >
      {copied ? "COPIED" : `/${slug}`}
    </button>
  );
}
