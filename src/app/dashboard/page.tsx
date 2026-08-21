"use client";

import Image from "next/image";
import Link from "next/link";
import { FormEvent, useEffect, useState } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button, buttonVariants } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { getToken, setToken, synclink, type LinkItem, type Page } from "@/lib/api";

type Gate = "login" | "forgot" | "reset";

export default function DashboardPage() {
  const [token, setTokenState] = useState<string | null>(null);
  const [ready, setReady] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [gate, setGate] = useState<Gate>("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [resetToken, setResetToken] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [currentPassword, setCurrentPassword] = useState("");
  const [nextPassword, setNextPassword] = useState("");
  const [page, setPage] = useState<Page | null>(null);
  const [links, setLinks] = useState<LinkItem[]>([]);
  const [slug, setSlug] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [bio, setBio] = useState("");
  const [avatarUrl, setAvatarUrl] = useState("");
  const [theme, setTheme] = useState("light");
  const [title, setTitle] = useState("");
  const [url, setUrl] = useState("");
  const [detail, setDetail] = useState("");

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
      setAvatarUrl(nextPage?.avatarUrl ?? "");
      setTheme(nextPage?.theme || "light");
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

  async function onRegister() {
    try {
      const result = await synclink.register(email, password);
      setToken(result.token);
      setTokenState(result.token);
      await load(result.token);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Sign up failed");
    }
  }

  async function onForgot(event: FormEvent) {
    event.preventDefault();
    try {
      const result = await synclink.forgotPassword(email);
      setNotice("If that email exists, a reset token is ready.");
      if (result.resetToken) {
        setResetToken(result.resetToken);
        setGate("reset");
        setNotice("Demo token filled. Set a new password.");
      }
      setError("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Forgot failed");
    }
  }

  async function onReset(event: FormEvent) {
    event.preventDefault();
    try {
      await synclink.resetPassword(email, resetToken, newPassword);
      setNotice("Password updated. Sign in.");
      setGate("login");
      setPassword("");
      setNewPassword("");
      setError("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Reset failed");
    }
  }

  async function onChangePassword(event: FormEvent) {
    event.preventDefault();
    if (!token) return;
    try {
      await synclink.changePassword(token, currentPassword, nextPassword);
      setNotice("Password changed.");
      setCurrentPassword("");
      setNextPassword("");
      setError("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Change failed");
    }
  }

  async function onSavePage(event: FormEvent) {
    event.preventDefault();
    if (!token) return;
    try {
      const next = await synclink.upsertPage(token, {
        slug,
        displayName,
        bio,
        avatarUrl: avatarUrl || null,
        theme,
      });
      setPage(next);
      setNotice("Page saved. Preview is current.");
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
    setNotice("");
  }

  const previewName = displayName || "Your name";
  const previewInitial = previewName.slice(0, 1).toUpperCase();
  const previewLinks = links.filter((link) => link.active);

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
            <Image src="/stations/login.png" alt="" width={1600} height={1200} className="h-auto w-full object-cover" />
          </div>
        </div>
        <div className="flex flex-col justify-center px-6 py-16">
          <Card className="mx-auto w-full max-w-md border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <p className="text-xs tracking-[0.28em] text-neutral-400">STATION 02 — STUDIO</p>
              <CardTitle className="text-2xl font-medium">
                {gate === "login" ? "SyncLink" : gate === "forgot" ? "Forgot password" : "Reset password"}
              </CardTitle>
              <CardDescription>
                {gate === "login"
                  ? "Sign in or create an account."
                  : gate === "forgot"
                    ? "We send a demo token. No mail yet."
                    : "Paste the token and a new password."}
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {gate === "login" ? (
                <form className="space-y-4" onSubmit={onLogin}>
                  <div className="space-y-2"><Label htmlFor="email">Email</Label><Input id="email" type="email" value={email} onChange={(event) => setEmail(event.target.value)} required /></div>
                  <div className="space-y-2"><Label htmlFor="password">Password</Label><Input id="password" type="password" minLength={8} value={password} onChange={(event) => setPassword(event.target.value)} required /></div>
                  {error ? <p className="text-sm text-red-600">{error}</p> : null}
                  {notice ? <p className="text-sm text-neutral-600">{notice}</p> : null}
                  <Button type="submit" className="w-full">Sign in</Button>
                  <Button type="button" variant="outline" className="w-full" onClick={() => void onRegister()}>Sign up</Button>
                  <Button type="button" variant="ghost" className="w-full" onClick={() => { setGate("forgot"); setError(""); setNotice(""); }}>Forgot password</Button>
                </form>
              ) : null}
              {gate === "forgot" ? (
                <form className="space-y-4" onSubmit={onForgot}>
                  <div className="space-y-2"><Label htmlFor="forgot-email">Email</Label><Input id="forgot-email" type="email" value={email} onChange={(event) => setEmail(event.target.value)} required /></div>
                  {error ? <p className="text-sm text-red-600">{error}</p> : null}
                  {notice ? <p className="text-sm text-neutral-600">{notice}</p> : null}
                  <Button type="submit" className="w-full">Get reset token</Button>
                  <Button type="button" variant="ghost" className="w-full" onClick={() => setGate("login")}>Back to sign in</Button>
                </form>
              ) : null}
              {gate === "reset" ? (
                <form className="space-y-4" onSubmit={onReset}>
                  <div className="space-y-2"><Label htmlFor="reset-email">Email</Label><Input id="reset-email" type="email" value={email} onChange={(event) => setEmail(event.target.value)} required /></div>
                  <div className="space-y-2"><Label htmlFor="reset-token">Token</Label><Input id="reset-token" value={resetToken} onChange={(event) => setResetToken(event.target.value)} required /></div>
                  <div className="space-y-2"><Label htmlFor="reset-pass">New password</Label><Input id="reset-pass" type="password" minLength={8} value={newPassword} onChange={(event) => setNewPassword(event.target.value)} required /></div>
                  {error ? <p className="text-sm text-red-600">{error}</p> : null}
                  {notice ? <p className="text-sm text-neutral-600">{notice}</p> : null}
                  <Button type="submit" className="w-full">Reset password</Button>
                  <Button type="button" variant="ghost" className="w-full" onClick={() => setGate("login")}>Back to sign in</Button>
                </form>
              ) : null}
            </CardContent>
          </Card>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-full bg-[#faf9f7]">
      <header className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 py-6">
        <div className="flex items-center gap-3"><Link href="/" className="text-xs tracking-[0.28em] text-neutral-500">SYNCLINK STUDIO</Link><Link href="/admin" className="text-xs tracking-[0.18em] text-neutral-400">ADMIN</Link></div>
        <div className="flex gap-2">
          {slug ? (
            <Link href={`/${slug}`} className={buttonVariants({ variant: "outline", size: "sm" })} target="_blank">
              Open /{slug}
            </Link>
          ) : null}
          <Button variant="ghost" onClick={signOut}>Sign out</Button>
        </div>
      </header>
      <div className="mx-auto grid w-full max-w-6xl gap-6 px-6 pb-16 lg:grid-cols-[0.9fr_1.1fr]">
        <Card className="h-fit border-neutral-200/80 bg-white shadow-none lg:sticky lg:top-6">
          <CardHeader>
            <p className="text-xs tracking-[0.24em] text-neutral-400">LIVE PREVIEW</p>
            <CardTitle className="font-medium">How the page builds</CardTitle>
            <CardDescription>Same layout as /{slug || "your-slug"}.</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="rounded-3xl border border-neutral-200 bg-[#faf9f7] px-6 py-10">
              <div className="flex flex-col items-center gap-5">
                <Avatar className="size-16 border border-neutral-200 bg-white">
                  {avatarUrl ? <AvatarImage src={avatarUrl} alt="" /> : null}
                  <AvatarFallback className="bg-white">{previewInitial}</AvatarFallback>
                </Avatar>
                <div className="space-y-1 text-center">
                  <p className="text-xl font-medium">{previewName}</p>
                  <p className="text-sm text-neutral-600">{bio || "Bio appears here."}</p>
                </div>
                <ul className="w-full space-y-2">
                  {previewLinks.length ? previewLinks.map((link) => (
                    <li key={link.id} className="rounded-2xl border border-neutral-200 bg-white px-4 py-3 text-center text-sm">
                      {link.title}
                    </li>
                  )) : (
                    <li className="rounded-2xl border border-dashed border-neutral-300 px-4 py-3 text-center text-sm text-neutral-400">
                      Links land here as you add them
                    </li>
                  )}
                </ul>
                <p className="text-xs text-neutral-400">theme · {theme || "light"}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <div className="space-y-6">
          {error ? <p className="text-sm text-red-600">{error}</p> : null}
          {notice ? <p className="text-sm text-neutral-600">{notice}</p> : null}

          <Card className="border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <p className="text-xs tracking-[0.24em] text-neutral-400">PAGE</p>
              <CardTitle className="font-medium">Identity</CardTitle>
              <CardDescription>Slug, name, bio. This is the public header.</CardDescription>
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
              <p className="text-xs tracking-[0.24em] text-neutral-400">DETAILS</p>
              <CardTitle className="font-medium">How it should feel</CardTitle>
              <CardDescription>Avatar, theme, and notes you want on this page.</CardDescription>
            </CardHeader>
            <CardContent>
              <form className="space-y-4" onSubmit={onSavePage}>
                <div className="space-y-2"><Label htmlFor="avatarUrl">Avatar URL</Label><Input id="avatarUrl" type="url" value={avatarUrl} onChange={(event) => setAvatarUrl(event.target.value)} placeholder="https://" /></div>
                <div className="space-y-2"><Label htmlFor="theme">Theme</Label><Input id="theme" value={theme} onChange={(event) => setTheme(event.target.value)} placeholder="light" /></div>
                <div className="space-y-2"><Label htmlFor="detail">Notes for this page</Label><Textarea id="detail" value={detail} onChange={(event) => setDetail(event.target.value)} placeholder="What else should live here later." /></div>
                <Button type="submit">Save details</Button>
              </form>
            </CardContent>
          </Card>

          <Card className="border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <p className="text-xs tracking-[0.24em] text-neutral-400">LINKS</p>
              <CardTitle className="font-medium">What people tap</CardTitle>
              <CardDescription>Title and URL. Hidden links stay off the preview.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <form className="space-y-4" onSubmit={onAddLink}>
                <div className="space-y-2"><Label htmlFor="title">Title</Label><Input id="title" value={title} onChange={(event) => setTitle(event.target.value)} required /></div>
                <div className="space-y-2"><Label htmlFor="url">URL</Label><Input id="url" type="url" value={url} onChange={(event) => setUrl(event.target.value)} required /></div>
                <Button type="submit">Add link</Button>
              </form>
              <ul className="space-y-3">
                {links.length ? links.map((link) => (
                  <li key={link.id} className="flex items-center justify-between gap-3 rounded-xl border border-neutral-200 bg-[#faf9f7] px-3 py-2">
                    <span className={link.active ? "" : "text-neutral-400"}>{link.title}</span>
                    <span className="flex gap-2">
                      <Button size="sm" variant="outline" onClick={() => void onToggle(link)}>{link.active ? "Hide" : "Show"}</Button>
                      <Button size="sm" variant="ghost" onClick={() => void onDelete(link.id)}>Delete</Button>
                    </span>
                  </li>
                )) : (
                  <li className="rounded-xl border border-dashed border-neutral-300 px-3 py-4 text-sm text-neutral-400">No links yet. Add the first one above.</li>
                )}
              </ul>
            </CardContent>
          </Card>

          <Card className="border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <p className="text-xs tracking-[0.24em] text-neutral-400">ACCOUNT</p>
              <CardTitle className="font-medium">Password</CardTitle>
              <CardDescription>Change it here. Forgot/reset lives on the sign-in gate.</CardDescription>
            </CardHeader>
            <CardContent>
              <form className="space-y-4" onSubmit={onChangePassword}>
                <div className="space-y-2"><Label htmlFor="currentPassword">Current password</Label><Input id="currentPassword" type="password" value={currentPassword} onChange={(event) => setCurrentPassword(event.target.value)} required /></div>
                <div className="space-y-2"><Label htmlFor="nextPassword">New password</Label><Input id="nextPassword" type="password" minLength={8} value={nextPassword} onChange={(event) => setNextPassword(event.target.value)} required /></div>
                <Button type="submit">Change password</Button>
              </form>
            </CardContent>
          </Card>
        </div>
      </div>
    </main>
  );
}
