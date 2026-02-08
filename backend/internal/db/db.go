package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func Open(dbURL string) (*sql.DB, error) {
	if len(dbURL) > 5 && dbURL[:5] == "file:" {
		path := dbURL[5:]
		for i, c := range path {
			if c == '?' || c == '&' {
				path = path[:i]
				break
			}
		}
		if dir := filepath.Dir(path); dir != "." {
			_ = os.MkdirAll(dir, 0755)
		}
	}

	db, err := sql.Open("sqlite", dbURL)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB, migrationsDir string) error {
	migration := `-- Players: основная таблица игроков
CREATE TABLE IF NOT EXISTS players (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cftools_id TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    avatar TEXT,
    is_bot INTEGER DEFAULT 0,
    account_status INTEGER DEFAULT 0,
    playtime_sec INTEGER DEFAULT 0,
    sessions_count INTEGER DEFAULT 0,
    bans_count INTEGER DEFAULT 0,
    linked_accounts_count INTEGER DEFAULT 0,
    last_activity_at TEXT,
    last_seen_at TEXT,
    online INTEGER DEFAULT 0,
    raw_status TEXT,
    raw_overview TEXT,
    raw_structure TEXT,
    raw_play_state TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_players_cftools_id ON players(cftools_id);
CREATE INDEX IF NOT EXISTS idx_players_display_name ON players(display_name);

CREATE TABLE IF NOT EXISTS nicknames (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    nickname TEXT NOT NULL,
    source TEXT DEFAULT 'display_name',
    first_seen_at TEXT NOT NULL DEFAULT (datetime('now')),
    last_seen_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(player_id, nickname)
);
CREATE INDEX IF NOT EXISTS idx_nicknames_nickname ON nicknames(nickname);
CREATE INDEX IF NOT EXISTS idx_nicknames_player_id ON nicknames(player_id);

CREATE TABLE IF NOT EXISTS player_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    linked_cftools_id TEXT NOT NULL,
    confirmed INTEGER DEFAULT 0,
    trusted INTEGER DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(player_id, linked_cftools_id)
);
CREATE INDEX IF NOT EXISTS idx_player_links_player_id ON player_links(player_id);

CREATE TABLE IF NOT EXISTS bans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    server_id TEXT,
    server_name TEXT,
    reason TEXT,
    banned_at TEXT,
    expires_at TEXT,
    banned_by TEXT,
    source TEXT DEFAULT 'cftools',
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_bans_player_id ON bans(player_id);

CREATE TABLE IF NOT EXISTS player_servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    cftools_server_id TEXT NOT NULL,
    identifier TEXT,
    game_type INTEGER DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(player_id, cftools_server_id)
);
CREATE INDEX IF NOT EXISTS idx_player_servers_player_id ON player_servers(player_id);
`
	if _, err := db.Exec(migration); err != nil {
		return err
	}
	// Миграция 002: steam + groups (встроена, не зависит от cwd)
	migration002 := `
ALTER TABLE players ADD COLUMN steam64 TEXT;
ALTER TABLE players ADD COLUMN steam_avatar TEXT;
ALTER TABLE players ADD COLUMN steam_persona TEXT;
ALTER TABLE players ADD COLUMN steam_vac_bans INTEGER DEFAULT 0;
ALTER TABLE players ADD COLUMN steam_game_bans INTEGER DEFAULT 0;
ALTER TABLE players ADD COLUMN raw_bans TEXT;
ALTER TABLE players ADD COLUMN raw_battleye TEXT;
CREATE TABLE IF NOT EXISTS groups (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, created_at TEXT NOT NULL DEFAULT (datetime('now')), updated_at TEXT NOT NULL DEFAULT (datetime('now')));
CREATE TABLE IF NOT EXISTS group_members (id INTEGER PRIMARY KEY AUTOINCREMENT, group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE, player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE, alias TEXT DEFAULT '', created_at TEXT NOT NULL DEFAULT (datetime('now')), UNIQUE(group_id, player_id));
CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);
CREATE INDEX IF NOT EXISTS idx_group_members_player_id ON group_members(player_id);
`
	for _, stmt := range splitStatements(migration002) {
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			if strings.Contains(err.Error(), "duplicate column name") || strings.Contains(err.Error(), "already exists") {
				log.Printf("Migration 002: skip (already applied): %s", truncate(stmt, 60))
				continue
			}
			return fmt.Errorf("migration 002: %w", err)
		}
	}
	log.Println("Migration 002: ok")

	// Миграция 003: player_history + tracked_players
	migration003 := `
CREATE TABLE IF NOT EXISTS player_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    ts TEXT NOT NULL DEFAULT (datetime('now')),
    online INTEGER DEFAULT 0,
    playtime_sec INTEGER DEFAULT 0,
    sessions_count INTEGER DEFAULT 0,
    display_name TEXT
);
CREATE INDEX IF NOT EXISTS idx_player_history_player_id ON player_history(player_id);
CREATE INDEX IF NOT EXISTS idx_player_history_ts ON player_history(ts);

ALTER TABLE players ADD COLUMN last_server_identifier TEXT;
ALTER TABLE player_history ADD COLUMN server_name TEXT;
CREATE TABLE IF NOT EXISTS tracked_players (
    player_id INTEGER NOT NULL PRIMARY KEY REFERENCES players(id) ON DELETE CASCADE,
    added_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`
	for _, stmt := range splitStatements(migration003) {
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "duplicate column") {
				log.Printf("Migration 003: skip (already applied): %s", truncate(stmt, 60))
				continue
			}
			return fmt.Errorf("migration 003: %w", err)
		}
	}
	log.Println("Migration 003: ok")

	return runMigrationsFromDir(db, migrationsDir)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func runMigrationsFromDir(db *sql.DB, dir string) error {
	if dir == "" {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") || e.Name() == "001_init.sql" {
			continue
		}
		path := filepath.Join(dir, e.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		for _, stmt := range splitStatements(string(content)) {
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {
				if strings.Contains(err.Error(), "duplicate column name") || strings.Contains(err.Error(), "already exists") {
					continue
				}
				return fmt.Errorf("%s: %w", e.Name(), err)
			}
		}
		log.Printf("Migration %s: ok", e.Name())
	}
	return nil
}

func splitStatements(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ";") {
		stmt := strings.TrimSpace(stripLeadingCommentLines(part))
		if stmt != "" {
			out = append(out, stmt+";")
		}
	}
	return out
}

// stripLeadingCommentLines removes lines that are only SQL comments from the start.
func stripLeadingCommentLines(s string) string {
	lines := strings.Split(s, "\n")
	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line != "" && !strings.HasPrefix(line, "--") {
			break
		}
		i++
	}
	return strings.Join(lines[i:], "\n")
}

