const TOKEN_KEY = "synclink.token";

export function apiBase(): string {
  const url = process.env.NEXT_PUBLIC_MF_API_URL;
  if (!url) {
    throw new Error("NEXT_PUBLIC_MF_API_URL is not set");
  }
  return url.replace(/\/$/, "");
}

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string | null): void {
  if (typeof window === "undefined") return;
  if (token) window.localStorage.setItem(TOKEN_KEY, token);
  else window.localStorage.removeItem(TOKEN_KEY);
}

export const SOCIAL_NETWORKS = [
  "instagram",
  "x",
  "youtube",
  "tiktok",
  "github",
  "linkedin",
  "threads",
  "spotify",
  "whatsapp",
  "website",
  "email",
] as const;

export type Social = {
  network: string;
  url: string;
};

export type Subscriber = {
  id: string;
  email: string;
  createdAt: string;
};

export type Page = {
  id: string;
  slug: string;
  displayName: string;
  bio: string;
  avatarUrl: string | null;
  theme: string;
  avatarShape: string;
  accentColor: string;
  background: string;
  motion: string;
  socials?: Social[];
  verified?: boolean;
  pagePassword?: string;
  publishedAt?: string | null;
  coverUrl?: string | null;
  coverKind?: string;
  createdAt: string;
  updatedAt: string;
};

export type LinkItem = {
  id: string;
  title: string;
  url: string;
  icon: string | null;
  order: number;
  active: boolean;
  clicks?: number;
  lastClickedAt?: string | null;
  featured?: boolean;
  thumbnailUrl?: string | null;
  startsAt?: string | null;
  endsAt?: string | null;
  sensitive?: boolean;
  section?: string | null;
  embedUrl?: string | null;
  createdAt: string;
  updatedAt: string;
};

export type PublicLink = {
  id: string;
  title: string;
  url: string;
  icon: string | null;
  order: number;
  clicks?: number;
  lastClickedAt?: string | null;
  featured?: boolean;
  thumbnailUrl?: string | null;
  sensitive?: boolean;
  section?: string | null;
  embedUrl?: string | null;
};

export type PublicPage = {
  slug: string;
  displayName: string;
  bio: string;
  avatarUrl: string | null;
  theme: string;
  avatarShape: string;
  accentColor: string;
  background: string;
  motion: string;
  socials?: Social[];
  verified?: boolean;
  publishedAt?: string | null;
  coverUrl?: string | null;
  coverKind?: string;
  links: PublicLink[];
};

export type UpsertPageInput = {
  slug: string;
  displayName: string;
  bio: string;
  avatarUrl?: string | null;
  theme?: string;
  avatarShape?: string;
  accentColor?: string;
  background?: string;
  motion?: string;
  socials?: Social[];
  pagePassword?: string;
  coverUrl?: string | null;
  coverKind?: string;
};

export type CreateLinkInput = {
  title: string;
  url: string;
  icon?: string | null;
  active?: boolean;
  featured?: boolean;
  thumbnailUrl?: string | null;
  startsAt?: string | null;
  endsAt?: string | null;
  sensitive?: boolean;
  section?: string | null;
  embedUrl?: string | null;
};

export type UpdateLinkInput = {
  title?: string;
  url?: string;
  icon?: string | null;
  active?: boolean;
  featured?: boolean;
  thumbnailUrl?: string | null;
  startsAt?: string | null;
  endsAt?: string | null;
  sensitive?: boolean;
  section?: string | null;
  embedUrl?: string | null;
};

export type AdminUser = {
  id: string;
  email: string;
  role: string;
  status: string;
  createdAt: string;
};

export type NavItem = {
  label: string;
  href: string;
};

export type PlatformSettings = {
  siteName: string;
  tagline: string;
  about: string;
  supportEmail: string;
  signupEnabled: boolean;
  maintenance: boolean;
  metaTitle: string;
  metaDescription: string;
  ogImage: string;
  favicon: string;
  themeColor: string;
  heroTitle: string;
  heroSubtitle: string;
  heroCta: string;
  heroCtaHref: string;
  heroImage: string;
  demoSlug: string;
  nav: NavItem[];
};

export type MeStats = {
  totalClicks: number;
  links: { id: string; title?: string; clicks: number }[];
};

export type AdminStats = {
  users: number;
  pages: number;
  totalClicks?: number;
};

export type LoginResponse = {
  token: string;
  user: {
    id: string;
    email: string;
  };
};

type ApiError = {
  error?: string;
  message?: string;
  code?: number;
};

function isNotFound(data: ApiError, status: number): boolean {
  if (status === 404) return true;
  const text = `${data.error || ""} ${data.message || ""}`.toLowerCase();
  return text.includes("not found");
}

function humanApiError(data: ApiError, status: number, statusText: string): string {
  const message = typeof data.message === "string" ? data.message.trim() : "";
  const error = typeof data.error === "string" ? data.error.trim() : "";

  if (status === 404) {
    return message || error || statusText;
  }

  if (status === 401) {
    if (error.toLowerCase() === "locked" || message.toLowerCase() === "locked") {
      return "locked";
    }
    const generic = !message || /^unauthori[sz]ed$/i.test(message) || message.toLowerCase() === "validation";
    return generic ? "invalid credentials" : message;
  }

  const genericValidation =
    message.toLowerCase() === "validation" || error === "Unprocessable Entity";
  if (genericValidation && (!message || message.toLowerCase() === "validation")) {
    return "Check email, password (8+), slug, or URL.";
  }

  return message || error || statusText;
}

