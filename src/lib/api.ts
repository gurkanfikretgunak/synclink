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
};

export type CreateLinkInput = {
  title: string;
  url: string;
  icon?: string | null;
  active?: boolean;
};

export type UpdateLinkInput = {
  title?: string;
  url?: string;
  icon?: string | null;
  active?: boolean;
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
    throw new Error(data.message || data.error || res.statusText);
  }
  return data;
}

export const synclink = {
  getPublicPage(slug: string) {
    return request<PublicPage>(`/api/v1/public/pages/${encodeURIComponent(slug)}`);
  },
  recordClick(slug: string, id: string) {
    return request<{ ok: boolean }>(
      `/api/v1/public/pages/${encodeURIComponent(slug)}/links/${encodeURIComponent(id)}/click`,
      { method: "POST" },
    ).catch(() => ({ ok: false }));
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
  adminStats(token: string) {
    return request<{ users: number; pages: number }>("/api/v1/admin/stats", { token });
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
