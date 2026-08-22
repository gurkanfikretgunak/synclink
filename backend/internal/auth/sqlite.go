package auth

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

func scanUser(row interface{ Scan(dest ...any) error }) (*User, error) {
	var u User
	var created string
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.Status, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	u.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return &u, nil
}

func (s *SQLiteStore) UserCount(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&n)
	return n, err
}

func (s *SQLiteStore) CreateUser(ctx context.Context, u *User) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, role, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		u.ID.String(), u.Email, u.PasswordHash, u.Role, u.Status, u.CreatedAt.UTC().Format(time.RFC3339Nano),
	)
	if isUniqueErr(err) {
		return ErrExists
	}
	return err
}

func (s *SQLiteStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return scanUser(s.db.QueryRowContext(ctx, `SELECT id, email, password_hash, role, status, created_at FROM users WHERE email=?`, email))
}

func (s *SQLiteStore) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return scanUser(s.db.QueryRowContext(ctx, `SELECT id, email, password_hash, role, status, created_at FROM users WHERE id=?`, id.String()))
}

func (s *SQLiteStore) ListUsers(ctx context.Context) ([]*User, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, email, password_hash, role, status, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	if out == nil {
		out = []*User{}
	}
	return out, rows.Err()
}

func (s *SQLiteStore) UpdateUser(ctx context.Context, u *User) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE users SET email=?, password_hash=?, role=?, status=? WHERE id=?`,
		u.Email, u.PasswordHash, u.Role, u.Status, u.ID.String(),
	)
	if err != nil {
		if isUniqueErr(err) {
			return ErrExists
		}
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLiteStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE id=?`, id.String())
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLiteStore) GetSettings(ctx context.Context) (Settings, error) {
	var raw string
	err := s.db.QueryRowContext(ctx, `SELECT payload FROM settings WHERE id=1`).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return DefaultSettings(), nil
	}
	if err != nil {
		return Settings{}, err
	}
	var st Settings
	if err := json.Unmarshal([]byte(raw), &st); err != nil {
		return Settings{}, err
	}
	return st, nil
}

func (s *SQLiteStore) PutSettings(ctx context.Context, st Settings) error {
	raw, err := json.Marshal(st)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO settings (id, payload) VALUES (1, ?) ON CONFLICT(id) DO UPDATE SET payload=excluded.payload`, string(raw))
	return err
}

func (s *SQLiteStore) PutResetToken(ctx context.Context, token string, userID uuid.UUID, exp time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO password_reset_tokens (token, user_id, expires_at) VALUES (?, ?, ?)
		ON CONFLICT(token) DO UPDATE SET user_id=excluded.user_id, expires_at=excluded.expires_at`,
		token, userID.String(), exp.UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (s *SQLiteStore) GetResetToken(ctx context.Context, token string) (uuid.UUID, time.Time, bool, error) {
	var uid string
	var exp string
	err := s.db.QueryRowContext(ctx, `SELECT user_id, expires_at FROM password_reset_tokens WHERE token=?`, token).Scan(&uid, &exp)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, time.Time{}, false, nil
	}
	if err != nil {
		return uuid.Nil, time.Time{}, false, err
	}
	id, err := uuid.Parse(uid)
	if err != nil {
		return uuid.Nil, time.Time{}, false, err
	}
	t, err := time.Parse(time.RFC3339Nano, exp)
	if err != nil {
		t, err = time.Parse(time.RFC3339, exp)
		if err != nil {
			return uuid.Nil, time.Time{}, false, err
		}
	}
	return id, t, true, nil
}

func (s *SQLiteStore) DeleteResetToken(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM password_reset_tokens WHERE token=?`, token)
	return err
}
