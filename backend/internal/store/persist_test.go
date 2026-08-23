package store_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/gurkanfikretgunak/synclink/backend/internal/auth"
	"github.com/gurkanfikretgunak/synclink/backend/internal/page"
	"github.com/gurkanfikretgunak/synclink/backend/internal/store"
)

func TestSQLiteSurvivesReopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "synclink.db")
	ctx := context.Background()

	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	authSvc := auth.NewServiceWithStore(auth.NewSQLiteStore(db))
	pages := page.NewService(page.NewSQLiteStore(db))
	if err := authSvc.SeedIfEmpty(); err != nil {
		t.Fatal(err)
	}
	if err := authSvc.SeedIfEmpty(); err != nil {
		t.Fatal(err)
	}
	if n := len(authSvc.ListUsers(ctx)); n != 2 {
		t.Fatalf("seed users %d", n)
	}
	owner, ok := authSvc.UserByEmail(auth.SeedOwnerEmail)
	if !ok {
		t.Fatal("missing owner")
	}
	if err := pages.SeedDemoIfEmpty(ctx, owner.ID); err != nil {
		t.Fatal(err)
	}
	uid := uuid.New()
	if _, _, err := authSvc.Register(ctx, "later@s.com", "password1"); err != nil {
		t.Fatal(err)
	}
	if _, err := pages.UpsertPage(ctx, uid, page.UpsertPageInput{Slug: "extra", DisplayName: "Extra"}); err != nil {
		t.Fatal(err)
	}
	got := authSvc.UpdateSettings(auth.Settings{SiteName: "Persisted", SignupEnabled: true})
	if got.SiteName != "Persisted" {
		t.Fatalf("settings %#v", got)
	}
	_ = db.Close()

	db2, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db2.Close()
	auth2 := auth.NewServiceWithStore(auth.NewSQLiteStore(db2))
	pages2 := page.NewService(page.NewSQLiteStore(db2))
	if err := auth2.SeedIfEmpty(); err != nil {
		t.Fatal(err)
	}
	if n := len(auth2.ListUsers(ctx)); n != 3 {
		t.Fatalf("users after reopen %d", n)
	}
	if _, _, err := auth2.Login(ctx, auth.SeedOwnerEmail, "Gurkan123!!"); err != nil {
		t.Fatal(err)
	}
	if err := pages2.SeedDemoIfEmpty(ctx, uuid.New()); err != nil {
		t.Fatal(err)
	}
	all, err := pages2.ListAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Fatalf("pages after reopen %d", len(all))
	}
	pub, err := pages2.GetPublicPage(ctx, "gurkan")
	if err != nil || pub.DisplayName != "Gürkan" || len(pub.Links) != 3 {
		t.Fatalf("public %#v err=%v", pub, err)
	}
	if auth2.Settings().SiteName != "Persisted" {
		t.Fatalf("settings after reopen %#v", auth2.Settings())
	}
	n, err := pages2.RecordClick(ctx, "gurkan", pub.Links[0].ID)
	if err != nil || n != 1 {
		t.Fatalf("click n=%d err=%v", n, err)
	}
	if err := db2.Close(); err != nil {
		t.Fatal(err)
	}
	db3, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db3.Close()
	pages3 := page.NewService(page.NewSQLiteStore(db3))
	auth3 := auth.NewServiceWithStore(auth.NewSQLiteStore(db3))
	pub3, err := pages3.GetPublicPage(ctx, "gurkan")
	if err != nil || pub3.Links[0].Clicks != 1 {
		t.Fatalf("clicks after reopen %#v err=%v", pub3, err)
	}
	if pub3.Links[0].LastClickedAt == nil {
		t.Fatalf("lastClickedAt after reopen %#v", pub3)
	}
	if auth3.Settings().SiteName != "Persisted" {
		t.Fatalf("settings after reopen %#v", auth3.Settings())
	}
}

func TestPageSQLiteMatchesMemoryContract(t *testing.T) {
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := page.NewService(page.NewSQLiteStore(db))
	ctx := context.Background()
	uid := uuid.New()
	got, err := svc.GetMyPage(ctx, uid)
	if err != nil || got.Slug != "" || got.ID != uuid.Nil {
		t.Fatalf("empty page %#v err=%v", got, err)
	}
	links, err := svc.ListLinks(ctx, uid)
	if err != nil || links == nil || len(links) != 0 {
		t.Fatalf("empty links %#v err=%v", links, err)
	}
}
