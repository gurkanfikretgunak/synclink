package auth

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrExists       = errors.New("already exists")
	ErrValidation   = errors.New("validation")
	ErrDisabled     = errors.New("account disabled")
	ErrSignupClosed = errors.New("signup closed")
	ErrNotFound     = errors.New("not found")
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Service struct {
	mu     sync.Mutex
	store  Store
	secret []byte
}

func NewService() *Service {
	return NewServiceWithStore(NewMemoryStore())
}

func NewServiceWithStore(store Store) *Service {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	return &Service{
		store:  store,
		secret: []byte(secret),
	}
}

type Claims struct {
	UserID uuid.UUID `json:"uid"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func (s *Service) Register(ctx context.Context, email, password string) (*User, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, err := s.store.UserCount(ctx)
	if err != nil {
		return nil, "", err
	}
	st, err := s.store.GetSettings(ctx)
	if err != nil {
		return nil, "", err
	}
	if !st.SignupEnabled && n > 0 {
		return nil, "", ErrSignupClosed
	}
	if _, err := s.store.GetUserByEmail(ctx, email); err == nil {
		return nil, "", ErrExists
	} else if !errors.Is(err, ErrNotFound) {
		return nil, "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	role := RoleUser
	if n == 0 {
		role = RoleAdmin
	}
	u := &User{ID: uuid.New(), Email: email, PasswordHash: string(hash), Role: role, Status: StatusActive, CreatedAt: time.Now().UTC()}
	if err := s.store.CreateUser(ctx, u); err != nil {
		return nil, "", err
	}
	tok, err := s.token(u)
	return u, tok, err
}

func (s *Service) Login(ctx context.Context, email, password string) (*User, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCreds
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return nil, "", ErrInvalidCreds
	}
	if u.Status == StatusDisabled {
		return nil, "", ErrDisabled
	}
	tok, err := s.token(u)
	return u, tok, err
}

func (s *Service) User(id uuid.UUID) (*User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, err := s.store.GetUserByID(context.Background(), id)
	if err != nil {
		return nil, false
	}
	return u, true
}

func (s *Service) token(u *User) (string, error) {
	claims := Claims{
		UserID: u.ID,
		Email:  u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *Service) Parse(token string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
