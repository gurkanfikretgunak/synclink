package auth

import (
	"context"
	"testing"
)

func TestChangeAndResetPassword(t *testing.T) {
	s := NewService()
	ctx := context.Background()
	u, _, err := s.Register(ctx, "a@b.com", "oldpass12")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.ChangePassword(ctx, u.ID, "wrong", "newpass12"); err != ErrInvalidCreds {
		t.Fatalf("expected invalid current, got %v", err)
	}
	if err := s.ChangePassword(ctx, u.ID, "oldpass12", "newpass12"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Login(ctx, "a@b.com", "newpass12"); err != nil {
		t.Fatal(err)
	}
	tok, ok := s.ForgotPassword(ctx, "a@b.com")
	if !ok || tok == "" {
		t.Fatal("expected reset token")
	}
	if err := s.ResetPassword(ctx, "a@b.com", tok, "thirdpass"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Login(ctx, "a@b.com", "thirdpass"); err != nil {
		t.Fatal(err)
	}
}
