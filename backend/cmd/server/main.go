package main

import (
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
	pages := page.NewService(page.NewMemoryStore())
	log.Printf("synclink api listening on %s (in-memory store)", addr)
	if err := http.ListenAndServe(addr, httpapi.New(authSvc, pages)); err != nil {
		log.Fatal(err)
	}
}
