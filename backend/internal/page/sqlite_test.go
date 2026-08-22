package page

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/gurkanfikretgunak/synclink/backend/internal/store"
)

func TestSQLitePersistAndClick(t *testing.T) {
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	svc := NewService(NewSQLiteStore(db))
	ctx := context.Background()
	u := uuid.New()
	if _, err := svc.UpsertPage(ctx, u, UpsertPageInput{Slug: "gurkan", DisplayName: "G"}); err != nil {
		t.Fatal(err)
	}
	link, err := svc.CreateLink(ctx, u, CreateLinkInput{Title: "On", URL: "https://a.com"})
	if err != nil {
		t.Fatal(err)
	}
	n, err := svc.RecordClick(ctx, "gurkan", link.ID)
	if err != nil || n != 1 {
		t.Fatalf("click n=%d err=%v", n, err)
	}
	got, err := svc.GetPublicPage(ctx, "gurkan")
	if err != nil || len(got.Links) != 1 || got.Links[0].Clicks != 1 {
		t.Fatalf("persist public %#v err=%v", got, err)
	}
	if _, err := svc.RecordClick(ctx, "nope", link.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing page %v", err)
	}
}
