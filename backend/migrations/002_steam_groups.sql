-- Steam: данные Steam профиля
ALTER TABLE players ADD COLUMN steam64 TEXT;
ALTER TABLE players ADD COLUMN steam_avatar TEXT;
ALTER TABLE players ADD COLUMN steam_persona TEXT;
ALTER TABLE players ADD COLUMN steam_vac_bans INTEGER DEFAULT 0;
ALTER TABLE players ADD COLUMN steam_game_bans INTEGER DEFAULT 0;
ALTER TABLE players ADD COLUMN raw_bans TEXT;
ALTER TABLE players ADD COLUMN raw_battleye TEXT;

-- Groups: группы игроков (как в старом проекте)
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    alias TEXT DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(group_id, player_id)
);

CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);
CREATE INDEX IF NOT EXISTS idx_group_members_player_id ON group_members(player_id);
