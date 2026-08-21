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
  createdAt: string;
  updatedAt: string;
};

export type PublicLink = {
  id: string;
  title: string;
  url: string;
  icon: string | null;
  order: number;
};

export type PublicPage = {
  slug: string;
  displayName: string;
  bio: string;
  avatarUrl: string | null;
  theme: string;
  links: PublicLink[];
};

export type UpsertPageInput = {
  slug: string;
  displayName: string;
  bio: string;
  avatarUrl?: string | null;
  theme?: string;
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

export type LoginResponse = {
  token: string;
  user: {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    status: string;
    created_at: string;
  };
};

type ApiError = {
  error?: string;
  message?: string;
};

async function request<T>(
  path: string,
  init: RequestInit & { token?: string | null } = {},
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
    throw new Error(data.message || data.error || res.statusText);
  }
  return data;
}

export const synclink = {
  getPublicPage(slug: string) {
    return request<PublicPage>(`/api/v1/public/pages/${encodeURIComponent(slug)}`);
  },
  login(email: string, password: string) {
    return request<LoginResponse>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },
  getMyPage(token: string) {
    return request<Page>("/api/v1/me/page", { token });
  },
  upsertPage(token: string, input: UpsertPageInput) {
    return request<Page>("/api/v1/me/page", {
      method: "PUT",
      token,
      body: JSON.stringify(input),
    });
  },
  listLinks(token: string) {
    return request<LinkItem[]>("/api/v1/me/page/links", { token });
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
};
