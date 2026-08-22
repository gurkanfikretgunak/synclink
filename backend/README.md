# SyncLink backend

Go API for the Next app in this repo. Not masterfabric-go.

export JWT_SECRET=dev-secret
export ADDR=:8080
# optional; default ./data/synclink.db (parent dir is created)
export SYNCLINK_DB=./data/synclink.db
go test ./...
go run ./cmd/server

Health: GET http://localhost:8080/health/live
SQLite store (modernc.org/sqlite, no CGO). JWT from /api/v1/auth/register or /login.

Routes: /health/live, /api/v1/auth/register, /api/v1/auth/login, /api/v1/public/pages/{slug}, POST /api/v1/public/pages/{slug}/links/{id}/click (no JWT; 200 {ok:true, clicks:N} or 404), /api/v1/me, /api/v1/me/page, /api/v1/me/page/links, /api/v1/me/page/links/{id}, /api/v1/me/page/links/reorder

Docker: docker build -t synclink-api . && docker run -p 8080:8080 -e JWT_SECRET=dev-secret -e SYNCLINK_DB=/var/data/synclink.db synclink-api
Render: ../render.yaml (rootDir backend, Frankfurt). SYNCLINK_DB=/var/data/synclink.db. Process restart keeps data; a new instance without a disk still starts empty and seeds.

Password: PUT /me/password, POST /auth/forgot-password, POST /auth/reset-password (demo resetToken in forgot response).

Admin: first signup is admin if the users table is empty. Seed (admin@synclink.app, gurkanfikretgunak@gmail.com, /gurkan page) runs only when the DB is empty. Later signups stay user. /admin/* JWT. Public settings: GET /api/v1/public/settings.

GET /me/page returns 200 with an empty PageDTO (empty slug) when unsaved. GET /me/page/links returns 200 [].

Look fields on page: avatarShape, accentColor, background, motion.
Admin/public settings also carry metaTitle, metaDescription, ogImage, favicon, themeColor, heroTitle, heroSubtitle, heroCta, heroCtaHref, heroImage, demoSlug, nav [{label,href}].

Public click: POST /api/v1/public/pages/{slug}/links/{id}/click increments clicks for that active link. Links expose clicks (studio LinkDTO and public links).

