# SyncLink

Text-first pages. Next.js at root (Vercel) + Go API in backend/ (Render).
Client: src/lib/api.ts. Public: /{slug}. Editor: /dashboard.

Local API: cd backend && go run ./cmd/server (http://localhost:8080/health/live)
See CHANGELOG.md and backend/README.md.

## Live

- App: https://synclink-mocha.vercel.app
- API: https://synclink-api.onrender.com
Dashboard: sign in or Sign up (email + password, min 8).
v0.2.0 visual stations on /, /{slug}, /dashboard.
v0.3.0 Fira Code, live preview, password flows.
