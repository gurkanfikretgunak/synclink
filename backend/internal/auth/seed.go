package auth

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	SeedAdminEmail = "admin@synclink.app"
	SeedOwnerEmail = "gurkanfikretgunak@gmail.com"
)

func (s *Service) SeedIfEmpty() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ctx := context.Background()
	n, err := s.store.UserCount(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	seeds := []struct {
		email    string
		password string
	}{
		{SeedAdminEmail, "SyncLink-admin-1"},
		{SeedOwnerEmail, "Gurkan123!!"},
	}
	now := time.Now().UTC()
	for _, seed := range seeds {
		hash, err := s.hashPassword(seed.password)
		if err != nil {
			return err
		}
		u := &User{
			ID:           uuid.New(),
			Email:        seed.email,
			PasswordHash: hash,
			Role:         RoleAdmin,
			Status:       StatusActive,
			CreatedAt:    now,
		}
		if err := s.store.CreateUser(ctx, u); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) UserByEmail(email string) (*User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, err := s.store.GetUserByEmail(context.Background(), strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return nil, false
	}
	return u, true
}
