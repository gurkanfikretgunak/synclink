package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gurkanfikretgunak/synclink/backend/internal/auth"
	"github.com/gurkanfikretgunak/synclink/backend/internal/httpapi"
	"github.com/gurkanfikretgunak/synclink/backend/internal/page"
	"github.com/gurkanfikretgunak/synclink/backend/internal/store"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	dbPath := store.PathFromEnv()
	db, err := store.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authSvc := auth.NewServiceWithStore(auth.NewSQLiteStore(db))
	if err := authSvc.SeedIfEmpty(); err != nil {
		log.Fatal(err)
	}
	pages := page.NewService(page.NewSQLiteStore(db))
	if u, ok := authSvc.UserByEmail(auth.SeedOwnerEmail); ok {
		if err := pages.SeedDemoIfEmpty(context.Background(), u.ID); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("synclink api listening on %s (sqlite %s)", addr, dbPath)
	if err := http.ListenAndServe(addr, httpapi.New(authSvc, pages)); err != nil {
		log.Fatal(err)
	}
}
