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
	mu       sync.Mutex
	users    map[string]*User
	byID     map[uuid.UUID]*User
	reset    map[string]resetEntry
	secret   []byte
	settings Settings
}

func NewService() *Service {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	return &Service{
		users:  map[string]*User{},
		byID:   map[uuid.UUID]*User{},
		reset:  map[string]resetEntry{},
		secret: []byte(secret),
		settings: Settings{
			SiteName:        "SyncLink",
			Tagline:         "Your links, one page.",
			About:           "SyncLink is a text-first page for people and brands.",
			SupportEmail:    "hello@synclink.app",
			SignupEnabled:   true,
			Maintenance:     false,
			MetaTitle:       "SyncLink",
			MetaDescription: "Your links, one page.",
			ThemeColor:      "#111111",
			HeroTitle:       "One page. Every link.",
			HeroSubtitle:    "A quieter public page. White space, type, and a few stills. Edit from the dashboard.",
			HeroCta:         "Create your page",
			HeroCtaHref:     "/dashboard",
			HeroImage:       "/stations/hero.png",
			DemoSlug:        "gurkan",
			Nav:             defaultNav(),
		},
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
	if !s.settings.SignupEnabled && len(s.users) > 0 {
		return nil, "", ErrSignupClosed
	}
	if _, ok := s.users[email]; ok {
		return nil, "", ErrExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	role := RoleUser
	if len(s.users) == 0 {
		role = RoleAdmin
	}
	u := &User{ID: uuid.New(), Email: email, PasswordHash: string(hash), Role: role, Status: StatusActive, CreatedAt: time.Now().UTC()}
	s.users[email] = u
	s.byID[u.ID] = u
	tok, err := s.token(u)
	return u, tok, err
}

func (s *Service) Login(ctx context.Context, email, password string) (*User, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[email]
	if !ok {
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
	u, ok := s.byID[id]
	return u, ok
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
