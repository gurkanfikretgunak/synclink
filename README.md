# SyncLink

Your links, one page. People design a public page from the dashboard. Admins drive site-wide meta from /admin.

![Home](public/stations/hero.png)

Live app: https://synclink-mocha.vercel.app
API: https://synclink-api.onrender.com

Next.js at repo root (Vercel). Go API in backend/ (Render). Client: src/lib/api.ts.

## What it is

One public URL per person (/{slug}). Circle avatar, accent color, background, and motion live on the page and show on both the editor preview and the live page.

![Orbit](public/stations/orbit.png)

Sign in or sign up on /dashboard. The Go API stores users, pages, links, and settings in SQLite (SYNCLINK_DB, default ./data/synclink.db). Seed admins and /gurkan only when the DB is empty. Later signups stay user.

![Login](public/stations/login.png)

Routes: / public home from settings, /{slug} live page, /about, /dashboard editor, /admin meta, /reset password. Top nav is shared and reads `nav` from public settings (Home / About / Dashboard / Admin by default). On small screens it opens a bottom sheet. Public links send a click ping when that API exists.

## Look and meta

Page GET/PUT /api/v1/me/page and GET /api/v1/public/pages/{slug} now include avatarShape (circle|rounded|square), accentColor, background (cream|white|dark|motion), motion (none|subtle|lively). Defaults: circle, #111111, cream, subtle.

Admin GET/PUT /api/v1/admin/settings and GET /api/v1/public/settings add metaTitle, metaDescription, ogImage, favicon, themeColor plus home copy: heroTitle, heroSubtitle, heroCta, heroCtaHref, heroImage, demoSlug, nav [{label,href}]. GET /api/v1/me/page and /me/page/links return 200 empty when unsaved.

Local: cd backend && go test ./... && go run ./cmd/server
SQLite via SYNCLINK_DB. See CHANGELOG.md and backend/README.md.


Dashboard studio is 404-safe (empty /me/page is a draft). Homepage hero fields come from public settings.

Studio and admin show click totals from GET /api/v1/me/stats and GET /api/v1/admin/stats.

0.9.2: public page shows publishedAt when the API sends it.

0.9.1: public + studio QR share (copy URL / PNG). No new API.

0.9.0 app: public subscribe form, studio Inbox, look presets, link pin/thumbnail/schedule/18+, page password unlock. 0.9.0 API: public email subscribe POST /api/v1/public/pages/{slug}/subscribe; owner GET/DELETE /api/v1/me/subscribers. Links carry featured, thumbnailUrl, startsAt, endsAt, sensitive (public hides inactive/unstarted/ended). Pages have verified (admin PATCH /admin/pages/{id}) and optional pagePassword (public 401 locked unless X-Page-Password). Socials allow whatsapp (https://wa.me/...).
