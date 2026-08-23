"use client";

import Image from "next/image";
import Link from "next/link";
import { FormEvent, useEffect, useState } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";
import { SiteNav } from "@/components/site-nav";
import { getToken, setToken, synclink, type LinkItem, type Page } from "@/lib/api";

type Gate = "login" | "register" | "forgot" | "reset";
type Panel = "identity" | "look" | "links" | "account";

const emptyPage: Page = {
  id: "",
  slug: "",
  displayName: "",
  bio: "",
  avatarUrl: null,
  theme: "minimal",
  avatarShape: "circle",
  accentColor: "#111111",
  background: "cream",
  motion: "subtle",
  createdAt: "",
  updatedAt: "",
};

const bgClass: Record<string, string> = {
  cream: "bg-[#faf9f7] text-neutral-950",
  white: "bg-white text-neutral-950",
  dark: "bg-[#111111] text-white",
  motion: "bg-motion text-neutral-950",
};

const shapeClass: Record<string, string> = {
  circle: "rounded-full",
  rounded: "rounded-2xl",
  square: "rounded-none",
};

function slugFromEmail(value: string) {
  return value.split("@")[0].toLowerCase().replace(/[^a-z0-9-]/g, "") || "page";
}

function draftFrom(email: string): Page {
  const slug = slugFromEmail(email);
  return { ...emptyPage, slug, displayName: slug };
}

const ZERO_ID = "00000000-0000-0000-0000-000000000000";

function isLivePage(value: Page | null | undefined): value is Page {
  return Boolean(value && value.slug && value.id && value.id !== ZERO_ID);
}

