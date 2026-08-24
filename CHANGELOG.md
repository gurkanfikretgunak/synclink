# Changelog

All notable changes to SyncLink live in this file.

## [Unreleased]

## [0.10.12] - 2026-08-24

### App
- Studio/public 422 copy is path-aware: page save, link save, subscribe, auth, admin settings. If the API later sends `fields`, those win. Locked 401 unchanged. No new API.

## [0.10.11] - 2026-08-24

### App
- Client `request()` maps generic 422 (`message` exactly `validation` or `error` `Unprocessable Entity`) to `Check email, password (8+), slug, or URL.` Specific API messages are kept. Generic/empty 401 becomes `invalid credentials`. 404-safe behavior is unchanged. No new API.

## [0.10.10] - 2026-08-24

### App
- Marketing home: type-first hero, device frame around `heroImage` (monogram if empty), primary + ghost CTAs, `/{demo}` handle chip, type-only feature strip. Dropped Station 01 and the three-card grid. No new API.

## [0.10.9] - 2026-08-24

### App
- Public `/{slug}` quiet copyable handle under the display name (copies the full URL). No new API.

## [0.10.8] - 2026-08-24

### App
- App Router `sitemap.ts` and `robots.ts`: `/sitemap.xml` lists `/`, `/about`, and `/${demoSlug}` from public settings (omits the slug row if settings throw). `/robots.txt` allows `/`, disallows `/dashboard` `/admin` `/reset`, points at sitemap. No new API.

## [0.10.7] - 2026-08-24

### App
- Public `/{slug}` generateMetadata: `openGraph.url` matches canonical `/{slug}`; `appleWebApp` (capable, title from displayName, default status bar) and `applicationName` for Add to Home Screen. Locked pages stay noindex. No new API.

## [0.10.6] - 2026-08-24

### App
- Public `/{slug}` generateMetadata: `themeColor` from `accentColor`, tab/Apple icon from avatar (or cover) http(s) URL, canonical `/{slug}`. Locked pages stay noindex. No new API.

## [0.10.5] - 2026-08-24

### App
- Public + studio Share card: client-side Add to contacts (.vcf from displayName, bio, page URL, socials; PHOTO if avatar is http(s)). No new API.

## [0.10.4] - 2026-08-24

### App
- Public `/{slug}` emits schema.org ProfilePage JSON-LD (Person name/bio/image, sameAs from socials). No new API.

## [0.10.3] - 2026-08-24

### App
- Public + studio Share card: native share sheet when the browser supports it, copy URL fallback. No new API.

## [0.10.2] - 2026-08-24

### App
- Public `/{slug}` sets document title, description, and og/twitter image from `displayName` / `bio` / avatar (or cover). Locked pages stay noindex. No new API.

## [0.10.1] - 2026-08-23

### App
- Studio aside Share card: QR + copy URL + PNG when a slug is set (Linktree Share → QR). No new API.

## [0.10.0] - 2026-08-23

### App
- Studio + public cover (`coverUrl`, `coverKind` image|video), link `section` (max 40), and `embedUrl`. Hidden until API 0.10.0 is live.

### Added
- Page cover: `coverUrl` (http/https) and `coverKind` (`image`|`video`, default `image` when URL is set) on GET/PUT `/me/page` and public GET `/public/pages/{slug}`.
- Link `section` (string, max 40) and `embedUrl` (http/https) on create/update and public links. Empty section is ungrouped. Invalid embed URLs drop to null.
- SQLite `pages.cover_url`, `pages.cover_kind`, `links.section`, `links.embed_url`.

## [0.9.2] - 2026-08-23

### App
- Public `/{slug}` shows `publishedAt` when the API sends it (hidden until Render recuts 0.9.1).

## [0.9.1] - 2026-08-23

### App
- Public `/{slug}` and studio Share card: QR for the live URL (copy + PNG). No new API.

### Added
- `publishedAt` on owner GET/PUT `/api/v1/me/page` and public GET `/api/v1/public/pages/{slug}` (RFC3339; null only for unsaved empty page). Set once on first create; later upserts keep it. SQLite `pages.published_at`; existing rows backfill from `created_at`.

## [0.9.0] - 2026-08-23

