-- Players: основная таблица игроков
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

-- Nicknames: все ники для поиска
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

-- Player links: связанные аккаунты (альты)
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

-- Bans: баны (локальные + из CFtools когда будет endpoint)
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

-- Servers: серверы игрока
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
