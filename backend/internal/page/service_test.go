package page

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestEmptyGetMyPageAndListLinks(t *testing.T) {
	svc := NewService(NewMemoryStore())
	ctx := context.Background()
	uid := uuid.New()
	got, err := svc.GetMyPage(ctx, uid)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != uuid.Nil || got.Slug != "" || got.DisplayName != "" || got.Bio != "" {
		t.Fatalf("expected empty page, got %#v", got)
	}
	if got.AvatarURL != nil {
		t.Fatalf("avatarUrl should be null, got %#v", got.AvatarURL)
	}
	if got.Theme != ThemeDefault || got.AvatarShape != "circle" || got.AccentColor != "#111111" || got.Background != "cream" || got.Motion != "subtle" {
		t.Fatalf("expected default look, got %#v", got)
	}
	if got.Socials == nil || len(got.Socials) != 0 {
		t.Fatalf("expected empty socials slice, got %#v", got.Socials)
	}
	links, err := svc.ListLinks(ctx, uid)
	if err != nil {
		t.Fatal(err)
	}
	if links == nil || len(links) != 0 {
		t.Fatalf("expected empty slice, got %#v", links)
	}
}

func TestSeedDemoIfEmpty(t *testing.T) {
	svc := NewService(NewMemoryStore())
	ctx := context.Background()
	uid := uuid.New()
	if err := svc.SeedDemoIfEmpty(ctx, uid); err != nil {
		t.Fatal(err)
	}
	pub, err := svc.GetPublicPage(ctx, "gurkan")
	if err != nil {
		t.Fatal(err)
	}
	if pub.DisplayName != "Gürkan" || pub.Theme != ThemeDefault || len(pub.Links) != 3 {
		t.Fatalf("seeded page %#v", pub)
	}
	if pub.Socials == nil || len(pub.Socials) != 2 {
		t.Fatalf("seeded socials %#v", pub.Socials)
	}
	if pub.Socials[0].Network != "github" || pub.Socials[1].Network != "website" {
		t.Fatalf("seeded social networks %#v", pub.Socials)
	}
	if err := svc.SeedDemoIfEmpty(ctx, uuid.New()); err != nil {
		t.Fatal(err)
	}
	pages, err := svc.ListAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 1 {
		t.Fatalf("second seed should no-op, got %d pages", len(pages))
	}
}