### App
- Public `/{slug}` email form `POST /api/v1/public/pages/{slug}/subscribe` `{email}` (201). Hides itself if the route is missing.
- Studio Inbox lists `GET /api/v1/me/subscribers` and deletes via `DELETE /api/v1/me/subscribers/{id}`.
- Locked pages: public 401 `locked` shows an unlock form that sends `X-Page-Password`.
- Studio identity can set `pagePassword`. Look presets (cream/white/dark/motion). Social network select includes whatsapp.
- Link extras in studio + public: featured, thumbnailUrl, startsAt/endsAt, sensitive (18+ confirm). Verified badge when API sets it.

### Added
- Email capture: `subscribers` table (`id`, `page_id`, `email`, `created_at`, unique `page_id+email`).
- `POST /api/v1/public/pages/{slug}/subscribe` `{email}` → `201 {"ok":true}`; `400` invalid email; `404` missing page; `409` duplicate.
- `GET /api/v1/me/subscribers` (JWT) → `[{id,email,createdAt}]`. Empty page is `[]`.
- `DELETE /api/v1/me/subscribers/{id}` (JWT) → `204`.
- Link extras on `Link` / `LinkDTO` / public links: `featured`, `thumbnailUrl`, `startsAt`, `endsAt`, `sensitive`. Create/Update accept them. Public GET omits inactive and out-of-window links (`startsAt` in the future or `endsAt` in the past).
- Page extras: `verified` (default false; owner PUT cannot set it). `PATCH /api/v1/admin/pages/{id}` `{verified}` (admin JWT).
- Optional `pagePassword` on owner GET/PUT `/me/page`. If set, public GET `/public/pages/{slug}` returns `401 {"error":"locked"}` unless `X-Page-Password` matches. Public payload never includes the password.
- Social network `whatsapp` (`https://wa.me/...`).
- SQLite ALTER for new page/link columns plus `subscribers` table.

## [0.8.2] - 2026-08-23

### Added
- Studio + public `/{slug}` social-icon row from `socials: [{ network, url }]`. Saves with the page. Hidden until the API stores the field.
- API `socials` on `Page`, `PageDTO`, `PublicPage`, and `PUT /me/page` (`[{network,url}]`). Empty/null becomes `[]`. Allowed networks: instagram, x (twitter→x), youtube, tiktok, github, linkedin, threads, spotify, website, email. HTTP(S) URLs; email also `mailto:` or `user@host`. Invalid items dropped; max 12.
- SQLite `pages.socials` TEXT JSON; ALTER on migrate. GET/PUT `/me/page` and public GET `pages/{slug}` include `socials`. Seed `/gurkan` has github + website.

## [0.8.1] - 2026-08-23

### Added
- `lastClickedAt` on `Link`, `LinkDTO`, and public links (JSON camelCase; null if never clicked). Set by `RecordClick` / `IncrementClicks` in memory and SQLite (`last_clicked_at` TEXT nullable, ALTER on migrate).
- Studio and public `/{slug}` show last tap when `lastClickedAt` is set.

## [0.8.0] - 2026-08-22

### Added
- Studio shows total clicks from GET /api/v1/me/stats (0 if the route is missing). Per-link counts merge onto the preview.
- Admin overview has a Clicks card from GET /api/v1/admin/stats `totalClicks`.
- GET /api/v1/me/stats (JWT): 200 `{"totalClicks":N,"links":[{"id":"uuid","title":"...","clicks":N,"url":"..."}]}`. No page or empty links: `{"totalClicks":0,"links":[]}`.
- GET /api/v1/admin/stats now includes click sum as `totalClicks`: `{"users":N,"pages":N,"totalClicks":N}`.
- Memory and SQLite `SumClicks`; page service `MyStats` aggregates the caller's links.


## [0.7.0] - 2026-08-22

### Added
- SQLite persistence for the Go API (`modernc.org/sqlite`, no CGO). Path from `SYNCLINK_DB`, default `./data/synclink.db` (parent directory created). Render/Docker: `SYNCLINK_DB=/var/data/synclink.db`.
- Tables: users, pages, links (`clicks INTEGER NOT NULL DEFAULT 0`), settings (single JSON row), password_reset_tokens.
- Public `POST /api/v1/public/pages/{slug}/links/{id}/click` (no JWT). Increments clicks on that active link. Response `200 {"ok":true,"clicks":N}` or `404` if the page/link is missing or inactive.
- `clicks` on `Link`, `LinkDTO`, and public links (memory and SQLite).

