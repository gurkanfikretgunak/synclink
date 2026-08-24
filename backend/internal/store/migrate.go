package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

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
	socials TEXT NOT NULL DEFAULT '[]',
	verified INTEGER NOT NULL DEFAULT 0,
	page_password TEXT,
	published_at TEXT,
	cover_url TEXT,
	cover_kind TEXT NOT NULL DEFAULT '',
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
	last_clicked_at TEXT,
	featured INTEGER NOT NULL DEFAULT 0,
	thumbnail_url TEXT,
	starts_at TEXT,
	ends_at TEXT,
	sensitive INTEGER NOT NULL DEFAULT 0,
	section TEXT NOT NULL DEFAULT '',
	embed_url TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_links_page_order ON links(page_id, sort_order);

CREATE TABLE IF NOT EXISTS subscribers (
	id TEXT PRIMARY KEY,
	page_id TEXT NOT NULL,
	email TEXT NOT NULL,
	created_at TEXT NOT NULL,
	UNIQUE (page_id, email),
	FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_subscribers_page ON subscribers(page_id);

CREATE TABLE IF NOT EXISTS settings (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	payload TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
	token TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	expires_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS click_days (
	page_id TEXT NOT NULL,
	day TEXT NOT NULL,
	clicks INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (page_id, day),
	FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE
);
`

func Migrate(db *sql.DB) error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	alters := []struct {
		table  string
		column string
		ddl    string
	}{
		{"links", "clicks", "clicks INTEGER NOT NULL DEFAULT 0"},
		{"links", "last_clicked_at", "last_clicked_at TEXT"},
		{"links", "featured", "featured INTEGER NOT NULL DEFAULT 0"},
		{"links", "thumbnail_url", "thumbnail_url TEXT"},
		{"links", "starts_at", "starts_at TEXT"},
		{"links", "ends_at", "ends_at TEXT"},
		{"links", "sensitive", "sensitive INTEGER NOT NULL DEFAULT 0"},
		{"pages", "socials", "socials TEXT"},
		{"pages", "verified", "verified INTEGER NOT NULL DEFAULT 0"},
		{"pages", "page_password", "page_password TEXT"},
		{"pages", "published_at", "published_at TEXT"},
		{"pages", "cover_url", "cover_url TEXT"},
		{"pages", "cover_kind", "cover_kind TEXT NOT NULL DEFAULT ''"},
		{"links", "section", "section TEXT NOT NULL DEFAULT ''"},
		{"links", "embed_url", "embed_url TEXT"},
	}
	for _, a := range alters {
		if err := ensureColumn(db, a.table, a.column, a.ddl); err != nil {
			return err
		}
	}
	if _, err := db.Exec(`UPDATE pages SET published_at = created_at WHERE published_at IS NULL OR published_at = ''`); err != nil {
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

func ensureColumn(db *sql.DB, table, name, ddl string) error {
	var got string
	err := db.QueryRow(fmt.Sprintf(`SELECT name FROM pragma_table_info('%s') WHERE name=?`, table), name).Scan(&got)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	_, err = db.Exec(`ALTER TABLE ` + table + ` ADD COLUMN ` + ddl)
	return err
}
