package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Store interface {
	UserCount(ctx context.Context) (int, error)
	CreateUser(ctx context.Context, u *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetSettings(ctx context.Context) (Settings, error)
	PutSettings(ctx context.Context, st Settings) error
	PutResetToken(ctx context.Context, token string, userID uuid.UUID, exp time.Time) error
	GetResetToken(ctx context.Context, token string) (uuid.UUID, time.Time, bool, error)
	DeleteResetToken(ctx context.Context, token string) error
}

func DefaultSettings() Settings {
	return Settings{
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
	}
}
