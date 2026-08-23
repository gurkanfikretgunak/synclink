package page

import (
	"context"
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var slugRE = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,62}[a-z0-9])?$`)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func normalizeSlug(slug string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(slug))
	if !slugRE.MatchString(s) {
		return "", ErrValidation
	}
	return s, nil
}

func isNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func normalizePagePassword(in *string) *string {
	if in == nil {
		return nil
	}
	v := strings.TrimSpace(*in)
	if v == "" {
		return nil
	}
	return &v
}

func (s *Service) GetPublicPage(ctx context.Context, slug string) (*PublicPage, error) {
	return s.GetPublicPageWithPassword(ctx, slug, "")
}

func (s *Service) GetPublicPageWithPassword(ctx context.Context, slug, password string) (*PublicPage, error) {
	normalized, err := normalizeSlug(slug)
	if err != nil {
		return nil, ErrNotFound
	}
	p, err := s.store.GetPageBySlug(ctx, normalized)
	if err != nil {
		return nil, ErrNotFound
	}
	if p.PagePassword != nil && *p.PagePassword != "" {
		if password != *p.PagePassword {
			return nil, ErrLocked
		}
	}
	links, err := s.store.ListLinks(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	out := make([]PublicLink, 0, len(links))
	for _, l := range links {
		if !l.Active || !LinkInSchedule(l, now) {
			continue
		}
		out = append(out, ToPublicLink(l))
	}
	return &PublicPage{
		Slug: p.Slug, DisplayName: p.DisplayName, Bio: p.Bio,
		AvatarURL: p.AvatarURL, Theme: p.Theme, AvatarShape: p.AvatarShape,
		AccentColor: p.AccentColor, Background: p.Background, Motion: p.Motion,
		Socials: copySocials(p.Socials), Verified: p.Verified, PublishedAt: copyTime(p.PublishedAt), Links: out,
	}, nil
}

func emptyPageDTO() PageDTO {
	shape, accent, bg, motion := NormalizeLook("", "", "", "")
	return PageDTO{
		ID:          uuid.Nil,
		Slug:        "",
		DisplayName: "",
		Bio:         "",
		AvatarURL:   nil,
		Theme:       ThemeDefault,
		AvatarShape: shape,
		AccentColor: accent,
		Background:  bg,
		Motion:      motion,
		Socials:     emptySocials(),
		Verified:    false,
	}
}

func (s *Service) GetMyPage(ctx context.Context, userID uuid.UUID) (*PageDTO, error) {
	p, err := s.store.GetPageByUserID(ctx, userID)
	if err != nil {
		if isNotFound(err) {
			dto := emptyPageDTO()
			return &dto, nil
		}
		return nil, err
	}
	dto := ToPageDTO(p)
	return &dto, nil
}

func (s *Service) UpsertPage(ctx context.Context, userID uuid.UUID, in UpsertPageInput) (*PageDTO, error) {
	slug, err := normalizeSlug(in.Slug)
	if err != nil {
		return nil, ErrValidation
	}
	theme := strings.TrimSpace(in.Theme)
	if theme == "" {
		theme = ThemeDefault
	}
	if !ValidTheme(theme) {
		return nil, ErrValidation
	}
	shape, accent, bg, motion := NormalizeLook(strings.TrimSpace(in.AvatarShape), strings.TrimSpace(in.AccentColor), strings.TrimSpace(in.Background), strings.TrimSpace(in.Motion))
	socials := NormalizeSocials(in.Socials)
	pw := normalizePagePassword(in.PagePassword)
	bySlug, slugErr := s.store.GetPageBySlug(ctx, slug)
	if slugErr != nil && !isNotFound(slugErr) {
		return nil, slugErr
	}
	mine, mineErr := s.store.GetPageByUserID(ctx, userID)
	if mineErr != nil && !isNotFound(mineErr) {
		return nil, mineErr
	}
	if bySlug != nil && (mine == nil || bySlug.ID != mine.ID) {
		return nil, ErrConflict
	}
	if mine != nil {
		mine.Slug = slug
		mine.DisplayName = strings.TrimSpace(in.DisplayName)
		mine.Bio = strings.TrimSpace(in.Bio)
		mine.AvatarURL = in.AvatarURL
		mine.Theme = theme
		mine.AvatarShape = shape
		mine.AccentColor = accent
		mine.Background = bg
		mine.Motion = motion
		mine.Socials = socials
		if in.PagePassword != nil {
			mine.PagePassword = pw
		}
		if mine.PublishedAt == nil {
			now := time.Now().UTC()
			mine.PublishedAt = &now
		}
		if err := s.store.UpdatePage(ctx, mine); err != nil {
			return nil, err
		}
		dto := ToPageDTO(mine)
		return &dto, nil
	}
	p := &Page{
		ID: uuid.New(), UserID: userID, Slug: slug,
		DisplayName: strings.TrimSpace(in.DisplayName),
		Bio:         strings.TrimSpace(in.Bio), AvatarURL: in.AvatarURL, Theme: theme,
		AvatarShape: shape, AccentColor: accent, Background: bg, Motion: motion,
		Socials: socials, Verified: false, PagePassword: pw,
	}
	if err := s.store.CreatePage(ctx, p); err != nil {
		return nil, err
	}
	dto := ToPageDTO(p)
	return &dto, nil
}

func (s *Service) SetPageVerified(ctx context.Context, pageID uuid.UUID, verified bool) (*PageDTO, error) {
	p, err := s.store.GetPageByID(ctx, pageID)
	if err != nil {
		return nil, ErrNotFound
	}
	p.Verified = verified
	if err := s.store.UpdatePage(ctx, p); err != nil {
		return nil, err
	}
	dto := ToPageDTO(p)
	return &dto, nil
}

func (s *Service) requirePage(ctx context.Context, userID uuid.UUID) (*Page, error) {
	p, err := s.store.GetPageByUserID(ctx, userID)
	if err != nil {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *Service) ListLinks(ctx context.Context, userID uuid.UUID) ([]LinkDTO, error) {
	p, err := s.requirePage(ctx, userID)
	if err != nil {
		if isNotFound(err) {
			return []LinkDTO{}, nil
		}
		return nil, err
	}
	links, err := s.store.ListLinks(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	out := make([]LinkDTO, 0, len(links))
	for _, l := range links {
		out = append(out, ToLinkDTO(l))
	}
	return out, nil
}

func applyCreateLinkExtras(l *Link, in CreateLinkInput) {
	if in.Featured != nil {
		l.Featured = *in.Featured
	}
	l.ThumbnailURL = normalizeOptionalString(in.ThumbnailURL)
	l.StartsAt = in.StartsAt
	l.EndsAt = in.EndsAt
	if in.Sensitive != nil {
		l.Sensitive = *in.Sensitive
	}
}

func (s *Service) CreateLink(ctx context.Context, userID uuid.UUID, in CreateLinkInput) (*LinkDTO, error) {
	p, err := s.requirePage(ctx, userID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.URL) == "" {
		return nil, ErrValidation
	}
	max, err := s.store.MaxOrder(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	active := true
	if in.Active != nil {
		active = *in.Active
	}
	l := &Link{
		ID: uuid.New(), PageID: p.ID, Title: strings.TrimSpace(in.Title),
		URL: strings.TrimSpace(in.URL), Icon: in.Icon, Order: max + 1, Active: active,
	}
	applyCreateLinkExtras(l, in)
	if err := s.store.CreateLink(ctx, l); err != nil {
		return nil, err
	}
	dto := ToLinkDTO(l)
	return &dto, nil
}

func (s *Service) UpdateLink(ctx context.Context, userID, linkID uuid.UUID, in UpdateLinkInput) (*LinkDTO, error) {
	p, err := s.requirePage(ctx, userID)
	if err != nil {
		return nil, err
	}
	l, err := s.store.GetLinkByID(ctx, linkID)
	if err != nil || l.PageID != p.ID {
		return nil, ErrNotFound
	}
	if in.Title != nil {
		l.Title = strings.TrimSpace(*in.Title)
	}
	if in.URL != nil {
		l.URL = strings.TrimSpace(*in.URL)
	}
	if in.Icon != nil {
		l.Icon = in.Icon
	}
	if in.Active != nil {
		l.Active = *in.Active
	}
	if in.Featured != nil {
		l.Featured = *in.Featured
	}
	if in.ThumbnailURL != nil {
		l.ThumbnailURL = normalizeOptionalString(in.ThumbnailURL)
	}
	if in.StartsAt != nil {
		l.StartsAt = in.StartsAt
	}
	if in.EndsAt != nil {
		l.EndsAt = in.EndsAt
	}
	if in.Sensitive != nil {
		l.Sensitive = *in.Sensitive
	}
	if err := s.store.UpdateLink(ctx, l); err != nil {
		return nil, err
	}
	dto := ToLinkDTO(l)
	return &dto, nil
}

func (s *Service) DeleteLink(ctx context.Context, userID, linkID uuid.UUID) error {
	p, err := s.requirePage(ctx, userID)
	if err != nil {
		return err
	}
	l, err := s.store.GetLinkByID(ctx, linkID)
	if err != nil || l.PageID != p.ID {
		return ErrNotFound
	}
	return s.store.DeleteLink(ctx, l.ID)
}

func (s *Service) ReorderLinks(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) ([]LinkDTO, error) {
	p, err := s.requirePage(ctx, userID)
	if err != nil {
		return nil, err
	}
	existing, err := s.store.ListLinks(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	owned := map[uuid.UUID]struct{}{}
	for _, l := range existing {
		owned[l.ID] = struct{}{}
	}
	if len(ids) != len(existing) {
		return nil, ErrValidation
	}
	seen := map[uuid.UUID]struct{}{}
	for _, id := range ids {
		if _, ok := owned[id]; !ok {
			return nil, ErrNotFound
		}
		if _, dup := seen[id]; dup {
			return nil, ErrValidation
		}
		seen[id] = struct{}{}
	}
	if err := s.store.Reorder(ctx, p.ID, ids); err != nil {
		return nil, err
	}
	return s.ListLinks(ctx, userID)
}

func (s *Service) ListAll(ctx context.Context) ([]PageDTO, error) {
	pages, err := s.store.ListPages(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]PageDTO, 0, len(pages))
	for _, p := range pages {
		out = append(out, ToPageDTO(p))
	}
	return out, nil
}

func (s *Service) RecordClick(ctx context.Context, slug string, linkID uuid.UUID) (int, error) {
	normalized, err := normalizeSlug(slug)
	if err != nil {
		return 0, ErrNotFound
	}
	p, err := s.store.GetPageBySlug(ctx, normalized)
	if err != nil {
		return 0, ErrNotFound
	}
	l, err := s.store.GetLinkByID(ctx, linkID)
	if err != nil || l.PageID != p.ID {
		return 0, ErrNotFound
	}
	if !l.Active || !LinkInSchedule(l, time.Now().UTC()) {
		return 0, ErrNotFound
	}
	return s.store.IncrementClicks(ctx, p.ID, linkID)
}

func (s *Service) MyStats(ctx context.Context, userID uuid.UUID) (*MyStats, error) {
	links, err := s.ListLinks(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := &MyStats{TotalClicks: 0, Links: make([]MyStatsLink, 0, len(links))}
	for _, l := range links {
		out.TotalClicks += l.Clicks
		out.Links = append(out.Links, MyStatsLink{ID: l.ID, Title: l.Title, Clicks: l.Clicks, URL: l.URL})
	}
	return out, nil
}

func (s *Service) SumClicks(ctx context.Context) (int, error) {
	return s.store.SumClicks(ctx)
}

func normalizeSubscribeEmail(raw string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(raw))
	if email == "" || strings.ContainsAny(email, " \t\n\r") {
		return "", ErrValidation
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address == "" || !strings.Contains(addr.Address, ".") {
		return "", ErrValidation
	}
	return strings.ToLower(addr.Address), nil
}

func (s *Service) Subscribe(ctx context.Context, slug, email string) error {
	normalized, err := normalizeSlug(slug)
	if err != nil {
		return ErrNotFound
	}
	p, err := s.store.GetPageBySlug(ctx, normalized)
	if err != nil {
		return ErrNotFound
	}
	addr, err := normalizeSubscribeEmail(email)
	if err != nil {
		return err
	}
	sub := &Subscriber{ID: uuid.New(), PageID: p.ID, Email: addr}
	return s.store.CreateSubscriber(ctx, sub)
}

func (s *Service) ListMySubscribers(ctx context.Context, userID uuid.UUID) ([]SubscriberDTO, error) {
	p, err := s.requirePage(ctx, userID)
	if err != nil {
		if isNotFound(err) {
			return []SubscriberDTO{}, nil
		}
		return nil, err
	}
	subs, err := s.store.ListSubscribers(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	out := make([]SubscriberDTO, 0, len(subs))
	for _, sub := range subs {
		out = append(out, ToSubscriberDTO(sub))
	}
	return out, nil
}

func (s *Service) DeleteMySubscriber(ctx context.Context, userID, subID uuid.UUID) error {
	p, err := s.requirePage(ctx, userID)
	if err != nil {
		return err
	}
	sub, err := s.store.GetSubscriberByID(ctx, subID)
	if err != nil || sub.PageID != p.ID {
		return ErrNotFound
	}
	return s.store.DeleteSubscriber(ctx, sub.ID)
}
