package page

import (
	"time"

	"github.com/google/uuid"
)

const (
	ThemeDefault  = "default"
	ThemeDark     = "dark"
	ThemeLight    = "light"
	ThemeColorful = "colorful"
)

func ValidTheme(theme string) bool {
	switch theme {
	case ThemeDefault, ThemeDark, ThemeLight, ThemeColorful:
		return true
	default:
		return false
	}
}

type Page struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Slug        string
	DisplayName string
	Bio         string
	AvatarURL   *string
	Theme       string
	AvatarShape string
	AccentColor string
	Background  string
	Motion      string
	Socials     []Social
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Link struct {
	ID            uuid.UUID
	PageID        uuid.UUID
	Title         string
	URL           string
	Icon          *string
	Order         int
	Active        bool
	Clicks        int
	LastClickedAt *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PageDTO struct {
	ID          uuid.UUID `json:"id"`
	Slug        string    `json:"slug"`
	DisplayName string    `json:"displayName"`
	Bio         string    `json:"bio"`
	AvatarURL   *string   `json:"avatarUrl"`
	Theme       string    `json:"theme"`
	AvatarShape string    `json:"avatarShape"`
	AccentColor string    `json:"accentColor"`
	Background  string    `json:"background"`
	Motion      string    `json:"motion"`
	Socials     []Social  `json:"socials"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LinkDTO struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	URL           string     `json:"url"`
	Icon          *string    `json:"icon"`
	Order         int        `json:"order"`
	Active        bool       `json:"active"`
	Clicks        int        `json:"clicks"`
	LastClickedAt *time.Time `json:"lastClickedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type PublicLink struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	URL           string     `json:"url"`
	Icon          *string    `json:"icon"`
	Order         int        `json:"order"`
	Clicks        int        `json:"clicks"`
	LastClickedAt *time.Time `json:"lastClickedAt"`
}

type PublicPage struct {
	Slug        string       `json:"slug"`
	DisplayName string       `json:"displayName"`
	Bio         string       `json:"bio"`
	AvatarURL   *string      `json:"avatarUrl"`
	Theme       string       `json:"theme"`
	AvatarShape string       `json:"avatarShape"`
	AccentColor string       `json:"accentColor"`
	Background  string       `json:"background"`
	Motion      string       `json:"motion"`
	Socials     []Social     `json:"socials"`
	Links       []PublicLink `json:"links"`
}

type MyStatsLink struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Clicks int       `json:"clicks"`
	URL    string    `json:"url"`
}

type MyStats struct {
	TotalClicks int           `json:"totalClicks"`
	Links       []MyStatsLink `json:"links"`
}

type UpsertPageInput struct {
	Slug        string   `json:"slug"`
	DisplayName string   `json:"displayName"`
	Bio         string   `json:"bio"`
	AvatarURL   *string  `json:"avatarUrl"`
	Theme       string   `json:"theme"`
	AvatarShape string   `json:"avatarShape"`
	AccentColor string   `json:"accentColor"`
	Background  string   `json:"background"`
	Motion      string   `json:"motion"`
	Socials     []Social `json:"socials"`
}

type CreateLinkInput struct {
	Title  string  `json:"title"`
	URL    string  `json:"url"`
	Icon   *string `json:"icon"`
	Active *bool   `json:"active"`
}

type UpdateLinkInput struct {
	Title  *string `json:"title"`
	URL    *string `json:"url"`
	Icon   *string `json:"icon"`
	Active *bool   `json:"active"`
}

func ToPageDTO(p *Page) PageDTO {
	return PageDTO{
		ID: p.ID, Slug: p.Slug, DisplayName: p.DisplayName, Bio: p.Bio,
		AvatarURL: p.AvatarURL, Theme: p.Theme, AvatarShape: p.AvatarShape,
		AccentColor: p.AccentColor, Background: p.Background, Motion: p.Motion,
		Socials:   copySocials(p.Socials),
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func ToLinkDTO(l *Link) LinkDTO {
	return LinkDTO{
		ID: l.ID, Title: l.Title, URL: l.URL, Icon: l.Icon,
		Order: l.Order, Active: l.Active, Clicks: l.Clicks, LastClickedAt: l.LastClickedAt,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func NormalizeLook(shape, accent, bg, motion string) (string, string, string, string) {
	switch shape {
	case "circle", "rounded", "square":
	default:
		shape = "circle"
	}
	if accent == "" {
		accent = "#111111"
	}
	switch bg {
	case "cream", "white", "dark", "motion":
	default:
		bg = "cream"
	}
	switch motion {
	case "none", "subtle", "lively":
	default:
		motion = "subtle"
	}
	return shape, accent, bg, motion
}