async function request<T>(
  path: string,
  init: RequestInit & { token?: string | null; allow404?: boolean } = {},
): Promise<T> {
  const headers = new Headers(init.headers);
  if (init.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (init.token) {
    headers.set("Authorization", `Bearer ${init.token}`);
  }

  const res = await fetch(`${apiBase()}${path}`, {
    ...init,
    headers,
    cache: "no-store",
  });

  if (res.status === 204) {
    return undefined as T;
  }

  const text = await res.text();
  const data = text ? (JSON.parse(text) as T & ApiError) : ({} as T & ApiError);
  if (!res.ok) {
    if (init.allow404 && isNotFound(data, res.status)) {
      return null as T;
    }
    throw new Error(humanApiError(data, res.status, res.statusText));
  }
  return data;
}

export const synclink = {
  getPublicPage(slug: string, pagePassword?: string) {
    return request<PublicPage>(`/api/v1/public/pages/${encodeURIComponent(slug)}`, {
      headers: pagePassword ? { "X-Page-Password": pagePassword } : undefined,
    });
  },
  subscribe(slug: string, email: string) {
    return request<{ ok: boolean }>(`/api/v1/public/pages/${encodeURIComponent(slug)}/subscribe`, {
      method: "POST",
      body: JSON.stringify({ email }),
    });
  },
  listSubscribers(token: string) {
    return request<Subscriber[] | null>("/api/v1/me/subscribers", { token, allow404: true }).then(
      (items) => items || [],
    );
  },
  deleteSubscriber(token: string, id: string) {
    return request<void>(`/api/v1/me/subscribers/${encodeURIComponent(id)}`, {
      method: "DELETE",
      token,
    });
  },
  recordClick(slug: string, id: string) {
    return request<{ ok: boolean; clicks: number }>(
      `/api/v1/public/pages/${encodeURIComponent(slug)}/links/${encodeURIComponent(id)}/click`,
      { method: "POST" },
    );
  },
  login(email: string, password: string) {
    return request<LoginResponse>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },
  register(email: string, password: string) {
    return request<LoginResponse>("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },
  forgotPassword(email: string) {
    return request<{ ok: boolean; resetToken?: string; expiresIn?: number }>(
      "/api/v1/auth/forgot-password",
      { method: "POST", body: JSON.stringify({ email }) },
    );
  },
  resetPassword(email: string, token: string, newPassword: string) {
    return request<{ ok: boolean }>("/api/v1/auth/reset-password", {
      method: "POST",
      body: JSON.stringify({ email, token, newPassword }),
    });
  },
  changePassword(token: string, currentPassword: string, newPassword: string) {
    return request<{ ok: boolean }>("/api/v1/me/password", {
      method: "PUT",
      token,
      body: JSON.stringify({ currentPassword, newPassword }),
    });
  },
  getMyPage(token: string) {
    return request<Page | null>("/api/v1/me/page", { token, allow404: true });
  },
  upsertPage(token: string, input: UpsertPageInput) {
    return request<Page>("/api/v1/me/page", {
      method: "PUT",
      token,
      body: JSON.stringify(input),
    });
  },
  listLinks(token: string) {
    return request<LinkItem[] | null>("/api/v1/me/page/links", { token, allow404: true }).then(
      (items) => items || [],
    );
  },
  createLink(token: string, input: CreateLinkInput) {
    return request<LinkItem>("/api/v1/me/page/links", {
      method: "POST",
      token,
      body: JSON.stringify(input),
    });
  },
  updateLink(token: string, id: string, input: UpdateLinkInput) {
    return request<LinkItem>(`/api/v1/me/page/links/${id}`, {
      method: "PATCH",
      token,
      body: JSON.stringify(input),
    });
  },
  deleteLink(token: string, id: string) {
    return request<void>(`/api/v1/me/page/links/${id}`, {
      method: "DELETE",
      token,
    });
  },
  reorderLinks(token: string, ids: string[]) {
    return request<LinkItem[]>("/api/v1/me/page/links/reorder", {
      method: "PUT",
      token,
      body: JSON.stringify({ ids }),
    });
  },
  publicSettings() {
    return request<PlatformSettings>("/api/v1/public/settings");
  },
  adminMe(token: string) {
    return request<AdminUser>("/api/v1/admin/me", { token });
  },
  meStats(token: string) {
    return request<MeStats | null>("/api/v1/me/stats", { token, allow404: true }).then(
      (data) => data || { totalClicks: 0, links: [] },
    );
  },
  adminStats(token: string) {
    return request<AdminStats>("/api/v1/admin/stats", { token });
  },
  adminUsers(token: string) {
    return request<AdminUser[]>("/api/v1/admin/users", { token });
  },
  adminPatchUser(token: string, id: string, input: { role?: string; status?: string }) {
    return request<AdminUser>(`/api/v1/admin/users/${id}`, {
      method: "PATCH",
      token,
      body: JSON.stringify(input),
    });
  },
  adminDeleteUser(token: string, id: string) {
    return request<void>(`/api/v1/admin/users/${id}`, { method: "DELETE", token });
  },
  adminPages(token: string) {
    return request<Page[]>("/api/v1/admin/pages", { token });
  },
  adminPatchPage(token: string, id: string, input: { verified?: boolean }) {
    return request<Page>(`/api/v1/admin/pages/${id}`, {
      method: "PATCH",
      token,
      body: JSON.stringify(input),
    });
  },
  adminSettings(token: string) {
    return request<PlatformSettings>("/api/v1/admin/settings", { token });
  },
  adminPutSettings(token: string, input: PlatformSettings) {
    return request<PlatformSettings>("/api/v1/admin/settings", {
      method: "PUT",
      token,
      body: JSON.stringify(input),
    });
  },
};
