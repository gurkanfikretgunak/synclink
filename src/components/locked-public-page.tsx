"use client";

import { FormEvent, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { PublicPageBody } from "@/components/public-page-body";
import { synclink, type PublicPage } from "@/lib/api";

export function LockedPublicPage({ slug }: { slug: string }) {
  const [password, setPassword] = useState("");
  const [page, setPage] = useState<PublicPage | null>(null);
  const [error, setError] = useState("");

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    setError("");
    try {
      setPage(await synclink.getPublicPage(slug, password));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not unlock");
    }
  }

  if (page) return <PublicPageBody page={page} />;

  return (
    <main className="page-enter flex min-h-full items-center justify-center bg-[#faf9f7] px-6">
      <form onSubmit={onSubmit} className="w-full max-w-sm space-y-3 rounded-3xl border border-neutral-200 bg-white p-6">
        <p className="text-xs tracking-[0.24em] text-neutral-400">LOCKED</p>
        <h1 className="text-xl font-medium">This page is private</h1>
        <Input type="password" required value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Page password" />
        {error ? <p className="text-sm text-red-600">{error}</p> : null}
        <Button type="submit">Unlock</Button>
      </form>
    </main>
  );
}
