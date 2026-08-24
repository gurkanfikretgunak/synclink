"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";

type ContactBits = {
  displayName?: string;
  bio?: string;
  avatarUrl?: string | null;
  socials?: { network: string; url: string }[];
};

function escapeVcard(value: string): string {
  return value
    .replace(/\\/g, "\\\\")
    .replace(/\r\n/g, "\n")
    .replace(/\n/g, "\\n")
    .replace(/\r/g, "")
    .replace(/,/g, "\\,")
    .replace(/;/g, "\\;");
}

function vcardName(displayName: string): { fn: string; n: string } {
  const fn = displayName.trim();
  const parts = fn.split(/\s+/).filter(Boolean);
  if (parts.length <= 1) {
    return { fn: fn || displayName, n: `${escapeVcard(parts[0] || fn)};;;;` };
  }
  const family = parts[parts.length - 1];
  const given = parts.slice(0, -1).join(" ");
  return { fn, n: `${escapeVcard(family)};${escapeVcard(given)};;;` };
}

function buildVcard(slug: string, pageUrl: string, contact: ContactBits): string {
  const { fn, n } = vcardName(contact.displayName?.trim() || slug);
  const lines = ["BEGIN:VCARD", "VERSION:3.0", `FN:${escapeVcard(fn)}`, `N:${n}`];
  const note = (contact.bio || "").trim();
  if (note) lines.push(`NOTE:${escapeVcard(note)}`);
  if (pageUrl) lines.push(`URL:${pageUrl}`);
  const seen = new Set(pageUrl ? [pageUrl] : []);
  for (const item of contact.socials || []) {
    const url = (item.url || "").trim();
    if (!url || seen.has(url)) continue;
    seen.add(url);
    lines.push(`URL:${url}`);
  }
  const photo = (contact.avatarUrl || "").trim();
  if (/^https?:\/\//i.test(photo)) {
    lines.push(`PHOTO;VALUE=URI:${photo}`);
  }
  lines.push("END:VCARD");
  return `${lines.join("\r\n")}\r\n`;
}

export function ShareQr({
  slug,
  dark = false,
  displayName,
  bio,
  avatarUrl,
  socials,
}: { slug: string; dark?: boolean } & ContactBits) {
  const [href, setHref] = useState("");
  const [copied, setCopied] = useState(false);
  const [canShare, setCanShare] = useState(false);

  useEffect(() => {
    if (!slug) {
      setHref("");
      return;
    }
    setHref(`${window.location.origin}/${slug}`);
    setCanShare(typeof navigator !== "undefined" && typeof navigator.share === "function");
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

  async function share() {
    if (!href) return;
    try {
      await navigator.share({ title: document.title, url: href, text: href });
    } catch (err) {
      if (err instanceof Error && err.name === "AbortError") return;
      await copy();
    }
  }

  function addToContacts() {
    const pageUrl = href || `https://synclink-mocha.vercel.app/${slug}`;
    const vcf = buildVcard(slug, pageUrl, { displayName, bio, avatarUrl, socials });
    const blob = new Blob([vcf], { type: "text/vcard;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${slug}.vcf`;
    document.body.appendChild(a);
    a.click();
    a.remove();
    window.setTimeout(() => URL.revokeObjectURL(url), 1000);
  }

  return (
    <div className={`mt-8 w-full space-y-3 rounded-2xl border px-4 py-4 ${dark ? "border-white/15 bg-white/5" : "border-neutral-200 bg-white"}`}>
      <p className={`text-xs tracking-[0.16em] ${dark ? "text-white/50" : "text-neutral-400"}`}>SHARE</p>
      {src ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={src} alt={`QR for ${href}`} width={240} height={240} className="mx-auto h-40 w-40 bg-white p-2" />
      ) : null}
      <p className={`break-all text-center text-xs ${dark ? "text-white/60" : "text-neutral-500"}`}>{href || `/${slug}`}</p>
      <div className="flex flex-wrap justify-center gap-2">
        {canShare ? (
          <Button type="button" size="sm" onClick={() => void share()}>
            Share
          </Button>
        ) : null}
        <Button type="button" size="sm" variant="outline" onClick={() => void copy()}>
          {copied ? "Copied" : "Copy URL"}
        </Button>
        <button type="button" onClick={addToContacts} className="inline-flex h-7 items-center rounded-lg px-2.5 text-[0.8rem] underline">
          Add to contacts
        </button>
        {src ? (
          <a href={src} download={`synclink-${slug}.png`} target="_blank" rel="noreferrer" className="inline-flex h-7 items-center rounded-lg px-2.5 text-[0.8rem] underline">
            QR PNG
          </a>
        ) : null}
      </div>
    </div>
  );
}
