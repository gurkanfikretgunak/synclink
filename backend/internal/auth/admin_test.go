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
