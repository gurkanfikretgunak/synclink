"use client";

import { FormEvent, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
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
    return <main className="mx-auto w-full max-w-lg px-6 py-16"><p className="text-sm text-neutral-500">SyncLink</p></main>;
  }

  if (!token) {
    return (
      <main className="mx-auto flex min-h-full w-full max-w-md flex-col justify-center px-6 py-16">
        <Card>
          <CardHeader>
            <CardTitle>SyncLink</CardTitle>
            <CardDescription>Sign in to edit your page.</CardDescription>
          </CardHeader>
          <CardContent>
            <form className="space-y-4" onSubmit={onLogin}>
              <div className="space-y-2"><Label htmlFor="email">Email</Label><Input id="email" type="email" value={email} onChange={(event) => setEmail(event.target.value)} required /></div>
              <div className="space-y-2"><Label htmlFor="password">Password</Label><Input id="password" type="password" value={password} onChange={(event) => setPassword(event.target.value)} required /></div>
              {error ? <p className="text-sm text-red-600">{error}</p> : null}
              <Button type="submit" className="w-full">Sign in</Button>
            </form>
          </CardContent>
        </Card>
      </main>
    );
  }

  return (
    <main className="mx-auto flex min-h-full w-full max-w-lg flex-col gap-8 px-6 py-16">
      <div className="flex items-center justify-between">
        <p className="text-sm text-neutral-500">SyncLink</p>
        <Button variant="ghost" onClick={signOut}>Sign out</Button>
      </div>
      {error ? <p className="text-sm text-red-600">{error}</p> : null}
      <Card>
        <CardHeader>
          <CardTitle>Page</CardTitle>
          <CardDescription>Public at /{slug || "your-slug"}</CardDescription>
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
      <Separator />
      <Card>
        <CardHeader>
          <CardTitle>Links</CardTitle>
          <CardDescription>Title and URL only.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <form className="space-y-4" onSubmit={onAddLink}>
            <div className="space-y-2"><Label htmlFor="title">Title</Label><Input id="title" value={title} onChange={(event) => setTitle(event.target.value)} required /></div>
            <div className="space-y-2"><Label htmlFor="url">URL</Label><Input id="url" type="url" value={url} onChange={(event) => setUrl(event.target.value)} required /></div>
            <Button type="submit">Add link</Button>
          </form>
          <ul className="space-y-3">
            {links.map((link) => (
              <li key={link.id} className="flex items-center justify-between gap-3 border border-neutral-200 px-3 py-2">
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
    </main>
  );
}
