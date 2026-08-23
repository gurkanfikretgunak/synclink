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

func scanPage(row interface{ Scan(dest ...any) error }) (*Page, error) {
	var p Page
	var avatar sql.NullString
	var socials sql.NullString
	var created, updated string
	err := row.Scan(
		&p.ID, &p.UserID, &p.Slug, &p.DisplayName, &p.Bio, &avatar, &p.Theme,
		&p.AvatarShape, &p.AccentColor, &p.Background, &p.Motion, &socials, &created, &updated,
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
	var avatar any
	if p.AvatarURL != nil {
		avatar = *p.AvatarURL
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO pages (id, user_id, slug, display_name, bio, avatar_url, theme, avatar_shape, accent_color, background, motion, socials, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID.String(), p.UserID.String(), p.Slug, p.DisplayName, p.Bio, avatar, p.Theme,
		p.AvatarShape, p.AccentColor, p.Background, p.Motion, encodeSocials(p.Socials),
		p.CreatedAt.Format(time.RFC3339Nano), p.UpdatedAt.Format(time.RFC3339Nano),
	)
	if isUniqueErr(err) {
		return ErrConflict
	}
	return err
}

func (s *SQLiteStore) UpdatePage(ctx context.Context, p *Page) error {
	p.UpdatedAt = time.Now().UTC()
	var avatar any
	if p.AvatarURL != nil {
		avatar = *p.AvatarURL
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE pages SET slug=?, display_name=?, bio=?, avatar_url=?, theme=?, avatar_shape=?, accent_color=?, background=?, motion=?, socials=?, updated_at=?, user_id=?
		WHERE id=?`,
		p.Slug, p.DisplayName, p.Bio, avatar, p.Theme, p.AvatarShape, p.AccentColor, p.Background, p.Motion,
		encodeSocials(p.Socials), p.UpdatedAt.Format(time.RFC3339Nano), p.UserID.String(), p.ID.String(),
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

const pageCols = `id, user_id, slug, display_name, bio, avatar_url, theme, avatar_shape, accent_color, background, motion, socials, created_at, updated_at`

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
	var lastClicked sql.NullString
	var active int
	var created, updated string
	err := row.Scan(&l.ID, &l.PageID, &l.Title, &l.URL, &icon, &l.Order, &active, &l.Clicks, &lastClicked, &created, &updated)
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
	l.Active = active != 0
	if lastClicked.Valid && lastClicked.String != "" {
		if ts, err := time.Parse(time.RFC3339Nano, lastClicked.String); err == nil {
			l.LastClickedAt = &ts
		} else if ts, err := time.Parse(time.RFC3339, lastClicked.String); err == nil {
			l.LastClickedAt = &ts
		}
	}
	l.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	l.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return &l, nil
}

func (s *SQLiteStore) CreateLink(ctx context.Context, l *Link) error {
	now := time.Now().UTC()
	l.CreatedAt, l.UpdatedAt = now, now
	var icon any
	if l.Icon != nil {
		icon = *l.Icon
	}
	active := 0
	if l.Active {
		active = 1
	}
	var lastClicked any
	if l.LastClickedAt != nil {
		lastClicked = l.LastClickedAt.UTC().Format(time.RFC3339Nano)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO links (id, page_id, title, url, icon, sort_order, active, clicks, last_clicked_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		l.ID.String(), l.PageID.String(), l.Title, l.URL, icon, l.Order, active, l.Clicks, lastClicked,
		l.CreatedAt.Format(time.RFC3339Nano), l.UpdatedAt.Format(time.RFC3339Nano),
	)
	return err
}

func (s *SQLiteStore) UpdateLink(ctx context.Context, l *Link) error {
	l.UpdatedAt = time.Now().UTC()
	var icon any
	if l.Icon != nil {
		icon = *l.Icon
	}
	active := 0
	if l.Active {
		active = 1
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE links SET title=?, url=?, icon=?, sort_order=?, active=?, updated_at=? WHERE id=?`,
		l.Title, l.URL, icon, l.Order, active, l.UpdatedAt.Format(time.RFC3339Nano), l.ID.String(),
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

const linkCols = `id, page_id, title, url, icon, sort_order, active, clicks, last_clicked_at, created_at, updated_at`

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
	var n int
	now := time.Now().UTC().Format(time.RFC3339Nano)
	err := s.db.QueryRowContext(ctx, `
		UPDATE links SET clicks = clicks + 1, last_clicked_at = ?, updated_at = ?
		WHERE id=? AND page_id=? AND active=1
		RETURNING clicks`,
		now, now, linkID.String(), pageID.String(),
	).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	return n, nil
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
