# SyncLink backend

Go API for SyncLink. Lives in this repo, not masterfabric-go.

```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/synclink?sslmode=disable
export JWT_SECRET=dev-secret
go run ./cmd/server
```

Listens on `:8080`.
