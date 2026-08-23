package page

import (
	"strings"
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
	ID           uuid.UUID
	UserID       uuid.UUID
	Slug         string
	DisplayName  string
	Bio          string
	AvatarURL    *string
	Theme        string
	AvatarShape  string
	AccentColor  string
	Background   string
	Motion       string
	Socials      []Social
	Verified     bool
	PagePassword *string
	PublishedAt  *time.Time
	CoverURL     *string
	CoverKind    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	Featured      bool
	ThumbnailURL  *string
	StartsAt      *time.Time
	EndsAt        *time.Time
	Sensitive     bool
	Section       string
	EmbedURL      *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Subscriber struct {
	ID        uuid.UUID
	PageID    uuid.UUID
	Email     string
	CreatedAt time.Time
}

type PageDTO struct {
	ID           uuid.UUID `json:"id"`
	Slug         string    `json:"slug"`
	DisplayName  string    `json:"displayName"`
	Bio          string    `json:"bio"`
	AvatarURL    *string   `json:"avatarUrl"`
	Theme        string    `json:"theme"`
	AvatarShape  string    `json:"avatarShape"`
	AccentColor  string    `json:"accentColor"`
	Background   string    `json:"background"`
	Motion       string    `json:"motion"`
	Socials      []Social  `json:"socials"`
	Verified     bool      `json:"verified"`
	PagePassword *string    `json:"pagePassword"`
	PublishedAt  *time.Time `json:"publishedAt"`
	CoverURL     *string    `json:"coverUrl"`
	CoverKind    string     `json:"coverKind"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
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
	Featured      bool       `json:"featured"`
	ThumbnailURL  *string    `json:"thumbnailUrl"`
	StartsAt      *time.Time `json:"startsAt"`
	EndsAt        *time.Time `json:"endsAt"`
	Sensitive     bool       `json:"sensitive"`
	Section       string     `json:"section"`
	EmbedURL      *string    `json:"embedUrl"`
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
	Featured      bool       `json:"featured"`
	ThumbnailURL  *string    `json:"thumbnailUrl"`
	StartsAt      *time.Time `json:"startsAt"`
	EndsAt        *time.Time `json:"endsAt"`
	Sensitive     bool       `json:"sensitive"`
	Section       string     `json:"section"`
	EmbedURL      *string    `json:"embedUrl"`
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
	Verified    bool         `json:"verified"`
	PublishedAt *time.Time   `json:"publishedAt"`
	CoverURL    *string      `json:"coverUrl"`
	CoverKind   string       `json:"coverKind"`
	Links       []PublicLink `json:"links"`
}

type SubscriberDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
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
	Slug         string   `json:"slug"`
	DisplayName  string   `json:"displayName"`
	Bio          string   `json:"bio"`
	AvatarURL    *string  `json:"avatarUrl"`
	Theme        string   `json:"theme"`
	AvatarShape  string   `json:"avatarShape"`
	AccentColor  string   `json:"accentColor"`
	Background   string   `json:"background"`
	Motion       string   `json:"motion"`
	Socials      []Social `json:"socials"`
	PagePassword *string  `json:"pagePassword"`
	CoverURL     *string  `json:"coverUrl"`
	CoverKind    string   `json:"coverKind"`
}

type CreateLinkInput struct {
	Title        string     `json:"title"`
	URL          string     `json:"url"`
	Icon         *string    `json:"icon"`
	Active       *bool      `json:"active"`
	Featured     *bool      `json:"featured"`
	ThumbnailURL *string    `json:"thumbnailUrl"`
	StartsAt     *time.Time `json:"startsAt"`
	EndsAt       *time.Time `json:"endsAt"`
	Sensitive    *bool      `json:"sensitive"`
	Section      string     `json:"section"`
	EmbedURL     *string    `json:"embedUrl"`
}

type UpdateLinkInput struct {
	Title        *string    `json:"title"`
	URL          *string    `json:"url"`
	Icon         *string    `json:"icon"`
	Active       *bool      `json:"active"`
	Featured     *bool      `json:"featured"`
	ThumbnailURL *string    `json:"thumbnailUrl"`
	StartsAt     *time.Time `json:"startsAt"`
	EndsAt       *time.Time `json:"endsAt"`
	Sensitive    *bool      `json:"sensitive"`
	Section      *string    `json:"section"`
	EmbedURL     *string    `json:"embedUrl"`
}

func ToPageDTO(p *Page) PageDTO {
	return PageDTO{
		ID: p.ID, Slug: p.Slug, DisplayName: p.DisplayName, Bio: p.Bio,
		AvatarURL: p.AvatarURL, Theme: p.Theme, AvatarShape: p.AvatarShape,
		AccentColor: p.AccentColor, Background: p.Background, Motion: p.Motion,
		Socials: copySocials(p.Socials), Verified: p.Verified, PagePassword: copyStr(p.PagePassword),
		PublishedAt: copyTime(p.PublishedAt), CoverURL: copyStr(p.CoverURL), CoverKind: p.CoverKind,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func ToLinkDTO(l *Link) LinkDTO {
	return LinkDTO{
		ID: l.ID, Title: l.Title, URL: l.URL, Icon: l.Icon,
		Order: l.Order, Active: l.Active, Clicks: l.Clicks, LastClickedAt: l.LastClickedAt,
		Featured: l.Featured, ThumbnailURL: l.ThumbnailURL, StartsAt: l.StartsAt, EndsAt: l.EndsAt, Sensitive: l.Sensitive,
		Section: l.Section, EmbedURL: l.EmbedURL,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func ToPublicLink(l *Link) PublicLink {
	return PublicLink{
		ID: l.ID, Title: l.Title, URL: l.URL, Icon: l.Icon, Order: l.Order,
		Clicks: l.Clicks, LastClickedAt: l.LastClickedAt,
		Featured: l.Featured, ThumbnailURL: l.ThumbnailURL, StartsAt: l.StartsAt, EndsAt: l.EndsAt, Sensitive: l.Sensitive,
		Section: l.Section, EmbedURL: l.EmbedURL,
	}
}

func ToSubscriberDTO(s *Subscriber) SubscriberDTO {
	return SubscriberDTO{ID: s.ID, Email: s.Email, CreatedAt: s.CreatedAt}
}

func LinkInSchedule(l *Link, now time.Time) bool {
	now = now.UTC()
	if l.StartsAt != nil && now.Before(l.StartsAt.UTC()) {
		return false
	}
	if l.EndsAt != nil && now.After(l.EndsAt.UTC()) {
		return false
	}
	return true
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

func copyStr(in *string) *string {
	if in == nil {
		return nil
	}
	v := *in
	return &v
}

func copyTime(in *time.Time) *time.Time {
	if in == nil {
		return nil
	}
	v := *in
	return &v
}

func normalizeOptionalString(in *string) *string {
	if in == nil {
		return nil
	}
	v := strings.TrimSpace(*in)
	if v == "" {
		return nil
	}
	return &v
}

func normalizeHTTPURLPtr(in *string) *string {
	v := normalizeOptionalString(in)
	if v == nil {
		return nil
	}
	out, ok := normalizeHTTPURL(*v)
	if !ok {
		return nil
	}
	return &out
}

func normalizeSection(s string) string {
	s = strings.TrimSpace(s)
	if len([]rune(s)) > 40 {
		s = string([]rune(s)[:40])
	}
	return s
}

func normalizeCoverKind(k string) string {
	switch strings.ToLower(strings.TrimSpace(k)) {
	case "image", "video":
		return strings.ToLower(strings.TrimSpace(k))
	default:
		return ""
	}
}
