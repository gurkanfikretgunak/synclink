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
	u, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
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
	return s.store.UpdateUser(ctx, u)
}

func (s *Service) ForgotPassword(ctx context.Context, email string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	email = strings.ToLower(strings.TrimSpace(email))
	u, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return "", false
	}
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", false
	}
	token := hex.EncodeToString(raw)
	if err := s.store.PutResetToken(ctx, token, u.ID, time.Now().Add(30*time.Minute)); err != nil {
		return "", false
	}
	return token, true
}

func (s *Service) ResetPassword(ctx context.Context, email, token, next string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	email = strings.ToLower(strings.TrimSpace(email))
	userID, exp, ok, err := s.store.GetResetToken(ctx, token)
	if err != nil || !ok || time.Now().After(exp) {
		return ErrInvalidCreds
	}
	u, err := s.store.GetUserByEmail(ctx, email)
	if err != nil || u.ID != userID {
		return ErrInvalidCreds
	}
	hash, err := s.hashPassword(next)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	if err := s.store.UpdateUser(ctx, u); err != nil {
		return err
	}
	_ = s.store.DeleteResetToken(ctx, token)
	return nil
}