func TestRecordClick(t *testing.T) {
	svc := NewService(NewMemoryStore())
	ctx := context.Background()
	u := uuid.New()
	if _, err := svc.UpsertPage(ctx, u, UpsertPageInput{Slug: "gurkan", DisplayName: "G"}); err != nil {
		t.Fatal(err)
	}
	on, err := svc.CreateLink(ctx, u, CreateLinkInput{Title: "On", URL: "https://a.com"})
	if err != nil {
		t.Fatal(err)
	}
	off := false
	hidden, err := svc.CreateLink(ctx, u, CreateLinkInput{Title: "Off", URL: "https://b.com", Active: &off})
	if err != nil {
		t.Fatal(err)
	}
	if on.LastClickedAt != nil {
		t.Fatalf("never-clicked lastClickedAt should be null, got %#v", on.LastClickedAt)
	}
	if hidden.LastClickedAt != nil {
		t.Fatalf("inactive never-clicked lastClickedAt should be null, got %#v", hidden.LastClickedAt)
	}
	n, err := svc.RecordClick(ctx, "gurkan", on.ID)
	if err != nil || n != 1 {
		t.Fatalf("first click n=%d err=%v", n, err)
	}
	n, err = svc.RecordClick(ctx, "gurkan", on.ID)
	if err != nil || n != 2 {
		t.Fatalf("second click n=%d err=%v", n, err)
	}
	if _, err := svc.RecordClick(ctx, "gurkan", hidden.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("inactive should 404, got %v", err)
	}
	if _, err := svc.RecordClick(ctx, "missing", on.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing page should 404, got %v", err)
	}
	if _, err := svc.RecordClick(ctx, "gurkan", uuid.New()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing link should 404, got %v", err)
	}
	pub, err := svc.GetPublicPage(ctx, "gurkan")
	if err != nil {
		t.Fatal(err)
	}
	if len(pub.Links) != 1 || pub.Links[0].Clicks != 2 {
		t.Fatalf("public clicks %#v", pub.Links)
	}
	links, err := svc.ListLinks(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	var got int
	for _, l := range links {
		if l.ID == on.ID {
			got = l.Clicks
		}
	}
	if got != 2 {
		t.Fatalf("studio clicks %d", got)
	}
	var last *time.Time
	for _, l := range links {
		if l.ID == on.ID {
			last = l.LastClickedAt
		}
		if l.ID == hidden.ID && l.LastClickedAt != nil {
			t.Fatalf("never-clicked hidden lastClickedAt %#v", l.LastClickedAt)
		}
	}
	if last == nil {
		t.Fatal("after click lastClickedAt should be set")
	}
	if pub.Links[0].LastClickedAt == nil {
		t.Fatal("public lastClickedAt should be set after click")
	}
}

func TestRecordClickIncrementsAnd404(t *testing.T) {
	svc := NewService(NewMemoryStore())
	ctx := context.Background()
	u := uuid.New()
	if _, err := svc.UpsertPage(ctx, u, UpsertPageInput{Slug: "gurkan", DisplayName: "G"}); err != nil {
		t.Fatal(err)
	}
	on, err := svc.CreateLink(ctx, u, CreateLinkInput{Title: "On", URL: "https://a.com"})
	if err != nil {
		t.Fatal(err)
	}
	off := false
	hidden, err := svc.CreateLink(ctx, u, CreateLinkInput{Title: "Off", URL: "https://b.com", Active: &off})
	if err != nil {
		t.Fatal(err)
	}
	n, err := svc.RecordClick(ctx, "gurkan", on.ID)
	if err != nil || n != 1 {
		t.Fatalf("first click n=%d err=%v", n, err)
	}
	n, err = svc.RecordClick(ctx, "gurkan", on.ID)
	if err != nil || n != 2 {
		t.Fatalf("second click n=%d err=%v", n, err)
	}
	pub, err := svc.GetPublicPage(ctx, "gurkan")
	if err != nil || len(pub.Links) != 1 || pub.Links[0].Clicks != 2 {
		t.Fatalf("public clicks %#v err=%v", pub, err)
	}
	links, err := svc.ListLinks(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	var got int
	for _, l := range links {
		if l.ID == on.ID {
			got = l.Clicks
		}
	}
	if got != 2 {
		t.Fatalf("studio clicks %d", got)
	}
	if _, err := svc.RecordClick(ctx, "gurkan", hidden.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("inactive: %v", err)
	}
	if _, err := svc.RecordClick(ctx, "missing", on.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing page: %v", err)
	}
	if _, err := svc.RecordClick(ctx, "gurkan", uuid.New()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing link: %v", err)
	}
}

func TestMyStatsEmptyAndAfterRecordClick(t *testing.T) {
	svc := NewService(NewMemoryStore())
	ctx := context.Background()
	u := uuid.New()
	empty, err := svc.MyStats(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	if empty.TotalClicks != 0 || empty.Links == nil || len(empty.Links) != 0 {
		t.Fatalf("empty stats %#v", empty)
	}
	if _, err := svc.UpsertPage(ctx, u, UpsertPageInput{Slug: "gurkan", DisplayName: "G"}); err != nil {
		t.Fatal(err)
	}
	afterPage, err := svc.MyStats(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	if afterPage.TotalClicks != 0 || afterPage.Links == nil || len(afterPage.Links) != 0 {
		t.Fatalf("page without links %#v", afterPage)
	}
	link, err := svc.CreateLink(ctx, u, CreateLinkInput{Title: "On", URL: "https://a.com"})
	if err != nil {
		t.Fatal(err)
	}
	n, err := svc.RecordClick(ctx, "gurkan", link.ID)
	if err != nil || n != 1 {
		t.Fatalf("click n=%d err=%v", n, err)
	}
	n, err = svc.RecordClick(ctx, "gurkan", link.ID)
	if err != nil || n != 2 {
		t.Fatalf("click2 n=%d err=%v", n, err)
	}
	stats, err := svc.MyStats(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	if stats.TotalClicks != 2 || len(stats.Links) != 1 {
		t.Fatalf("stats %#v", stats)
	}
	if stats.Links[0].ID != link.ID || stats.Links[0].Title != "On" || stats.Links[0].URL != "https://a.com" || stats.Links[0].Clicks != 2 {
		t.Fatalf("link row %#v", stats.Links[0])
	}
	total, err := svc.SumClicks(ctx)
	if err != nil || total != 2 {
		t.Fatalf("admin sum %d err=%v", total, err)
	}
}

func TestNormalizeAndUpsertSocials(t *testing.T) {
	got := NormalizeSocials([]Social{
		{Network: "Twitter", URL: "https://x.com/g"},
		{Network: "email", URL: "hi@example.com"},
		{Network: "email", URL: "mailto:other@example.com"},
		{Network: "github", URL: "ftp://nope.example"},
		{Network: "myspace", URL: "https://myspace.com/x"},
		{Network: "website", URL: "not-a-url"},
		{Network: "instagram", URL: "https://instagram.com/g"},
	})
	if len(got) != 4 {
		t.Fatalf("expected 4 kept, got %#v", got)
	}
	if got[0].Network != "x" || got[0].URL != "https://x.com/g" {
		t.Fatalf("twitter→x %#v", got[0])
	}
	if got[1].Network != "email" || got[1].URL != "hi@example.com" {
		t.Fatalf("plain email %#v", got[1])
	}
	if got[2].URL != "mailto:other@example.com" {
		t.Fatalf("mailto %#v", got[2])
	}
	if got[3].Network != "instagram" {
		t.Fatalf("instagram %#v", got[3])
	}
	if n := NormalizeSocials(nil); n == nil || len(n) != 0 {
		t.Fatalf("nil → [] got %#v", n)
	}
	tooMany := make([]Social, 20)
	for i := range tooMany {
		tooMany[i] = Social{Network: "website", URL: "https://example.com/" + string(rune('a'+i%26))}
	}
	if len(NormalizeSocials(tooMany)) != 12 {
		t.Fatalf("max 12, got %d", len(NormalizeSocials(tooMany)))
	}

	svc := NewService(NewMemoryStore())
	ctx := context.Background()
	u := uuid.New()
	page, err := svc.UpsertPage(ctx, u, UpsertPageInput{
		Slug: "me", DisplayName: "Me",
		Socials: []Social{
			{Network: "twitter", URL: "https://x.com/g"},
			{Network: "nope", URL: "https://nope.com"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if page.Socials == nil || len(page.Socials) != 1 || page.Socials[0].Network != "x" {
		t.Fatalf("upsert socials %#v", page.Socials)
	}
	mine, err := svc.GetMyPage(ctx, u)
	if err != nil || len(mine.Socials) != 1 || mine.Socials[0].Network != "x" {
		t.Fatalf("get my page socials %#v err=%v", mine, err)
	}
	pub, err := svc.GetPublicPage(ctx, "me")
	if err != nil || len(pub.Socials) != 1 || pub.Socials[0].Network != "x" {
		t.Fatalf("public socials %#v err=%v", pub, err)
	}
	cleared, err := svc.UpsertPage(ctx, u, UpsertPageInput{Slug: "me", DisplayName: "Me"})
	if err != nil || cleared.Socials == nil || len(cleared.Socials) != 0 {
		t.Fatalf("clear socials %#v err=%v", cleared, err)
	}
}
