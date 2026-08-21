package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrWeakPassword = errors.New("password must be at least 8 characters")

type resetEntry struct {
	userID uuid.UUID
	exp    time.Time
}

func (s *Service) ensureReset() {
	if s.reset == nil {
		s.reset = map[string]resetEntry{}
	}
}

func (s *Service) hashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", ErrWeakPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, current, next string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.byID[userID]
	if !ok {
		return ErrInvalidCreds
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(current)) != nil {
		return ErrInvalidCreds
	}
	hash, err := s.hashPassword(next)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

func (s *Service) ForgotPassword(ctx context.Context, email string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureReset()
	email = strings.ToLower(strings.TrimSpace(email))
	u, ok := s.users[email]
	if !ok {
		return "", false
	}
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", false
	}
	token := hex.EncodeToString(raw)
	s.reset[token] = resetEntry{userID: u.ID, exp: time.Now().Add(30 * time.Minute)}
	return token, true
}

func (s *Service) ResetPassword(ctx context.Context, email, token, next string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureReset()
	email = strings.ToLower(strings.TrimSpace(email))
	entry, ok := s.reset[token]
	if !ok || time.Now().After(entry.exp) {
		return ErrInvalidCreds
	}
	u, found := s.users[email]
	if !found || u.ID != entry.userID {
		return ErrInvalidCreds
	}
	hash, err := s.hashPassword(next)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	delete(s.reset, token)
	return nil
}
