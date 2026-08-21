# Changelog

All notable changes to SyncLink live in this file.

## [Unreleased]
## [0.6.1] - 2026-08-22

### Added
- Shared top nav from public settings (`nav: [{label,href}]`) on `/`, `/about`, `/dashboard`, `/admin`, `/reset`. Slim bar on `/{slug}`. Fallback Home / About / Dashboard / Admin until API ships nav.
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