export default function DashboardPage() {
  const [token, setTok] = useState<string | null>(null);
  const [gate, setGate] = useState<Gate>("login");
  const [panel, setPanel] = useState<Panel>("identity");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [resetToken, setResetToken] = useState("");
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [saving, setSaving] = useState(false);
  const [page, setPage] = useState<Page>(emptyPage);
  const [links, setLinks] = useState<LinkItem[]>([]);
  const [totalClicks, setTotalClicks] = useState(0);
  const [linkTitle, setLinkTitle] = useState("");
  const [linkUrl, setLinkUrl] = useState("");

  const saved = isLivePage(page);
  const canOpen = saved && Boolean(page.slug);

  async function boot(next: string, fallbackEmail = "") {
    setError("");
    const [mine, items, stats] = await Promise.all([
      synclink.getMyPage(next),
      synclink.listLinks(next),
      synclink.meStats(next).catch(() => ({ totalClicks: 0, links: [] as { id: string; clicks: number }[] })),
    ]);
    const clicksById = Object.fromEntries((stats.links || []).map((row) => [row.id, row.clicks]));
    setTotalClicks(stats.totalClicks || items.reduce((sum, item) => sum + (item.clicks || clicksById[item.id] || 0), 0));
    if (isLivePage(mine)) {
      setPage({
        ...emptyPage,
        ...mine,
        avatarShape: mine.avatarShape || "circle",
        accentColor: mine.accentColor || "#111111",
        background: mine.background || "cream",
        motion: mine.motion || "subtle",
      });
      setLinks(items.map((item) => ({ ...item, clicks: item.clicks ?? clicksById[item.id] })));
      return;
    }
    setPage(draftFrom(fallbackEmail || email));
    setLinks([]);
  }

  useEffect(() => {
    const savedToken = getToken();
    if (!savedToken) return;
    setTok(savedToken);
    void boot(savedToken).catch((err) => {
      setError(err instanceof Error ? err.message : "Could not load page");
      setPage(emptyPage);
      setLinks([]);
    });
  }, []);

  async function onAuth(event: FormEvent) {
    event.preventDefault();
    setError("");
    setNotice("");
    try {
      const res = gate === "register" ? await synclink.register(email, password) : await synclink.login(email, password);
      setToken(res.token);
      setTok(res.token);
      await boot(res.token, res.user.email || email);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not sign in");
    }
  }

  async function onForgot(event: FormEvent) {
    event.preventDefault();
    setError("");
    try {
      const res = await synclink.forgotPassword(email);
      setResetToken(res.resetToken || "");
      setNotice(res.resetToken ? `Reset token: ${res.resetToken}` : "If that email exists, a reset was issued.");
      setGate("reset");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not start reset");
    }
  }

  async function onReset(event: FormEvent) {
    event.preventDefault();
    setError("");
    try {
      await synclink.resetPassword(email, resetToken, newPassword);
      setNotice("Password reset. Sign in.");
      setGate("login");
      setPassword("");
      setNewPassword("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not reset password");
    }
  }

  async function savePage(): Promise<Page | null> {
    if (!token) return null;
    if (!page.slug.trim() || !page.displayName.trim()) {
      setError("Slug and display name are required to publish.");
      setPanel("identity");
      return null;
    }
    setSaving(true);
    setError("");
    try {
      const next = await synclink.upsertPage(token, {
        slug: page.slug.trim(),
        displayName: page.displayName.trim(),
        bio: page.bio,
        avatarUrl: page.avatarUrl,
        theme: page.theme,
        avatarShape: page.avatarShape,
        accentColor: page.accentColor,
        background: page.background,
        motion: page.motion,
      });
      setPage(next);
      setNotice(`Live at /${next.slug}`);
      return next;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not save page");
      return null;
    } finally {
      setSaving(false);
    }
  }

  async function onSavePage(event: FormEvent) {
    event.preventDefault();
    await savePage();
  }

  async function addLink(event: FormEvent) {
    event.preventDefault();
    if (!token) return;
    setError("");
    try {
      let live = page;
      if (!live.id) {
        const published = await savePage();
        if (!published) return;
        live = published;
      }
      const created = await synclink.createLink(token, { title: linkTitle, url: linkUrl });
      setLinks((current) => [...current, created]);
      setLinkTitle("");
      setLinkUrl("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not add link");
    }
  }

  async function patchLink(id: string, input: { title?: string; url?: string; active?: boolean }) {
    if (!token) return;
    try {
      const next = await synclink.updateLink(token, id, input);
      setLinks((current) => current.map((item) => (item.id === next.id ? next : item)));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not update link");
    }
  }

  async function move(id: string, dir: -1 | 1) {
    if (!token) return;
    const i = links.findIndex((item) => item.id === id);
    const j = i + dir;
    if (i < 0 || j < 0 || j >= links.length) return;
    const next = [...links];
    [next[i], next[j]] = [next[j], next[i]];
    setLinks(next);
    await synclink.reorderLinks(
      token,
      next.map((item) => item.id),
    );
  }

  async function changePassword(event: FormEvent) {
    event.preventDefault();
    if (!token) return;
    setError("");
    try {
      await synclink.changePassword(token, currentPassword, newPassword);
      setNotice("Password updated.");
      setCurrentPassword("");
      setNewPassword("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not change password");
    }
  }

  function signOut() {
    setToken(null);
    setTok(null);
    setPage(emptyPage);
    setLinks([]);
  }

  const shape = shapeClass[page.avatarShape] || shapeClass.circle;
  const tone = bgClass[page.background] || bgClass.cream;
  const dark = page.background === "dark";
  const activeLinks = links.filter((item) => item.active);

  if (!token) {
    return (
      <main className="page-enter min-h-full bg-[#faf9f7]">
        <SiteNav />
        <div className="mx-auto grid w-full max-w-5xl items-center gap-10 px-6 py-16 md:grid-cols-2">
        <div>
          <p className="text-xs tracking-[0.28em] text-neutral-400">SYNCLINK</p>
          <h1 className="mt-3 text-4xl font-medium tracking-tight">Studio</h1>
          <p className="mt-3 text-sm leading-6 text-neutral-600">Sign in. Preview stays on the right. Tap a card to edit.</p>
          <Card className="mt-8 border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <CardTitle className="font-medium">{gate === "register" ? "Create account" : gate === "forgot" || gate === "reset" ? "Reset password" : "Sign in"}</CardTitle>
            </CardHeader>
            <CardContent>
              {gate === "forgot" ? (
                <form className="space-y-4" onSubmit={onForgot}>
                  <div className="space-y-2"><Label htmlFor="email">Email</Label><Input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required /></div>
                  {error ? <p className="text-sm text-red-600">{error}</p> : null}
                  {notice ? <p className="text-sm text-neutral-600">{notice}</p> : null}
                  <Button type="submit">Send reset</Button>
                  <button type="button" className="block text-sm text-neutral-500 underline" onClick={() => setGate("login")}>Back</button>
                </form>
              ) : gate === "reset" ? (
                <form className="space-y-4" onSubmit={onReset}>
                  <div className="space-y-2"><Label htmlFor="email">Email</Label><Input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required /></div>
                  <div className="space-y-2"><Label htmlFor="resetToken">Token</Label><Input id="resetToken" value={resetToken} onChange={(e) => setResetToken(e.target.value)} required /></div>
                  <div className="space-y-2"><Label htmlFor="newPassword">New password</Label><Input id="newPassword" type="password" minLength={8} value={newPassword} onChange={(e) => setNewPassword(e.target.value)} required /></div>
                  {error ? <p className="text-sm text-red-600">{error}</p> : null}
                  <Button type="submit">Reset</Button>
                </form>
              ) : (
                <form className="space-y-4" onSubmit={onAuth}>
                  <div className="space-y-2"><Label htmlFor="email">Email</Label><Input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required /></div>
                  <div className="space-y-2"><Label htmlFor="password">Password</Label><Input id="password" type="password" minLength={8} value={password} onChange={(e) => setPassword(e.target.value)} required /></div>
                  {error ? <p className="text-sm text-red-600">{error}</p> : null}
                  <Button type="submit">{gate === "register" ? "Create account" : "Continue"}</Button>
                  <div className="flex flex-wrap gap-3 text-sm text-neutral-500">
                    <button type="button" className="underline" onClick={() => setGate(gate === "register" ? "login" : "register")}>{gate === "register" ? "Sign in" : "Sign up"}</button>
                    <button type="button" className="underline" onClick={() => setGate("forgot")}>Forgot password</button>
                  </div>
                </form>
              )}
            </CardContent>
          </Card>
        </div>
        <div className="img-hover overflow-hidden rounded-3xl border border-neutral-200/80">
          <Image src="/stations/login.png" alt="" width={1400} height={1750} className="h-auto w-full object-cover" priority />
        </div>
        </div>
      </main>
    );
  }

  return (
    <main className="page-enter min-h-full bg-[#faf9f7]">
      <SiteNav />
      <div className="mx-auto grid w-full max-w-6xl gap-8 px-6 py-12 lg:grid-cols-[minmax(0,1fr)_360px]">
      <section className="space-y-4">
        <div className="flex items-center justify-between gap-4">
          <div>
            <p className="text-xs tracking-[0.28em] text-neutral-400">DASHBOARD</p>
            <h1 className="mt-2 text-3xl font-medium tracking-tight">Studio</h1>
            <p className="mt-1 text-sm text-neutral-500">{saved ? `Saved · /${page.slug} · ${totalClicks} clicks` : "Unsaved draft"}</p>
          </div>
          <div className="flex gap-2">
            {canOpen ? <Link href={`/${page.slug}`} className="text-sm underline">Open live</Link> : null}
            <button type="button" className="text-sm text-neutral-500 underline" onClick={signOut}>Sign out</button>
          </div>
        </div>
        {error ? <p className="text-sm text-red-600">{error}</p> : null}
        {notice ? <p className="text-sm text-neutral-600">{notice}</p> : null}

        {panel === "identity" ? (
          <Card className="border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <CardTitle className="font-medium">Identity</CardTitle>
              <CardDescription>Slug, name, bio, avatar. This is what people see first.</CardDescription>
            </CardHeader>
            <CardContent>
              <form className="space-y-4" onSubmit={onSavePage}>
                <div className="space-y-2"><Label htmlFor="slug">Slug</Label><Input id="slug" value={page.slug} onChange={(e) => setPage({ ...page, slug: e.target.value })} required /></div>
                <div className="space-y-2"><Label htmlFor="displayName">Display name</Label><Input id="displayName" value={page.displayName} onChange={(e) => setPage({ ...page, displayName: e.target.value })} required /></div>
                <div className="space-y-2"><Label htmlFor="bio">Bio</Label><Textarea id="bio" value={page.bio} onChange={(e) => setPage({ ...page, bio: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="avatarUrl">Avatar URL</Label><Input id="avatarUrl" value={page.avatarUrl || ""} onChange={(e) => setPage({ ...page, avatarUrl: e.target.value || null })} /></div>
                <Button type="submit" disabled={saving}>{saving ? "Saving…" : saved ? "Save identity" : "Publish page"}</Button>
              </form>
            </CardContent>
          </Card>
        ) : null}

        {panel === "look" ? (
          <Card className="border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <CardTitle className="font-medium">Look</CardTitle>
              <CardDescription>Preview updates as you tap. Save to persist.</CardDescription>
            </CardHeader>
            <CardContent>
              <form className="space-y-4" onSubmit={onSavePage}>
                <div className="space-y-2">
                  <Label>Avatar shape</Label>
                  <div className="flex flex-wrap gap-2">
                    {["circle", "rounded", "square"].map((value) => (
                      <Button key={value} type="button" variant={page.avatarShape === value ? "default" : "outline"} onClick={() => setPage({ ...page, avatarShape: value })}>{value}</Button>
                    ))}
                  </div>
                </div>
                <div className="space-y-2"><Label htmlFor="accentColor">Accent</Label><Input id="accentColor" type="color" value={page.accentColor || "#111111"} onChange={(e) => setPage({ ...page, accentColor: e.target.value })} /></div>
                <div className="space-y-2">
                  <Label>Background</Label>
                  <div className="flex flex-wrap gap-2">
                    {["cream", "white", "dark", "motion"].map((value) => (
                      <Button key={value} type="button" variant={page.background === value ? "default" : "outline"} onClick={() => setPage({ ...page, background: value })}>{value}</Button>
                    ))}
                  </div>
                </div>
                <div className="space-y-2">
                  <Label>Motion</Label>
                  <div className="flex flex-wrap gap-2">
                    {["none", "subtle", "lively"].map((value) => (
                      <Button key={value} type="button" variant={page.motion === value ? "default" : "outline"} onClick={() => setPage({ ...page, motion: value })}>{value}</Button>
                    ))}
                  </div>
                </div>
                <Button type="submit" disabled={saving}>{saving ? "Saving…" : "Save look"}</Button>
              </form>
            </CardContent>
          </Card>
        ) : null}

        {panel === "links" ? (
          <Card className="border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <CardTitle className="font-medium">Links</CardTitle>
              <CardDescription>{saved ? "Add, edit, hide, reorder." : "Publishes the page first, then adds the link."}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <form className="grid gap-3 md:grid-cols-[1fr_1fr_auto]" onSubmit={addLink}>
                <Input placeholder="Title" value={linkTitle} onChange={(e) => setLinkTitle(e.target.value)} required />
                <Input placeholder="https://" type="url" value={linkUrl} onChange={(e) => setLinkUrl(e.target.value)} required />
                <Button type="submit">Add</Button>
              </form>
              <Separator />
              <ul className="space-y-3">
                {links.map((item, index) => (
                  <li key={item.id} className="space-y-2 rounded-2xl border border-neutral-200 px-4 py-3">
                    <Input defaultValue={item.title} onBlur={(e) => e.target.value !== item.title && void patchLink(item.id, { title: e.target.value })} />
                    <Input defaultValue={item.url} onBlur={(e) => e.target.value !== item.url && void patchLink(item.id, { url: e.target.value })} />
                    <p className="text-xs text-neutral-500">{item.clicks ?? 0} taps{item.lastClickedAt ? ` · last ${new Date(item.lastClickedAt).toLocaleString("tr-TR", { day: "2-digit", month: "short", hour: "2-digit", minute: "2-digit" })}` : ""}</p>
                    <div className="flex flex-wrap gap-1">
                      <Button type="button" size="sm" variant="ghost" disabled={index === 0} onClick={() => void move(item.id, -1)}>Up</Button>
                      <Button type="button" size="sm" variant="ghost" disabled={index === links.length - 1} onClick={() => void move(item.id, 1)}>Down</Button>
                      <Button type="button" size="sm" variant="ghost" onClick={() => void patchLink(item.id, { active: !item.active })}>{item.active ? "Hide" : "Show"}</Button>
                      <Button type="button" size="sm" variant="ghost" onClick={() => token && synclink.deleteLink(token, item.id).then(() => setLinks((c) => c.filter((l) => l.id !== item.id)))}>Delete</Button>
                    </div>
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        ) : null}

        {panel === "account" ? (
          <Card className="border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <CardTitle className="font-medium">Password</CardTitle>
            </CardHeader>
            <CardContent>
              <form className="space-y-4" onSubmit={changePassword}>
                <div className="space-y-2"><Label htmlFor="currentPassword">Current</Label><Input id="currentPassword" type="password" value={currentPassword} onChange={(e) => setCurrentPassword(e.target.value)} required /></div>
                <div className="space-y-2"><Label htmlFor="newPassword">New</Label><Input id="newPassword" type="password" minLength={8} value={newPassword} onChange={(e) => setNewPassword(e.target.value)} required /></div>
                <Button type="submit">Update</Button>
              </form>
            </CardContent>
          </Card>
        ) : null}
      </section>

      <aside className="order-first space-y-4 lg:order-none lg:sticky lg:top-8 lg:self-start">
        <button type="button" onClick={() => setPanel("identity")} className={`w-full overflow-hidden rounded-3xl border text-left ${tone} ${panel === "identity" ? "ring-2 ring-neutral-900" : "border-neutral-200"}`}>
          <div className="px-5 py-6">
            <p className="text-[10px] tracking-[0.24em] opacity-50">{saved ? "LIVE PREVIEW" : "DRAFT PREVIEW"}</p>
            <div className="mt-4 flex items-center gap-3">
              <Avatar className={`size-12 border ${shape} ${dark ? "border-white/15" : "border-neutral-200"}`}>
                {page.avatarUrl ? <AvatarImage src={page.avatarUrl} alt="" /> : null}
                <AvatarFallback>{(page.displayName || "S").slice(0, 1).toUpperCase()}</AvatarFallback>
              </Avatar>
              <div>
                <p className="text-sm font-medium">{page.displayName || "Your name"}</p>
                <p className="text-xs opacity-60">/{page.slug || "slug"}</p>
              </div>
            </div>
            <p className="mt-3 text-xs leading-5 opacity-70">{page.bio || "Bio appears here."}</p>
            <ul className="mt-4 space-y-2">
              {activeLinks.map((item) => (
                <li key={item.id} className={`rounded-xl border px-3 py-2 text-center text-xs ${dark ? "border-white/15" : "border-neutral-200 bg-white/80"}`} style={{ boxShadow: `0 0 0 1px ${page.accentColor || "#111111"}14` }}>{item.title}{typeof item.clicks === "number" ? ` · ${item.clicks}` : ""}</li>
              ))}
              {activeLinks.length === 0 ? <li className="text-xs opacity-50">No links yet</li> : null}
            </ul>
          </div>
        </button>
        <Card className={`cursor-pointer border-neutral-200/80 bg-white shadow-none transition hover:-translate-y-0.5 ${panel === "look" ? "ring-2 ring-neutral-900" : ""}`} onClick={() => setPanel("look")}>
          <CardHeader>
            <CardTitle className="text-base font-medium">Look</CardTitle>
            <CardDescription>{page.avatarShape} · {page.background} · {page.motion}</CardDescription>
          </CardHeader>
        </Card>
        <Card className={`cursor-pointer border-neutral-200/80 bg-white shadow-none transition hover:-translate-y-0.5 ${panel === "links" ? "ring-2 ring-neutral-900" : ""}`} onClick={() => setPanel("links")}>
          <CardHeader>
            <CardTitle className="text-base font-medium">Links</CardTitle>
            <CardDescription>{links.length} on this page</CardDescription>
          </CardHeader>
        </Card>
        <Card className={`cursor-pointer border-neutral-200/80 bg-white shadow-none transition hover:-translate-y-0.5 ${panel === "account" ? "ring-2 ring-neutral-900" : ""}`} onClick={() => setPanel("account")}>
          <CardHeader>
            <CardTitle className="text-base font-medium">Account</CardTitle>
            <CardDescription>Change password</CardDescription>
          </CardHeader>
        </Card>
      </aside>
      </div>
    </main>
  );
}