### Changed
- `SeedIfEmpty` / `SeedDemoIfEmpty` run only when users or pages are empty so a restart does not re-seed.

## [0.6.2] - 2026-08-22

### Added
- Mobile top nav (shadcn Sheet). Slim public-page nav falls back to Home / Dashboard if settings omit those links.
- Optional `clicks` on public/studio links. Public taps POST `/api/v1/public/pages/{slug}/links/{id}/click` and ignore a missing route until Go lands counts.

## [0.6.1] - 2026-08-22

### Added
- Shared top nav from public settings (`nav: [{label,href}]`) on `/`, `/about`, `/dashboard`, `/admin`, `/reset`. Slim bar on `/{slug}`.
- Fallback Home / About / Dashboard / Admin until API ships nav.
- Admin settings can edit those nav links.

## [0.6.0] - 2026-08-22

### Changed
- GET /api/v1/me/page returns 200 with an empty PageDTO (nil id, empty slug) instead of 404 when the user has no page. Next treats empty slug as unsaved.
- GET /api/v1/me/page/links returns 200 [] when the user has no page.

### Added
- In-memory seed after Render wipe: admin@synclink.app and gurkanfikretgunak@gmail.com (admins), plus a demo /gurkan page with links so the site has data at process start.
- Home/design copy on GET /api/v1/public/settings and GET/PUT /api/v1/admin/settings: heroTitle, heroSubtitle, heroCta, heroImage, demoSlug.
- Next dashboard studio: 404 or empty /me/page is an unsaved draft. Live preview on the right, tap a card to edit. Inline link edit. Look fields persist.
- Homepage hero reads public settings (title, subtitle, CTA, image, demoSlug) with local fallbacks. Admin can edit those fields.


## [0.5.0] - 2026-08-21

### Added
- Page look on GET/PUT /api/v1/me/page and public GET /api/v1/public/pages/{slug}: avatarShape (circle|rounded|square), accentColor, background (cream|white|dark|motion), motion (none|subtle|lively). Defaults: circle, #111111, cream, subtle.
- Site meta on GET/PUT /api/v1/admin/settings and GET /api/v1/public/settings: metaTitle, metaDescription, ogImage, favicon, themeColor. Next head now reads public settings.
- Dashboard studio: preview cards on the right, tap to open identity / look / links / account on the left.
- Public /{slug} honors avatarShape, accentColor, background, motion. Page-enter and image hover on home, dashboard, public.



## [0.4.0] - 2026-08-21

### Added
- Admin console at `/admin` (stats, users table + mobile sheets, pages, platform settings).
- Public `/about` from `GET /api/v1/public/settings`.
- Login illustration on dashboard gate. Mobile-first tables/sheets.

## [0.3.0] - 2026-08-21

### Added
- Fira Code across the app.
- Dashboard live preview, details (avatar/theme/notes), change password, forgot/reset gates, `/reset` page.

## [0.2.0] - 2026-08-21

### Changed
- Richer UI: visual stations on home, public page, and dashboard. Cream canvas, generated stills in public/stations, shadcn Avatar on /{slug}.
- Dashboard Sign up on `/dashboard` plus `synclink.register` (`POST /api/v1/auth/register`).

## [0.1.0] - 2026-08-21

### Added
- Next.js app at repo root: public page `/{slug}`, editor `/dashboard`, typed client in `src/lib/api.ts`.
- Go API in `backend/`: JWT auth plus page/links REST contract (camelCase).
- Local env via `NEXT_PUBLIC_MF_API_URL=http://localhost:8080`.
- Render blueprint (`render.yaml`) for `synclink-api` in Frankfurt from `backend/Dockerfile`.
- Vercel hosts the Next app from this same repo.

### Changed
- Backend is this repo, not masterfabric-go.

### Added
- Live URLs: app https://synclink-mocha.vercel.app , API https://synclink-api.onrender.com

- Password: PUT /api/v1/me/password, POST /api/v1/auth/forgot-password, POST /api/v1/auth/reset-password.

## [0.4.0] - 2026-08-21

### Added
- Admin API: first registered user is admin. GET/PUT /api/v1/admin/settings, users, pages, stats. Public GET /api/v1/public/settings (about, tagline, signup, maintenance).
