# Changelog

All notable changes to SyncLink live in this file.

## [Unreleased]
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
