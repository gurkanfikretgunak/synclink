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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Link struct {
	ID        uuid.UUID
	PageID    uuid.UUID
	Title     string
	URL       string
	Icon      *string
	Order     int
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PageDTO struct {
	ID          uuid.UUID `json:"id"`
	Slug        string    `json:"slug"`
	DisplayName string    `json:"displayName"`
	Bio         string    `json:"bio"`
	AvatarURL   *string   `json:"avatarUrl"`
	Theme       string    `json:"theme"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LinkDTO struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Icon      *string   `json:"icon"`
	Order     int       `json:"order"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PublicLink struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	URL   string    `json:"url"`
	Icon  *string   `json:"icon"`
	Order int       `json:"order"`
}

type PublicPage struct {
	Slug        string       `json:"slug"`
	DisplayName string       `json:"displayName"`
	Bio         string       `json:"bio"`
	AvatarURL   *string      `json:"avatarUrl"`
	Theme       string       `json:"theme"`
	Links       []PublicLink `json:"links"`
}

type UpsertPageInput struct {
	Slug        string  `json:"slug"`
	DisplayName string  `json:"displayName"`
	Bio         string  `json:"bio"`
	AvatarURL   *string `json:"avatarUrl"`
	Theme       string  `json:"theme"`
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
		AvatarURL: p.AvatarURL, Theme: p.Theme, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func ToLinkDTO(l *Link) LinkDTO {
	return LinkDTO{
		ID: l.ID, Title: l.Title, URL: l.URL, Icon: l.Icon,
		Order: l.Order, Active: l.Active, CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}
