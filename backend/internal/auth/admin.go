package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	RoleUser       = "user"
	RoleAdmin      = "admin"
	StatusActive   = "active"
	StatusDisabled = "disabled"
)

var ErrForbidden = errors.New("forbidden")

type UserInfo struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *User) Info() UserInfo {
	return UserInfo{ID: u.ID, Email: u.Email, Role: u.Role, Status: u.Status, CreatedAt: u.CreatedAt}
}

func (s *Service) ListUsers(ctx context.Context) []UserInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, err := s.store.ListUsers(ctx)
	if err != nil {
		return []UserInfo{}
	}
	out := make([]UserInfo, 0, len(users))
	for _, u := range users {
		out = append(out, u.Info())
	}
	return out
}

func (s *Service) IsAdmin(id uuid.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, err := s.store.GetUserByID(context.Background(), id)
	return err == nil && u.Role == RoleAdmin && u.Status == StatusActive
}

func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, role, status *string) (*UserInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return nil, ErrInvalidCreds
	}
	if role != nil {
		r := strings.TrimSpace(*role)
		if r != RoleUser && r != RoleAdmin {
			return nil, ErrValidation
		}
		u.Role = r
	}
	if status != nil {
		st := strings.TrimSpace(*status)
		if st != StatusActive && st != StatusDisabled {
			return nil, ErrValidation
		}
		u.Status = st
	}
	if err := s.store.UpdateUser(ctx, u); err != nil {
		return nil, err
	}
	info := u.Info()
	return &info, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.store.DeleteUser(ctx, id); err != nil {
		return ErrInvalidCreds
	}
	return nil
}

type NavItem struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

func defaultNav() []NavItem {
	return []NavItem{
		{Label: "About", Href: "/about"},
		{Label: "Admin", Href: "/admin"},
		{Label: "Dashboard", Href: "/dashboard"},
	}
}

func trimNav(in []NavItem) []NavItem {
	if len(in) == 0 {
		return nil
	}
	out := make([]NavItem, 0, len(in))
	for _, item := range in {
		label := strings.TrimSpace(item.Label)
		href := strings.TrimSpace(item.Href)
		if label == "" || href == "" {
			continue
		}
		out = append(out, NavItem{Label: label, Href: href})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func navOrDefault(nav []NavItem) []NavItem {
	if len(nav) == 0 {
		return defaultNav()
	}
	out := make([]NavItem, len(nav))
	copy(out, nav)
	return out
}

type Settings struct {
	SiteName        string    `json:"siteName"`
	Tagline         string    `json:"tagline"`
	About           string    `json:"about"`
	SupportEmail    string    `json:"supportEmail"`
	SignupEnabled   bool      `json:"signupEnabled"`
	Maintenance     bool      `json:"maintenance"`
	MetaTitle       string    `json:"metaTitle"`
	MetaDescription string    `json:"metaDescription"`
	OgImage         string    `json:"ogImage"`
	Favicon         string    `json:"favicon"`
	ThemeColor      string    `json:"themeColor"`
	HeroTitle       string    `json:"heroTitle"`
	HeroSubtitle    string    `json:"heroSubtitle"`
	HeroCta         string    `json:"heroCta"`
	HeroCtaHref     string    `json:"heroCtaHref"`
	HeroImage       string    `json:"heroImage"`
	DemoSlug        string    `json:"demoSlug"`
	Nav             []NavItem `json:"nav"`
}

func (s *Service) Settings() Settings {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, err := s.store.GetSettings(context.Background())
	if err != nil {
		return DefaultSettings()
	}
	return st
}

func (s *Service) PublicSettings() map[string]any {
	st := s.Settings()
	return map[string]any{
		"siteName":        st.SiteName,
		"tagline":         st.Tagline,
		"about":           st.About,
		"supportEmail":    st.SupportEmail,
		"signupEnabled":   st.SignupEnabled,
		"maintenance":     st.Maintenance,
		"metaTitle":       st.MetaTitle,
		"metaDescription": st.MetaDescription,
		"ogImage":         st.OgImage,
		"favicon":         st.Favicon,
		"themeColor":      st.ThemeColor,
		"heroTitle":       st.HeroTitle,
		"heroSubtitle":    st.HeroSubtitle,
		"heroCta":         st.HeroCta,
		"heroCtaHref":     st.HeroCtaHref,
		"heroImage":       st.HeroImage,
		"demoSlug":        st.DemoSlug,
		"nav":             navOrDefault(st.Nav),
	}
}

func (s *Service) UpdateSettings(in Settings) Settings {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, err := s.store.GetSettings(context.Background())
	if err != nil {
		st = DefaultSettings()
	}
	if strings.TrimSpace(in.SiteName) != "" {
		st.SiteName = strings.TrimSpace(in.SiteName)
	}
	st.Tagline = strings.TrimSpace(in.Tagline)
	st.About = strings.TrimSpace(in.About)
	st.SupportEmail = strings.TrimSpace(in.SupportEmail)
	st.SignupEnabled = in.SignupEnabled
	st.Maintenance = in.Maintenance
	st.MetaTitle = strings.TrimSpace(in.MetaTitle)
	st.MetaDescription = strings.TrimSpace(in.MetaDescription)
	st.OgImage = strings.TrimSpace(in.OgImage)
	st.Favicon = strings.TrimSpace(in.Favicon)
	st.ThemeColor = strings.TrimSpace(in.ThemeColor)
	if v := strings.TrimSpace(in.HeroTitle); v != "" {
		st.HeroTitle = v
	}
	if v := strings.TrimSpace(in.HeroSubtitle); v != "" {
		st.HeroSubtitle = v
	}
	if v := strings.TrimSpace(in.HeroCta); v != "" {
		st.HeroCta = v
	}
	if v := strings.TrimSpace(in.HeroImage); v != "" {
		st.HeroImage = v
	}
	if v := strings.TrimSpace(in.DemoSlug); v != "" {
		st.DemoSlug = v
	}
	if v := strings.TrimSpace(in.HeroCtaHref); v != "" {
		st.HeroCtaHref = v
	}
	if nav := trimNav(in.Nav); len(nav) > 0 {
		st.Nav = nav
	}
	_ = s.store.PutSettings(context.Background(), st)
	return st
}
