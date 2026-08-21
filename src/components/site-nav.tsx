"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { buttonVariants } from "@/components/ui/button";
import { synclink, type NavItem, type PlatformSettings } from "@/lib/api";

export const defaultNav: NavItem[] = [
  { label: "Home", href: "/" },
  { label: "About", href: "/about" },
  { label: "Dashboard", href: "/dashboard" },
  { label: "Admin", href: "/admin" },
];

function itemsFrom(settings?: Partial<PlatformSettings> | null): NavItem[] {
  const raw = settings?.nav;
  if (raw && raw.length) {
    return raw.filter((item) => item.label && item.href);
  }
  return defaultNav;
}

export function SiteNav({
  variant = "full",
  tone = "light",
  settings,
}: {
  variant?: "full" | "slim";
  tone?: "light" | "dark";
  settings?: Partial<PlatformSettings>;
}) {
  const [live, setLive] = useState<Partial<PlatformSettings>>(settings || {});

  useEffect(() => {
    if (settings?.siteName || settings?.nav) {
      setLive(settings);
      return;
    }
    synclink.publicSettings().then(setLive).catch(() => undefined);
  }, [settings]);

  const name = (live.siteName || "SYNCLINK").toUpperCase();
  const items = itemsFrom(live);
  const shown = variant === "slim" ? items.filter((item) => item.href === "/" || item.href === "/dashboard").slice(0, 2) : items;
  const dark = tone === "dark";

  return (
    <header className={`mx-auto flex w-full max-w-6xl items-center justify-between gap-4 px-6 py-5 ${dark ? "text-white" : "text-neutral-950"}`}>
      <Link href="/" className={`text-sm font-medium tracking-[0.22em] ${dark ? "text-white/80" : ""}`}>
        {name}
      </Link>
      <nav className="flex flex-wrap items-center justify-end gap-1">
        {shown.map((item) => (
          <Link
            key={`${item.label}-${item.href}`}
            href={item.href}
            className={buttonVariants({
              variant: item.href === "/dashboard" && variant === "full" ? "outline" : "ghost",
              size: "sm",
            })}
          >
            {item.label}
          </Link>
        ))}
      </nav>
    </header>
  );
}
