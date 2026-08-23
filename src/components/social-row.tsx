import { Github, Globe, Instagram, Linkedin, Mail, Youtube } from "lucide-react";
import type { Social } from "@/lib/api";

const icons: Record<string, typeof Github> = {
  github: Github,
  instagram: Instagram,
  linkedin: Linkedin,
  youtube: Youtube,
  email: Mail,
  mail: Mail,
  x: Globe,
  twitter: Globe,
  tiktok: Globe,
  website: Globe,
};

export function SocialRow({ socials, dark = false }: { socials?: Social[] | null; dark?: boolean }) {
  const items = (socials || []).filter((item) => item.network && item.url);
  if (!items.length) return null;
  return (
    <ul className="mt-4 flex flex-wrap items-center justify-center gap-2">
      {items.map((item) => {
        const Icon = icons[item.network.toLowerCase()] || Globe;
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
