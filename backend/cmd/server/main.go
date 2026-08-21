package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gurkanfikretgunak/synclink/backend/internal/auth"
	"github.com/gurkanfikretgunak/synclink/backend/internal/httpapi"
	"github.com/gurkanfikretgunak/synclink/backend/internal/page"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	authSvc := auth.NewService()
	if err := authSvc.SeedIfEmpty(); err != nil {
		log.Fatal(err)
	}
	pages := page.NewService(page.NewMemoryStore())
	if u, ok := authSvc.UserByEmail(auth.SeedOwnerEmail); ok {
		if err := pages.SeedDemoIfEmpty(context.Background(), u.ID); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("synclink api listening on %s (in-memory store)", addr)
	if err := http.ListenAndServe(addr, httpapi.New(authSvc, pages)); err != nil {
		log.Fatal(err)
	}
}
