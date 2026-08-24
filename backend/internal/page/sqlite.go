package page

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

func isUniqueErr(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "unique constraint") || strings.Contains(s, "constraint failed")
}

func nullableStr(p *string) any {
	if p == nil {
		return nil
	}
	return *p
}

func nullableTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.RFC3339Nano)
}

func parseTimePtr(raw sql.NullString) *time.Time {
	if !raw.Valid || raw.String == "" {
		return nil
	}
	if ts, err := time.Parse(time.RFC3339Nano, raw.String); err == nil {
		return &ts
	}
	if ts, err := time.Parse(time.RFC3339, raw.String); err == nil {
		return &ts
	}
	return nil
}

func scanPage(row interface{ Scan(dest ...any) error }) (*Page, error) {
	var p Page
	var avatar sql.NullString
	var socials sql.NullString
	var password sql.NullString
	var published sql.NullString
	var cover sql.NullString
	var verified int
	var created, updated string
	err := row.Scan(
		&p.ID, &p.UserID, &p.Slug, &p.DisplayName, &p.Bio, &avatar, &p.Theme,
		&p.AvatarShape, &p.AccentColor, &p.Background, &p.Motion, &socials,
		&verified, &password, &published, &cover, &p.CoverKind, &created, &updated,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if avatar.Valid {
		v := avatar.String
		p.AvatarURL = &v
	}
	p.Socials = decodeSocials(socials)
	p.Verified = verified != 0
	if password.Valid && password.String != "" {
		v := password.String
		p.PagePassword = &v
	}
	p.PublishedAt = parseTimePtr(published)
	if cover.Valid && cover.String != "" {
		v := cover.String
		p.CoverURL = &v
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return &p, nil
}

func decodeSocials(raw sql.NullString) []Social {
	if !raw.Valid || strings.TrimSpace(raw.String) == "" || raw.String == "null" {
		return emptySocials()
	}
	var out []Social
	if err := json.Unmarshal([]byte(raw.String), &out); err != nil {
		return emptySocials()
	}
	return NormalizeSocials(out)
}

func encodeSocials(in []Social) string {
	b, err := json.Marshal(NormalizeSocials(in))
	if err != nil {
		return "[]"
	}
	return string(b)
}

func (s *SQLiteStore) CreatePage(ctx context.Context, p *Page) error {
	now := time.Now().UTC()
	p.CreatedAt, p.UpdatedAt = now, now
	if p.PublishedAt == nil {
		p.PublishedAt = &now
	}
	verified := 0
	if p.Verified {
		verified = 1
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO pages (id, user_id, slug, display_name, bio, avatar_url, theme, avatar_shape, accent_color, background, motion, socials, verified, page_password, published_at, cover_url, cover_kind, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID.String(), p.UserID.String(), p.Slug, p.DisplayName, p.Bio, nullableStr(p.AvatarURL), p.Theme,
		p.AvatarShape, p.AccentColor, p.Background, p.Motion, encodeSocials(p.Socials),
		verified, nullableStr(p.PagePassword), nullableTime(p.PublishedAt), nullableStr(p.CoverURL), p.CoverKind,
		p.CreatedAt.Format(time.RFC3339Nano), p.UpdatedAt.Format(time.RFC3339Nano),
	)
	if isUniqueErr(err) {
		return ErrConflict
	}
	return err
}

func (s *SQLiteStore) UpdatePage(ctx context.Context, p *Page) error {
	p.UpdatedAt = time.Now().UTC()
	if p.PublishedAt == nil {
		now := p.UpdatedAt
		p.PublishedAt = &now
	}
	verified := 0
	if p.Verified {
		verified = 1
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE pages SET slug=?, display_name=?, bio=?, avatar_url=?, theme=?, avatar_shape=?, accent_color=?, background=?, motion=?, socials=?, verified=?, page_password=?, published_at=?, cover_url=?, cover_kind=?, updated_at=?, user_id=?
		WHERE id=?`,
		p.Slug, p.DisplayName, p.Bio, nullableStr(p.AvatarURL), p.Theme, p.AvatarShape, p.AccentColor, p.Background, p.Motion,
		encodeSocials(p.Socials), verified, nullableStr(p.PagePassword), nullableTime(p.PublishedAt), nullableStr(p.CoverURL), p.CoverKind, p.UpdatedAt.Format(time.RFC3339Nano), p.UserID.String(), p.ID.String(),
	)
	if err != nil {
		if isUniqueErr(err) {
			return ErrConflict
		}
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

const pageCols = `id, user_id, slug, display_name, bio, avatar_url, theme, avatar_shape, accent_color, background, motion, socials, verified, page_password, published_at, cover_url, cover_kind, created_at, updated_at`

func (s *SQLiteStore) GetPageByID(ctx context.Context, id uuid.UUID) (*Page, error) {
	return scanPage(s.db.QueryRowContext(ctx, `SELECT `+pageCols+` FROM pages WHERE id=?`, id.String()))
}

func (s *SQLiteStore) GetPageByUserID(ctx context.Context, userID uuid.UUID) (*Page, error) {
	return scanPage(s.db.QueryRowContext(ctx, `SELECT `+pageCols+` FROM pages WHERE user_id=?`, userID.String()))
}

func (s *SQLiteStore) GetPageBySlug(ctx context.Context, slug string) (*Page, error) {
	return scanPage(s.db.QueryRowContext(ctx, `SELECT `+pageCols+` FROM pages WHERE slug=?`, slug))
}

func (s *SQLiteStore) ListPages(ctx context.Context) ([]*Page, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+pageCols+` FROM pages`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Page
	for rows.Next() {
		p, err := scanPage(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if out == nil {
		out = []*Page{}
	}
	return out, rows.Err()
}

func scanLink(row interface{ Scan(dest ...any) error }) (*Link, error) {
	var l Link
	var icon sql.NullString
	var thumb sql.NullString
	var lastClicked sql.NullString
	var starts, ends sql.NullString
	var embed sql.NullString
	var active, featured, sensitive int
	var created, updated string
	err := row.Scan(
		&l.ID, &l.PageID, &l.Title, &l.URL, &icon, &l.Order, &active, &l.Clicks, &lastClicked,
		&featured, &thumb, &starts, &ends, &sensitive, &l.Section, &embed, &created, &updated,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if icon.Valid {
		v := icon.String
		l.Icon = &v
	}
	if thumb.Valid && thumb.String != "" {
		v := thumb.String
		l.ThumbnailURL = &v
	}
	l.Active = active != 0
	l.Featured = featured != 0
	l.Sensitive = sensitive != 0
	l.LastClickedAt = parseTimePtr(lastClicked)
	l.StartsAt = parseTimePtr(starts)
	l.EndsAt = parseTimePtr(ends)
	if embed.Valid && embed.String != "" {
		v := embed.String
		l.EmbedURL = &v
	}
	l.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	l.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return &l, nil
}

func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func (s *SQLiteStore) CreateLink(ctx context.Context, l *Link) error {
	now := time.Now().UTC()
	l.CreatedAt, l.UpdatedAt = now, now
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO links (id, page_id, title, url, icon, sort_order, active, clicks, last_clicked_at, featured, thumbnail_url, starts_at, ends_at, sensitive, section, embed_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		l.ID.String(), l.PageID.String(), l.Title, l.URL, nullableStr(l.Icon), l.Order, boolInt(l.Active), l.Clicks,
		nullableTime(l.LastClickedAt), boolInt(l.Featured), nullableStr(l.ThumbnailURL),
		nullableTime(l.StartsAt), nullableTime(l.EndsAt), boolInt(l.Sensitive), l.Section, nullableStr(l.EmbedURL),
		l.CreatedAt.Format(time.RFC3339Nano), l.UpdatedAt.Format(time.RFC3339Nano),
	)
	return err
}

func (s *SQLiteStore) UpdateLink(ctx context.Context, l *Link) error {
	l.UpdatedAt = time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `
		UPDATE links SET title=?, url=?, icon=?, sort_order=?, active=?, featured=?, thumbnail_url=?, starts_at=?, ends_at=?, sensitive=?, section=?, embed_url=?, updated_at=? WHERE id=?`,
		l.Title, l.URL, nullableStr(l.Icon), l.Order, boolInt(l.Active), boolInt(l.Featured), nullableStr(l.ThumbnailURL),
		nullableTime(l.StartsAt), nullableTime(l.EndsAt), boolInt(l.Sensitive), l.Section, nullableStr(l.EmbedURL),
		l.UpdatedAt.Format(time.RFC3339Nano), l.ID.String(),
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLiteStore) DeleteLink(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM links WHERE id=?`, id.String())
	return err
}

const linkCols = `id, page_id, title, url, icon, sort_order, active, clicks, last_clicked_at, featured, thumbnail_url, starts_at, ends_at, sensitive, section, embed_url, created_at, updated_at`

func (s *SQLiteStore) GetLinkByID(ctx context.Context, id uuid.UUID) (*Link, error) {
	return scanLink(s.db.QueryRowContext(ctx, `SELECT `+linkCols+` FROM links WHERE id=?`, id.String()))
}

func (s *SQLiteStore) ListLinks(ctx context.Context, pageID uuid.UUID) ([]*Link, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+linkCols+` FROM links WHERE page_id=? ORDER BY sort_order ASC`, pageID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Link
	for rows.Next() {
		l, err := scanLink(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	if out == nil {
		out = []*Link{}
	}
	return out, rows.Err()
}

func (s *SQLiteStore) MaxOrder(ctx context.Context, pageID uuid.UUID) (int, error) {
	var max sql.NullInt64
	err := s.db.QueryRowContext(ctx, `SELECT MAX(sort_order) FROM links WHERE page_id=?`, pageID.String()).Scan(&max)
	if err != nil {
		return -1, err
	}
	if !max.Valid {
		return -1, nil
	}
	return int(max.Int64), nil
}

func (s *SQLiteStore) Reorder(ctx context.Context, pageID uuid.UUID, ids []uuid.UUID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	for i, id := range ids {
		var gotPage string
		err := tx.QueryRowContext(ctx, `SELECT page_id FROM links WHERE id=?`, id.String()).Scan(&gotPage)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if gotPage != pageID.String() {
			return ErrNotFound
		}
		if _, err := tx.ExecContext(ctx, `UPDATE links SET sort_order=?, updated_at=? WHERE id=?`, i, now, id.String()); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *SQLiteStore) IncrementClicks(ctx context.Context, pageID, linkID uuid.UUID) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()
	var n int
	now := time.Now().UTC()
	err = tx.QueryRowContext(ctx, `
		UPDATE links SET clicks = clicks + 1, last_clicked_at = ?, updated_at = ?
		WHERE id=? AND page_id=? AND active=1
		RETURNING clicks`,
		now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), linkID.String(), pageID.String(),
	).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO click_days (page_id, day, clicks) VALUES (?, ?, 1)
		ON CONFLICT(page_id, day) DO UPDATE SET clicks = clicks + 1`,
		pageID.String(), now.Format("2006-01-02"),
	)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return n, nil
}

func (s *SQLiteStore) DailyClicks(ctx context.Context, pageID uuid.UUID) (map[string]int, error) {
	since := time.Now().UTC().AddDate(0, 0, -13).Format("2006-01-02")
	rows, err := s.db.QueryContext(ctx, `
		SELECT day, clicks FROM click_days WHERE page_id=? AND day>=?`,
		pageID.String(), since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var day string
		var n int
		if err := rows.Scan(&day, &n); err != nil {
			return nil, err
		}
		out[day] = n
	}
	return out, rows.Err()
}

func (s *SQLiteStore) BumpReferrer(ctx context.Context, pageID uuid.UUID, host string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO click_referrers (page_id, host, clicks) VALUES (?, ?, 1)
		ON CONFLICT(page_id, host) DO UPDATE SET clicks = clicks + 1`,
		pageID.String(), host,
	)
	return err
}

func (s *SQLiteStore) Referrers(ctx context.Context, pageID uuid.UUID) (map[string]int, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT host, clicks FROM click_referrers WHERE page_id=?`,
		pageID.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var host string
		var n int
		if err := rows.Scan(&host, &n); err != nil {
			return nil, err
		}
		out[host] = n
	}
	return out, rows.Err()
}

func (s *SQLiteStore) SumClicks(ctx context.Context) (int, error) {
	var n sql.NullInt64
	err := s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(clicks), 0) FROM links`).Scan(&n)
	if err != nil {
		return 0, err
	}
	if !n.Valid {
		return 0, nil
	}
	return int(n.Int64), nil
}

func (s *SQLiteStore) ClicksByPage(ctx context.Context) (map[uuid.UUID]int, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT p.id, COALESCE(SUM(l.clicks), 0)
		FROM pages p
		LEFT JOIN links l ON l.page_id = p.id
		GROUP BY p.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[uuid.UUID]int{}
	for rows.Next() {
		var idStr string
		var n int
		if err := rows.Scan(&idStr, &n); err != nil {
			return nil, err
		}
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		out[id] = n
	}
	return out, rows.Err()
}

func (s *SQLiteStore) CreateSubscriber(ctx context.Context, sub *Subscriber) error {
	sub.CreatedAt = time.Now().UTC()
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO subscribers (id, page_id, email, created_at) VALUES (?, ?, ?, ?)`,
		sub.ID.String(), sub.PageID.String(), sub.Email, sub.CreatedAt.Format(time.RFC3339Nano),
	)
	if isUniqueErr(err) {
		return ErrConflict
	}
	return err
}

func scanSubscriber(row interface{ Scan(dest ...any) error }) (*Subscriber, error) {
	var sub Subscriber
	var created string
	err := row.Scan(&sub.ID, &sub.PageID, &sub.Email, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	sub.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return &sub, nil
}

func (s *SQLiteStore) ListSubscribers(ctx context.Context, pageID uuid.UUID) ([]*Subscriber, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, page_id, email, created_at FROM subscribers WHERE page_id=? ORDER BY created_at DESC`, pageID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Subscriber
	for rows.Next() {
		sub, err := scanSubscriber(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sub)
	}
	if out == nil {
		out = []*Subscriber{}
	}
	return out, rows.Err()
}

func (s *SQLiteStore) GetSubscriberByID(ctx context.Context, id uuid.UUID) (*Subscriber, error) {
	return scanSubscriber(s.db.QueryRowContext(ctx, `SELECT id, page_id, email, created_at FROM subscribers WHERE id=?`, id.String()))
}

func (s *SQLiteStore) DeleteSubscriber(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM subscribers WHERE id=?`, id.String())
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
