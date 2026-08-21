# Changelog

All notable changes to SyncLink live in this file.

## [Unreleased]

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
