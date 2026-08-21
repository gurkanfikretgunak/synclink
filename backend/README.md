# SyncLink backend

Go API for the Next app in this repo. Not masterfabric-go.

export JWT_SECRET=dev-secret
export ADDR=:8080
go test ./...
go run ./cmd/server

Health: GET http://localhost:8080/health/live
In-memory store. JWT from /api/v1/auth/register or /login.

Routes: /health/live, /api/v1/auth/register, /api/v1/auth/login, /api/v1/public/pages/{slug}, /api/v1/me, /api/v1/me/page, /api/v1/me/page/links, /api/v1/me/page/links/{id}, /api/v1/me/page/links/reorder

Docker: docker build -t synclink-api . && docker run -p 8080:8080 -e JWT_SECRET=dev-secret synclink-api
Render: ../render.yaml (rootDir backend, Frankfurt).

Password: PUT /me/password, POST /auth/forgot-password, POST /auth/reset-password (demo resetToken in forgot response).

Admin: first signup is admin if the store is empty. Process start seeds admin@synclink.app and gurkanfikretgunak@gmail.com plus a /gurkan demo page so a Render wipe still has data. Later signups stay user. /admin/* JWT. Public settings: GET /api/v1/public/settings.

GET /me/page returns 200 with an empty PageDTO (empty slug) when unsaved. GET /me/page/links returns 200 [].

Look fields on page: avatarShape, accentColor, background, motion.
Admin/public settings also carry metaTitle, metaDescription, ogImage, favicon, themeColor, heroTitle, heroSubtitle, heroCta, heroImage, demoSlug.
