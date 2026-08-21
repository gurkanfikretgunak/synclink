package auth

import (
	"context"
	"testing"
)

func TestFirstUserIsAdminAndSettings(t *testing.T) {
	s := NewService()
	ctx := context.Background()
	a, _, err := s.Register(ctx, "admin@s.com", "password1")
	if err != nil {
		t.Fatal(err)
	}
	if a.Role != RoleAdmin {
		t.Fatalf("first user should be admin, got %s", a.Role)
	}
	u, _, err := s.Register(ctx, "user@s.com", "password1")
	if err != nil {
		t.Fatal(err)
	}
	if u.Role != RoleUser {
		t.Fatalf("second user should be user, got %s", u.Role)
	}
	if !s.IsAdmin(a.ID) || s.IsAdmin(u.ID) {
		t.Fatal("admin flags wrong")
	}
	st := StatusDisabled
	if _, err := s.UpdateUser(ctx, u.ID, nil, &st); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Login(ctx, "user@s.com", "password1"); err != ErrDisabled {
		t.Fatalf("expected disabled, got %v", err)
	}
	got := s.UpdateSettings(Settings{SiteName: "X", About: "Y", SignupEnabled: true})
	if got.SiteName != "X" || got.About != "Y" {
		t.Fatalf("settings %#v", got)
	}
}

func TestPublicSettingsHeroKeys(t *testing.T) {
	s := NewService()
	pub := s.PublicSettings()
	if pub["heroTitle"] != "One page. Every link." {
		t.Fatalf("heroTitle %v", pub["heroTitle"])
	}
	if pub["heroSubtitle"] != "A quieter public page. White space, type, and a few stills. Edit from the dashboard." {
		t.Fatalf("heroSubtitle %v", pub["heroSubtitle"])
	}
	if pub["heroCta"] != "Create your page" {
		t.Fatalf("heroCta %v", pub["heroCta"])
	}
	if pub["heroImage"] != "/stations/hero.png" {
		t.Fatalf("heroImage %v", pub["heroImage"])
	}
	if pub["demoSlug"] != "gurkan" {
		t.Fatalf("demoSlug %v", pub["demoSlug"])
	}
	if pub["heroCtaHref"] != "/dashboard" {
		t.Fatalf("heroCtaHref %v", pub["heroCtaHref"])
	}
	nav, ok := pub["nav"].([]NavItem)
	if !ok || len(nav) != 3 {
		t.Fatalf("nav %v", pub["nav"])
	}
	if nav[0] != (NavItem{Label: "About", Href: "/about"}) ||
		nav[1] != (NavItem{Label: "Admin", Href: "/admin"}) ||
		nav[2] != (NavItem{Label: "Dashboard", Href: "/dashboard"}) {
		t.Fatalf("default nav %#v", nav)
	}
	got := s.UpdateSettings(Settings{SiteName: "X", HeroTitle: "New hero", HeroCtaHref: "/about", SignupEnabled: true})
	if got.HeroTitle != "New hero" || got.HeroCta != "Create your page" || got.HeroCtaHref != "/about" {
		t.Fatalf("update %#v", got)
	}
	got = s.UpdateSettings(Settings{SiteName: "X", SignupEnabled: true})
	if got.HeroTitle != "New hero" {
		t.Fatalf("empty heroTitle should keep previous, got %#v", got)
	}
	if got.HeroCtaHref != "/about" {
		t.Fatalf("empty heroCtaHref should keep previous, got %#v", got)
	}
	got = s.UpdateSettings(Settings{
		SiteName:      "X",
		SignupEnabled: true,
		Nav:           []NavItem{{Label: " Home ", Href: " / "}},
	})
	if len(got.Nav) != 1 || got.Nav[0].Label != "Home" || got.Nav[0].Href != "/" {
		t.Fatalf("nav replace %#v", got.Nav)
	}
	got = s.UpdateSettings(Settings{SiteName: "X", SignupEnabled: true})
	if len(got.Nav) != 1 || got.Nav[0] != (NavItem{Label: "Home", Href: "/"}) {
		t.Fatalf("empty nav should keep previous, got %#v", got.Nav)
	}
	got = s.UpdateSettings(Settings{
		SiteName:      "X",
		SignupEnabled: true,
		Nav:           []NavItem{{Label: "  ", Href: "  "}},
	})
	if len(got.Nav) != 1 || got.Nav[0] != (NavItem{Label: "Home", Href: "/"}) {
		t.Fatalf("trimmed-empty nav should keep previous, got %#v", got.Nav)
	}
	pub = s.PublicSettings()
	if pub["heroCtaHref"] != "/about" {
		t.Fatalf("public heroCtaHref after update %v", pub["heroCtaHref"])
	}
	nav, ok = pub["nav"].([]NavItem)
	if !ok || len(nav) != 1 || nav[0] != (NavItem{Label: "Home", Href: "/"}) {
		t.Fatalf("public nav after update %v", pub["nav"])
	}
}

func TestSeedIfEmpty(t *testing.T) {
	s := NewService()
	ctx := context.Background()
	if err := s.SeedIfEmpty(); err != nil {
		t.Fatal(err)
	}
	admin, ok := s.UserByEmail(SeedAdminEmail)
	if !ok || admin.Role != RoleAdmin || admin.Status != StatusActive {
		t.Fatalf("admin seed %#v ok=%v", admin, ok)
	}
	owner, ok := s.UserByEmail(SeedOwnerEmail)
	if !ok || owner.Role != RoleAdmin || owner.Status != StatusActive {
		t.Fatalf("owner seed %#v ok=%v", owner, ok)
	}
	if _, _, err := s.Login(ctx, SeedAdminEmail, "SyncLink-admin-1"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Login(ctx, SeedOwnerEmail, "Gurkan123!!"); err != nil {
		t.Fatal(err)
	}
	if err := s.SeedIfEmpty(); err != nil {
		t.Fatal(err)
	}
	if len(s.ListUsers(ctx)) != 2 {
		t.Fatalf("second seed should no-op, got %d users", len(s.ListUsers(ctx)))
	}
	u, _, err := s.Register(ctx, "later@s.com", "password1")
	if err != nil {
		t.Fatal(err)
	}
	if u.Role != RoleUser {
		t.Fatalf("signup after seed should stay user, got %s", u.Role)
	}
}
