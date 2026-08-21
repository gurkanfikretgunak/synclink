package page

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestPublicPageOmitsInactiveAndSlugConflict(t *testing.T) {
	svc := NewService(NewMemoryStore())
	ctx := context.Background()
	u1, u2 := uuid.New(), uuid.New()
	if _, err := svc.UpsertPage(ctx, u1, UpsertPageInput{Slug: "gurkan", DisplayName: "G"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateLink(ctx, u1, CreateLinkInput{Title: "On", URL: "https://a.com"}); err != nil {
		t.Fatal(err)
	}
	off := false
	if _, err := svc.CreateLink(ctx, u1, CreateLinkInput{Title: "Off", URL: "https://b.com", Active: &off}); err != nil {
		t.Fatal(err)
	}
	pub, err := svc.GetPublicPage(ctx, "gurkan")
	if err != nil {
		t.Fatal(err)
	}
	if len(pub.Links) != 1 || pub.Links[0].Title != "On" {
		t.Fatalf("expected one active link, got %#v", pub.Links)
	}
	if _, err := svc.GetPublicPage(ctx, "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
	if _, err := svc.UpsertPage(ctx, u2, UpsertPageInput{Slug: "gurkan", DisplayName: "X"}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestLinkReorder(t *testing.T) {
	svc := NewService(NewMemoryStore())
	ctx := context.Background()
	u := uuid.New()
	if _, err := svc.UpsertPage(ctx, u, UpsertPageInput{Slug: "me", DisplayName: "Me"}); err != nil {
		t.Fatal(err)
	}
	a, _ := svc.CreateLink(ctx, u, CreateLinkInput{Title: "A", URL: "https://a.com"})
	b, _ := svc.CreateLink(ctx, u, CreateLinkInput{Title: "B", URL: "https://b.com"})
	out, err := svc.ReorderLinks(ctx, u, []uuid.UUID{b.ID, a.ID})
	if err != nil {
		t.Fatal(err)
	}
	if out[0].ID != b.ID || out[0].Order != 0 {
		t.Fatalf("reorder failed: %#v", out)
	}
}
