package store

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/gurkanfikretgunak/synclink/backend/internal/auth"
)

const schema = `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	role TEXT NOT NULL,
	status TEXT NOT NULL,
	created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS pages (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL UNIQUE,
	slug TEXT NOT NULL UNIQUE,
	display_name TEXT NOT NULL,
	bio TEXT NOT NULL DEFAULT '',
	avatar_url TEXT,
	theme TEXT NOT NULL DEFAULT 'default',
	avatar_shape TEXT NOT NULL DEFAULT 'circle',
	accent_color TEXT NOT NULL DEFAULT '#111111',
	background TEXT NOT NULL DEFAULT 'cream',
	motion TEXT NOT NULL DEFAULT 'subtle',
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS links (
	id TEXT PRIMARY KEY,
	page_id TEXT NOT NULL,
	title TEXT NOT NULL,
	url TEXT NOT NULL,
	icon TEXT,
	sort_order INTEGER NOT NULL DEFAULT 0,
	active INTEGER NOT NULL DEFAULT 1,
	clicks INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_links_page_order ON links(page_id, sort_order);

CREATE TABLE IF NOT EXISTS settings (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	payload TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
	token TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	expires_at TEXT NOT NULL
);
`

func Migrate(db *sql.DB) error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	if err := ensureLinksClicks(db); err != nil {
		return err
	}
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM settings`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	raw, err := json.Marshal(auth.DefaultSettings())
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO settings (id, payload) VALUES (1, ?)`, string(raw))
	return err
}

func ensureLinksClicks(db *sql.DB) error {
	var name string
	err := db.QueryRow(`SELECT name FROM pragma_table_info('links') WHERE name='clicks'`).Scan(&name)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	_, err = db.Exec(`ALTER TABLE links ADD COLUMN clicks INTEGER NOT NULL DEFAULT 0`)
	return err
}
