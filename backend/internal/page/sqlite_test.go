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
	if link.LastClickedAt != nil {
		t.Fatalf("never-clicked lastClickedAt should be null, got %#v", link.LastClickedAt)
	}
	n, err := svc.RecordClick(ctx, "gurkan", link.ID)
	if err != nil || n != 1 {
		t.Fatalf("click n=%d err=%v", n, err)
	}
	got, err := svc.GetPublicPage(ctx, "gurkan")
	if err != nil || len(got.Links) != 1 || got.Links[0].Clicks != 1 {
		t.Fatalf("persist public %#v err=%v", got, err)
	}
	if got.Links[0].LastClickedAt == nil {
		t.Fatal("after click lastClickedAt should be set")
	}
	if _, err := svc.RecordClick(ctx, "nope", link.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing page %v", err)
	}
}

func TestSQLiteMyStatsAndSumClicks(t *testing.T) {
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	svc := NewService(NewSQLiteStore(db))
	ctx := context.Background()
	u := uuid.New()
	empty, err := svc.MyStats(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	if empty.TotalClicks != 0 || empty.Links == nil || len(empty.Links) != 0 {
		t.Fatalf("empty %#v", empty)
	}
	if _, err := svc.UpsertPage(ctx, u, UpsertPageInput{Slug: "gurkan", DisplayName: "G"}); err != nil {
		t.Fatal(err)
	}
	afterPage, err := svc.MyStats(ctx, u)
	if err != nil || afterPage.TotalClicks != 0 || afterPage.Links == nil || len(afterPage.Links) != 0 {
		t.Fatalf("page without links %#v err=%v", afterPage, err)
	}
	link, err := svc.CreateLink(ctx, u, CreateLinkInput{Title: "On", URL: "https://a.com"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordClick(ctx, "gurkan", link.ID); err != nil {
		t.Fatal(err)
	}
	stats, err := svc.MyStats(ctx, u)
	if err != nil || stats.TotalClicks != 1 || len(stats.Links) != 1 || stats.Links[0].Clicks != 1 {
		t.Fatalf("sqlite mystats %#v err=%v", stats, err)
	}
	total, err := svc.SumClicks(ctx)
	if err != nil || total != 1 {
		t.Fatalf("sqlite sum %d err=%v", total, err)
	}
}
