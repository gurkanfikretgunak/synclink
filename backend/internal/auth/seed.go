package auth

import (
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
	if len(s.users) > 0 {
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
		s.users[u.Email] = u
		s.byID[u.ID] = u
	}
	return nil
}

func (s *Service) UserByEmail(email string) (*User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[strings.ToLower(strings.TrimSpace(email))]
	return u, ok
}
