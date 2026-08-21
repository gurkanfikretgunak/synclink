"use client";

import Image from "next/image";
import { FormEvent, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { getToken, setToken, synclink, type LinkItem, type Page } from "@/lib/api";

export default function DashboardPage() {
  const [token, setTokenState] = useState<string | null>(null);
  const [ready, setReady] = useState(false);
  const [error, setError] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [page, setPage] = useState<Page | null>(null);
  const [links, setLinks] = useState<LinkItem[]>([]);
  const [slug, setSlug] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [bio, setBio] = useState("");
  const [title, setTitle] = useState("");
  const [url, setUrl] = useState("");

  async function load(nextToken: string) {
    try {
      const [nextPage, nextLinks] = await Promise.all([
        synclink.getMyPage(nextToken).catch(() => null),
        synclink.listLinks(nextToken).catch(() => []),
      ]);
      setPage(nextPage);
      setLinks(nextLinks);
      setSlug(nextPage?.slug ?? "");
      setDisplayName(nextPage?.displayName ?? "");
      setBio(nextPage?.bio ?? "");
      setError("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load");
    }
  }

  useEffect(() => {
    const existing = getToken();
    setTokenState(existing);
    setReady(true);
    if (existing) void load(existing);
  }, []);

  async function onLogin(event: FormEvent) {
    event.preventDefault();
    try {
      const result = await synclink.login(email, password);
      setToken(result.token);
      setTokenState(result.token);
      await load(result.token);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    }
  }

  async function onRegister(event?: { preventDefault: () => void }) {
    event?.preventDefault();
    try {
      const result = await synclink.register(email, password);
      setToken(result.token);
      setTokenState(result.token);
      await load(result.token);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Sign up failed");
    }
  }

  async function onSavePage(event: FormEvent) {
    event.preventDefault();
    if (!token) return;
    try {
      const next = await synclink.upsertPage(token, { slug, displayName, bio, theme: page?.theme || "light" });
      setPage(next);
      setError("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Save failed");
    }
  }

  async function onAddLink(event: FormEvent) {
    event.preventDefault();
    if (!token) return;
    try {
      const link = await synclink.createLink(token, { title, url, active: true });
      setLinks((current) => [...current, link]);
      setTitle("");
      setUrl("");
      setError("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not add link");
    }
  }

  async function onToggle(link: LinkItem) {
    if (!token) return;
    const next = await synclink.updateLink(token, link.id, { active: !link.active });
    setLinks((current) => current.map((item) => (item.id === next.id ? next : item)));
  }

  async function onDelete(id: string) {
    if (!token) return;
    await synclink.deleteLink(token, id);
    setLinks((current) => current.filter((item) => item.id !== id));
  }

  function signOut() {
    setToken(null);
    setTokenState(null);
    setPage(null);
    setLinks([]);
  }

  if (!ready) {
    return (
      <main className="flex min-h-full items-center justify-center bg-[#faf9f7]">
        <p className="text-xs tracking-[0.28em] text-neutral-400">SYNCLINK</p>
      </main>
    );
  }

  if (!token) {
    return (
      <main className="grid min-h-full bg-[#faf9f7] md:grid-cols-2">
        <div className="hidden items-center justify-center p-12 md:flex">
          <div className="w-full max-w-md overflow-hidden rounded-3xl border border-neutral-200 bg-white">
            <Image src="/stations/hero.png" alt="" width={1600} height={1200} className="h-auto w-full object-cover" />
          </div>
        </div>
        <div className="flex flex-col justify-center px-6 py-16">
          <Card className="mx-auto w-full max-w-md border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <p className="text-xs tracking-[0.28em] text-neutral-400">STATION 02 — STUDIO</p>
              <CardTitle className="text-2xl font-medium">SyncLink</CardTitle>
              <CardDescription>Sign in or create an account.</CardDescription>
            </CardHeader>
            <CardContent>
              <form className="space-y-4" onSubmit={onLogin}>
                <div className="space-y-2">
                  <Label htmlFor="email">Email</Label>
                  <Input id="email" type="email" value={email} onChange={(event) => setEmail(event.target.value)} required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password">Password</Label>
                  <Input id="password" type="password" minLength={8} value={password} onChange={(event) => setPassword(event.target.value)} required />
                </div>
                {error ? <p className="text-sm text-red-600">{error}</p> : null}
                <Button type="submit" className="w-full">Sign in</Button>
                <Button type="button" variant="outline" className="w-full" onClick={onRegister}>Sign up</Button>
              </form>
            </CardContent>
          </Card>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-full bg-[#faf9f7]">
      <header className="mx-auto flex w-full max-w-5xl items-center justify-between px-6 py-6">
        <p className="text-xs tracking-[0.28em] text-neutral-500">SYNCLINK STUDIO</p>
        <Button variant="ghost" onClick={signOut}>Sign out</Button>
      </header>
      <div className="mx-auto grid w-full max-w-5xl gap-6 px-6 pb-16 md:grid-cols-[1.1fr_0.9fr]">
        {error ? <p className="text-sm text-red-600 md:col-span-2">{error}</p> : null}
        <Card className="border-neutral-200/80 bg-white shadow-none">
          <CardHeader>
            <p className="text-xs tracking-[0.24em] text-neutral-400">PAGE</p>
            <CardTitle className="font-medium">Public at /{slug || "your-slug"}</CardTitle>
            <CardDescription>Name and bio for the live page.</CardDescription>
          </CardHeader>
          <CardContent>
            <form className="space-y-4" onSubmit={onSavePage}>
              <div className="space-y-2"><Label htmlFor="slug">Slug</Label><Input id="slug" value={slug} onChange={(event) => setSlug(event.target.value)} required /></div>
              <div className="space-y-2"><Label htmlFor="displayName">Display name</Label><Input id="displayName" value={displayName} onChange={(event) => setDisplayName(event.target.value)} required /></div>
              <div className="space-y-2"><Label htmlFor="bio">Bio</Label><Textarea id="bio" value={bio} onChange={(event) => setBio(event.target.value)} /></div>
              <Button type="submit">Save page</Button>
            </form>
          </CardContent>
        </Card>
        <Card className="border-neutral-200/80 bg-white shadow-none">
          <CardHeader>
            <p className="text-xs tracking-[0.24em] text-neutral-400">LINKS</p>
            <CardTitle className="font-medium">Title and URL</CardTitle>
            <CardDescription>What people tap on your page.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <form className="space-y-4" onSubmit={onAddLink}>
              <div className="space-y-2"><Label htmlFor="title">Title</Label><Input id="title" value={title} onChange={(event) => setTitle(event.target.value)} required /></div>
              <div className="space-y-2"><Label htmlFor="url">URL</Label><Input id="url" type="url" value={url} onChange={(event) => setUrl(event.target.value)} required /></div>
              <Button type="submit">Add link</Button>
            </form>
            <ul className="space-y-3">
              {links.map((link) => (
                <li key={link.id} className="flex items-center justify-between gap-3 rounded-xl border border-neutral-200 bg-[#faf9f7] px-3 py-2">
                  <span className={link.active ? "" : "text-neutral-400"}>{link.title}</span>
                  <span className="flex gap-2">
                    <Button size="sm" variant="outline" onClick={() => void onToggle(link)}>{link.active ? "Hide" : "Show"}</Button>
                    <Button size="sm" variant="ghost" onClick={() => void onDelete(link.id)}>Delete</Button>
                  </span>
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>
      </div>
    </main>
  );
}
