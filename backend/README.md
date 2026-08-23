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

Routes: /health/live, /api/v1/auth/register, /api/v1/auth/login, /api/v1/public/pages/{slug}, POST /api/v1/public/pages/{slug}/links/{id}/click (no JWT; 200 {ok:true, clicks:N} or 404), /api/v1/me, GET /api/v1/me/stats (JWT; 200 {totalClicks, links:[{id,title,clicks,url}]} or empty {totalClicks:0,links:[]}), /api/v1/me/page, /api/v1/me/page/links, /api/v1/me/page/links/{id}, /api/v1/me/page/links/reorder, GET /api/v1/admin/stats (JWT admin; {users,pages,totalClicks})

Docker: docker build -t synclink-api . && docker run -p 8080:8080 -e JWT_SECRET=dev-secret -e SYNCLINK_DB=/var/data/synclink.db synclink-api
Render: ../render.yaml (rootDir backend, Frankfurt). SYNCLINK_DB=/var/data/synclink.db. Process restart keeps data; a new instance without a disk still starts empty and seeds.

Password: PUT /me/password, POST /auth/forgot-password, POST /auth/reset-password (demo resetToken in forgot response).

Admin: first signup is admin if the users table is empty. Seed (admin@synclink.app, gurkanfikretgunak@gmail.com, /gurkan page) runs only when the DB is empty. Later signups stay user. /admin/* JWT. Public settings: GET /api/v1/public/settings.

GET /me/page returns 200 with an empty PageDTO (empty slug) when unsaved. GET /me/page/links returns 200 [].

Look fields on page: avatarShape, accentColor, background, motion.
Admin/public settings also carry metaTitle, metaDescription, ogImage, favicon, themeColor, heroTitle, heroSubtitle, heroCta, heroCtaHref, heroImage, demoSlug, nav [{label,href}].

Public click: POST /api/v1/public/pages/{slug}/links/{id}/click increments clicks for that active link. Links expose clicks (studio LinkDTO and public links).
Links also expose lastClickedAt (RFC3339 or null if never clicked); IncrementClicks sets it in memory and SQLite.

Studio stats: GET /api/v1/me/stats (JWT) returns totalClicks plus each link id/title/clicks/url. Missing page is 200 {totalClicks:0,links:[]}.
Admin stats: GET /api/v1/admin/stats includes users, pages, and totalClicks (SumClicks on memory and SQLite).

Socials on GET/PUT /api/v1/me/page and public GET /api/v1/public/pages/{slug}: `socials: [{ "network": "github", "url": "https://github.com/..." }]`. Empty/null stores and returns `[]`.
Allowed networks (lowercase): instagram, x (twitter normalizes to x), youtube, tiktok, github, linkedin, threads, spotify, website, email.
URLs must be http(s) except email, which also accepts mailto: or an address like name@host. Invalid items are dropped; at most 12 are kept.
SQLite column `pages.socials` TEXT (JSON array); existing DBs get ALTER TABLE on migrate.
