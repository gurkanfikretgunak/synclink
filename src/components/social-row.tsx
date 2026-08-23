import { Globe, Link as LinkIcon, Mail } from "lucide-react";
import type { Social } from "@/lib/api";

function iconFor(network: string) {
  const key = network.toLowerCase();
  if (key === "email" || key === "mail") return Mail;
  if (key === "website" || key === "url") return LinkIcon;
  return Globe;
}

export function SocialRow({ socials, dark = false }: { socials?: Social[] | null; dark?: boolean }) {
  const items = (socials || []).filter((item) => item.network && item.url);
  if (!items.length) return null;
  return (
    <ul className="mt-4 flex flex-wrap items-center justify-center gap-2">
      {items.map((item) => {
        const Icon = iconFor(item.network);
        return (
          <li key={`${item.network}-${item.url}`}>
            <a
              href={item.url}
              target="_blank"
              rel="noreferrer"
              aria-label={item.network}
              className={`inline-flex size-9 items-center justify-center rounded-full border text-sm ${dark ? "border-white/20" : "border-neutral-200 bg-white"}`}
            >
              <Icon className="size-4" />
            </a>
          </li>
        );
      })}
    </ul>
  );
}
