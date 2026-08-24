"use client";

import Link from "next/link";
import { FormEvent, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle, SheetTrigger } from "@/components/ui/sheet";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Textarea } from "@/components/ui/textarea";
import { SiteNav } from "@/components/site-nav";
import { getToken, synclink, type AdminStats, type AdminUser, type Page, type PlatformSettings } from "@/lib/api";

type Tab = "overview" | "users" | "pages" | "settings";

const emptySettings: PlatformSettings = {
  siteName: "SyncLink",
  tagline: "",
  about: "",
  supportEmail: "",
  signupEnabled: true,
  maintenance: false,
  metaTitle: "SyncLink",
  metaDescription: "One page. Every link.",
  ogImage: "",
  favicon: "",
  themeColor: "#111111",
  heroTitle: "",
  heroSubtitle: "",
  heroCta: "",
  heroCtaHref: "/dashboard",
  heroImage: "",
  demoSlug: "gurkan",
  nav: [
    { label: "Home", href: "/" },
    { label: "About", href: "/about" },
    { label: "Dashboard", href: "/dashboard" },
    { label: "Admin", href: "/admin" },
  ],
};

export default function AdminPage() {
  const [token, setTok] = useState<string | null>(null);
  const [tab, setTab] = useState<Tab>("overview");
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [me, setMe] = useState<AdminUser | null>(null);
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [pages, setPages] = useState<Page[]>([]);
  const [settings, setSettings] = useState<PlatformSettings>(emptySettings);

  async function boot(next: string) {
    try {
      const self = await synclink.adminMe(next);
      setMe(self);
      const [st, us, pg, se] = await Promise.all([
        synclink.adminStats(next),
        synclink.adminUsers(next),
        synclink.adminPages(next),
        synclink.adminSettings(next),
      ]);
      setStats(st);
      setUsers(us);
      setPages(pg);
      setSettings({ ...emptySettings, ...se });
      setError("");
    } catch (err) {
      setMe(null);
      setError(err instanceof Error ? err.message : "Not an admin");
    }
  }

  useEffect(() => {
    const existing = getToken();
    setTok(existing);
    if (existing) void boot(existing);
  }, []);

  async function saveSettings(event: FormEvent) {
    event.preventDefault();
    if (!token) return;
    try {
      setSettings(await synclink.adminPutSettings(token, settings));
      setNotice("Settings saved.");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Save failed");
    }
  }

  async function setStatus(user: AdminUser, status: string) {
    if (!token) return;
    const next = await synclink.adminPatchUser(token, user.id, { status });
    setUsers((list) => list.map((item) => (item.id === next.id ? next : item)));
  }

  async function removeUser(id: string) {
    if (!token) return;
    await synclink.adminDeleteUser(token, id);
    setUsers((list) => list.filter((item) => item.id !== id));
  }

  const pageClicks = [...(stats?.pageClicks ?? [])].sort((a, b) => {
    if (b.clicks !== a.clicks) return b.clicks - a.clicks;
    return a.slug.localeCompare(b.slug);
  });

  if (!token) {
    return (
      <main className="min-h-full bg-[#faf9f7]">
        <SiteNav />
        <div className="flex flex-col items-center justify-center px-6 py-20">
          <p className="mb-4 text-xs tracking-[0.28em] text-neutral-400">ADMIN</p>
          <p className="mb-6 text-neutral-600">Sign in first. First account is admin.</p>
          <Button asChild><Link href="/dashboard">Go to studio</Link></Button>
        </div>
      </main>
    );
  }

  if (error && !me) {
    return (
      <main className="min-h-full bg-[#faf9f7]">
        <SiteNav />
        <div className="flex flex-col items-center justify-center px-6 py-20">
          <p className="mb-4 text-xs tracking-[0.28em] text-neutral-400">ADMIN</p>
          <p className="mb-6 text-neutral-600">{error}</p>
          <Button asChild variant="outline"><Link href="/dashboard">Back</Link></Button>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-full bg-[#faf9f7]">
      <SiteNav settings={settings} />
      <header className="mx-auto flex w-full max-w-6xl flex-wrap items-center justify-between gap-3 px-6 py-6">
        <p className="text-xs tracking-[0.28em] text-neutral-500">SYNCLINK ADMIN</p>
        <div className="flex flex-wrap gap-2">
          {(["overview", "users", "pages", "settings"] as Tab[]).map((item) => (
            <Button key={item} size="sm" variant={tab === item ? "default" : "outline"} onClick={() => setTab(item)}>
              {item}
            </Button>
          ))}
          <Button asChild size="sm" variant="ghost"><Link href="/dashboard">Studio</Link></Button>
        </div>
      </header>

      <div className="mx-auto w-full max-w-6xl space-y-6 px-6 pb-16">
        {notice ? <p className="text-sm text-neutral-600">{notice}</p> : null}
        {error ? <p className="text-sm text-red-600">{error}</p> : null}

        {tab === "overview" ? (
          <div className="space-y-6">
            <div className="grid gap-4 sm:grid-cols-3">
              <Card className="border-neutral-200/80 bg-white shadow-none">
                <CardHeader>
                  <CardTitle className="font-medium">Users</CardTitle>
                  <CardDescription>Signed-up accounts</CardDescription>
                </CardHeader>
                <CardContent className="text-4xl">{stats?.users ?? 0}</CardContent>
              </Card>
              <Card className="border-neutral-200/80 bg-white shadow-none">
                <CardHeader>
                  <CardTitle className="font-medium">Pages</CardTitle>
                  <CardDescription>Public SyncLink pages</CardDescription>
                </CardHeader>
                <CardContent className="text-4xl">{stats?.pages ?? 0}</CardContent>
              </Card>
              <Card className="border-neutral-200/80 bg-white shadow-none">
                <CardHeader>
                  <CardTitle className="font-medium">Clicks</CardTitle>
                  <CardDescription>All public taps</CardDescription>
                </CardHeader>
                <CardContent className="text-4xl">{stats?.totalClicks ?? 0}</CardContent>
              </Card>
            </div>
            <Card className="border-neutral-200/80 bg-white shadow-none">
              <CardHeader>
                <CardTitle className="font-medium">Clicks by page</CardTitle>
                <CardDescription>Public taps on each page</CardDescription>
              </CardHeader>
              <CardContent>
                {pageClicks.length === 0 ? (
                  <p className="text-xs text-neutral-500">No pages yet.</p>
                ) : (
                  <Table className="text-xs">
                    <TableHeader>
                      <TableRow>
                        <TableHead>Slug</TableHead>
                        <TableHead className="text-right">Clicks</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {pageClicks.map((row) => (
                        <TableRow key={row.id || row.slug}>
                          <TableCell>
                            <Link href={`/${row.slug}`}>/{row.slug}</Link>
                          </TableCell>
                          <TableCell className="text-right">{row.clicks}</TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </CardContent>
            </Card>
          </div>
        ) : null}

        {tab === "users" ? (
          <Card className="border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <CardTitle className="font-medium">People</CardTitle>
              <CardDescription>Disable or remove accounts. Table on desktop, sheets on mobile.</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="hidden md:block">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Email</TableHead>
                      <TableHead>Role</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead></TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {users.map((user) => (
                      <TableRow key={user.id}>
                        <TableCell>{user.email}</TableCell>
                        <TableCell>{user.role}</TableCell>
                        <TableCell>{user.status}</TableCell>
                        <TableCell className="text-right">
                          <Button size="sm" variant="outline" onClick={() => void setStatus(user, user.status === "active" ? "disabled" : "active")}>
                            {user.status === "active" ? "Disable" : "Enable"}
                          </Button>
                          <Button size="sm" variant="ghost" onClick={() => void removeUser(user.id)}>Delete</Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
              <ul className="space-y-3 md:hidden">
                {users.map((user) => (
                  <li key={user.id}>
                    <Sheet>
                      <SheetTrigger asChild>
                        <button className="w-full rounded-2xl border border-neutral-200 bg-[#faf9f7] px-4 py-3 text-left">
                          <p>{user.email}</p>
                          <p className="text-xs text-neutral-500">{user.role} · {user.status}</p>
                        </button>
                      </SheetTrigger>
                      <SheetContent side="bottom" className="rounded-t-3xl">
                        <SheetHeader>
                          <SheetTitle>{user.email}</SheetTitle>
                          <SheetDescription>{user.role} · {user.status}</SheetDescription>
                        </SheetHeader>
                        <div className="flex gap-2 pt-4">
                          <Button onClick={() => void setStatus(user, user.status === "active" ? "disabled" : "active")}>
                            {user.status === "active" ? "Disable" : "Enable"}
                          </Button>
                          <Button variant="outline" onClick={() => void removeUser(user.id)}>Delete</Button>
                        </div>
                      </SheetContent>
                    </Sheet>
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        ) : null}

        {tab === "pages" ? (
          <Card className="border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <CardTitle className="font-medium">Pages</CardTitle>
              <CardDescription>Every public page on the platform.</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="hidden md:block">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Slug</TableHead>
                      <TableHead>Name</TableHead>
                      <TableHead>Theme</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {pages.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell><Link href={`/${item.slug}`}>/{item.slug}</Link></TableCell>
                        <TableCell>{item.displayName}</TableCell>
                        <TableCell>{item.theme}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
              <ul className="space-y-3 md:hidden">
                {pages.map((item) => (
                  <li key={item.id} className="rounded-2xl border border-neutral-200 px-4 py-3">
                    <Link href={`/${item.slug}`}>/{item.slug}</Link>
                    <p className="text-sm text-neutral-500">{item.displayName}</p>
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        ) : null}

        {tab === "settings" ? (
          <Card className="border-neutral-200/80 bg-white shadow-none">
            <CardHeader>
              <CardTitle className="font-medium">Platform</CardTitle>
              <CardDescription>Name, about, and every head tag. The public site reads this live.</CardDescription>
            </CardHeader>
            <CardContent>
              <form className="space-y-4" onSubmit={saveSettings}>
                <div className="space-y-2"><Label htmlFor="siteName">Site name</Label><Input id="siteName" value={settings.siteName} onChange={(e) => setSettings({ ...settings, siteName: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="tagline">Tagline</Label><Input id="tagline" value={settings.tagline} onChange={(e) => setSettings({ ...settings, tagline: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="about">About</Label><Textarea id="about" value={settings.about} onChange={(e) => setSettings({ ...settings, about: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="supportEmail">Support email</Label><Input id="supportEmail" type="email" value={settings.supportEmail} onChange={(e) => setSettings({ ...settings, supportEmail: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="metaTitle">Meta title</Label><Input id="metaTitle" value={settings.metaTitle} onChange={(e) => setSettings({ ...settings, metaTitle: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="metaDescription">Meta description</Label><Textarea id="metaDescription" value={settings.metaDescription} onChange={(e) => setSettings({ ...settings, metaDescription: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="ogImage">OG image URL</Label><Input id="ogImage" value={settings.ogImage} onChange={(e) => setSettings({ ...settings, ogImage: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="favicon">Favicon URL</Label><Input id="favicon" value={settings.favicon} onChange={(e) => setSettings({ ...settings, favicon: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="themeColor">Theme color</Label><Input id="themeColor" value={settings.themeColor} onChange={(e) => setSettings({ ...settings, themeColor: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="heroTitle">Hero title</Label><Input id="heroTitle" value={settings.heroTitle || ""} onChange={(e) => setSettings({ ...settings, heroTitle: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="heroSubtitle">Hero subtitle</Label><Textarea id="heroSubtitle" value={settings.heroSubtitle || ""} onChange={(e) => setSettings({ ...settings, heroSubtitle: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="heroCta">Hero CTA</Label><Input id="heroCta" value={settings.heroCta || ""} onChange={(e) => setSettings({ ...settings, heroCta: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="heroCtaHref">Hero CTA href</Label><Input id="heroCtaHref" value={settings.heroCtaHref || ""} onChange={(e) => setSettings({ ...settings, heroCtaHref: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="heroImage">Hero image URL</Label><Input id="heroImage" value={settings.heroImage || ""} onChange={(e) => setSettings({ ...settings, heroImage: e.target.value })} /></div>
                <div className="space-y-2"><Label htmlFor="demoSlug">Demo slug</Label><Input id="demoSlug" value={settings.demoSlug || ""} onChange={(e) => setSettings({ ...settings, demoSlug: e.target.value })} /></div>
                <div className="space-y-2">
                  <Label>Top nav</Label>
                  {(settings.nav || []).map((item, index) => (
                    <div key={index} className="grid gap-2 md:grid-cols-[1fr_1fr_auto]">
                      <Input placeholder="Label" value={item.label} onChange={(e) => {
                        const nav = [...(settings.nav || [])];
                        nav[index] = { ...nav[index], label: e.target.value };
                        setSettings({ ...settings, nav });
                      }} />
                      <Input placeholder="/path" value={item.href} onChange={(e) => {
                        const nav = [...(settings.nav || [])];
                        nav[index] = { ...nav[index], href: e.target.value };
                        setSettings({ ...settings, nav });
                      }} />
                      <Button type="button" variant="ghost" onClick={() => setSettings({ ...settings, nav: (settings.nav || []).filter((_, i) => i !== index) })}>Remove</Button>
                    </div>
                  ))}
                  <Button type="button" variant="outline" onClick={() => setSettings({ ...settings, nav: [...(settings.nav || []), { label: "", href: "" }] })}>Add nav link</Button>
                </div>
                <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={settings.signupEnabled} onChange={(e) => setSettings({ ...settings, signupEnabled: e.target.checked })} /> Signup enabled</label>
                <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={settings.maintenance} onChange={(e) => setSettings({ ...settings, maintenance: e.target.checked })} /> Maintenance</label>
                <Button type="submit">Save settings</Button>
              </form>
            </CardContent>
          </Card>
        ) : null}
      </div>
    </main>
  );
}
