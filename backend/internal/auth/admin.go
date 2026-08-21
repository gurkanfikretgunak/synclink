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
	out := make([]UserInfo, 0, len(s.byID))
	for _, u := range s.byID {
		out = append(out, u.Info())
	}
	return out
}

func (s *Service) IsAdmin(id uuid.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.byID[id]
	return ok && u.Role == RoleAdmin && u.Status == StatusActive
}

func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, role, status *string) (*UserInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.byID[id]
	if !ok {
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
	info := u.Info()
	return &info, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.byID[id]
	if !ok {
		return ErrInvalidCreds
	}
	delete(s.byID, id)
	delete(s.users, u.Email)
	return nil
}

type Settings struct {
	SiteName        string `json:"siteName"`
	Tagline         string `json:"tagline"`
	About           string `json:"about"`
	SupportEmail    string `json:"supportEmail"`
	SignupEnabled   bool   `json:"signupEnabled"`
	Maintenance     bool   `json:"maintenance"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	OgImage         string `json:"ogImage"`
	Favicon         string `json:"favicon"`
	ThemeColor      string `json:"themeColor"`
	HeroTitle       string `json:"heroTitle"`
	HeroSubtitle    string `json:"heroSubtitle"`
	HeroCta         string `json:"heroCta"`
	HeroImage       string `json:"heroImage"`
	DemoSlug        string `json:"demoSlug"`
}

func (s *Service) Settings() Settings {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.settings
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
		"heroImage":       st.HeroImage,
		"demoSlug":        st.DemoSlug,
	}
}

func (s *Service) UpdateSettings(in Settings) Settings {
	s.mu.Lock()
	defer s.mu.Unlock()
	if strings.TrimSpace(in.SiteName) != "" {
		s.settings.SiteName = strings.TrimSpace(in.SiteName)
	}
	s.settings.Tagline = strings.TrimSpace(in.Tagline)
	s.settings.About = strings.TrimSpace(in.About)
	s.settings.SupportEmail = strings.TrimSpace(in.SupportEmail)
	s.settings.SignupEnabled = in.SignupEnabled
	s.settings.Maintenance = in.Maintenance
	s.settings.MetaTitle = strings.TrimSpace(in.MetaTitle)
	s.settings.MetaDescription = strings.TrimSpace(in.MetaDescription)
	s.settings.OgImage = strings.TrimSpace(in.OgImage)
	s.settings.Favicon = strings.TrimSpace(in.Favicon)
	s.settings.ThemeColor = strings.TrimSpace(in.ThemeColor)
	if v := strings.TrimSpace(in.HeroTitle); v != "" {
		s.settings.HeroTitle = v
	}
	if v := strings.TrimSpace(in.HeroSubtitle); v != "" {
		s.settings.HeroSubtitle = v
	}
	if v := strings.TrimSpace(in.HeroCta); v != "" {
		s.settings.HeroCta = v
	}
	if v := strings.TrimSpace(in.HeroImage); v != "" {
		s.settings.HeroImage = v
	}
	if v := strings.TrimSpace(in.DemoSlug); v != "" {
		s.settings.DemoSlug = v
	}
	return s.settings
}
