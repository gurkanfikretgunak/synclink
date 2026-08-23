"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";

export function ShareQr({ slug, dark = false }: { slug: string; dark?: boolean }) {
  const [href, setHref] = useState("");
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    if (!slug) {
      setHref("");
      return;
    }
    setHref(`${window.location.origin}/${slug}`);
  }, [slug]);

  if (!slug) return null;

  const src = href
    ? `https://api.qrserver.com/v1/create-qr-code/?size=240x240&margin=8&data=${encodeURIComponent(href)}`
    : "";

  async function copy() {
    if (!href) return;
    try {
      await navigator.clipboard.writeText(href);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 1600);
    } catch {
      setCopied(false);
    }
  }

  return (
    <div className={`mt-8 w-full space-y-3 rounded-2xl border px-4 py-4 ${dark ? "border-white/15 bg-white/5" : "border-neutral-200 bg-white"}`}>
      <p className={`text-xs tracking-[0.16em] ${dark ? "text-white/50" : "text-neutral-400"}`}>SHARE</p>
      {src ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={src} alt={`QR for ${href}`} width={240} height={240} className="mx-auto h-40 w-40 bg-white p-2" />
      ) : null}
      <p className={`break-all text-center text-xs ${dark ? "text-white/60" : "text-neutral-500"}`}>{href || `/${slug}`}</p>
      <div className="flex justify-center gap-2">
        <Button type="button" size="sm" variant="outline" onClick={() => void copy()}>
          {copied ? "Copied" : "Copy URL"}
        </Button>
        {src ? (
          <a href={src} download={`synclink-${slug}.png`} target="_blank" rel="noreferrer" className="inline-flex h-7 items-center rounded-lg px-2.5 text-[0.8rem] underline">
            QR PNG
          </a>
        ) : null}
      </div>
    </div>
  );
}
