package page

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) SeedDemoIfEmpty(ctx context.Context, userID uuid.UUID) error {
	pages, err := s.store.ListPages(ctx)
	if err != nil {
		return err
	}
	if len(pages) > 0 {
		return nil
	}
	if _, err := s.UpsertPage(ctx, userID, UpsertPageInput{
		Slug:        "gurkan",
		DisplayName: "Gürkan",
		Bio:         "Your links, one page.",
		Theme:       ThemeDefault,
		AvatarShape: "circle",
		AccentColor: "#111111",
		Background:  "cream",
		Motion:      "subtle",
		Socials: []Social{
			{Network: "github", URL: "https://github.com/gurkanfikretgunak"},
			{Network: "website", URL: "https://synclink-mocha.vercel.app"},
		},
	}); err != nil {
		return err
	}
	links := []CreateLinkInput{
		{Title: "GitHub", URL: "https://github.com/gurkanfikretgunak"},
		{Title: "Site", URL: "https://synclink-mocha.vercel.app"},
		{Title: "About", URL: "https://synclink-mocha.vercel.app/about"},
	}
	for _, in := range links {
		if _, err := s.CreateLink(ctx, userID, in); err != nil {
			return err
		}
	}
	return nil
}
