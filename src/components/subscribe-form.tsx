"use client";

import { FormEvent, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { synclink } from "@/lib/api";

export function SubscribeForm({ slug, dark = false }: { slug: string; dark?: boolean }) {
  const [email, setEmail] = useState("");
  const [hidden, setHidden] = useState(false);
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  if (hidden) return null;

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    setBusy(true);
    setError("");
    setNotice("");
    try {
      await synclink.subscribe(slug, email.trim());
      setNotice("You're on the list.");
      setEmail("");
    } catch (err) {
      const message = err instanceof Error ? err.message : "Could not subscribe";
      if (/not found/i.test(message)) {
        setHidden(true);
        return;
      }
      setError(message);
    } finally {
      setBusy(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className={`mt-10 w-full space-y-2 rounded-2xl border px-4 py-4 ${dark ? "border-white/15 bg-white/5" : "border-neutral-200 bg-white"}`}>
      <p className={`text-xs tracking-[0.16em] ${dark ? "text-white/50" : "text-neutral-400"}`}>EMAIL</p>
      <div className="flex gap-2">
        <Input
          type="email"
          required
          placeholder="you@mail"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className={dark ? "border-white/15 bg-transparent text-white" : ""}
        />
        <Button type="submit" disabled={busy}>{busy ? "…" : "Join"}</Button>
      </div>
      {notice ? <p className="text-xs opacity-70">{notice}</p> : null}
      {error ? <p className="text-xs text-red-500">{error}</p> : null}
    </form>
  );
}
