# Changelog

All notable changes to SyncLink live in this file.

## [Unreleased]

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
