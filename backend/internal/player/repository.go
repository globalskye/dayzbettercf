package player

import (
	"database/sql"
	"time"
)

type Player struct {
	ID                  int64      `json:"id"`
	CftoolsID           string     `json:"cftools_id"`
	DisplayName         string     `json:"display_name"`
	Avatar              string     `json:"avatar,omitempty"`
	IsBot               bool       `json:"is_bot"`
	AccountStatus       int        `json:"account_status"`
	PlaytimeSec         int64      `json:"playtime_sec"`
	SessionsCount       int        `json:"sessions_count"`
	BansCount           int        `json:"bans_count"`
	LinkedAccountsCount int        `json:"linked_accounts_count"`
	LastActivityAt      *time.Time `json:"last_activity_at,omitempty"`
	LastSeenAt          *time.Time `json:"last_seen_at,omitempty"`
	Online              bool       `json:"online"`
	RawStatus           string     `json:"raw_status,omitempty"`
	RawOverview         string     `json:"raw_overview,omitempty"`
	RawStructure        string     `json:"raw_structure,omitempty"`
	RawPlayState        string     `json:"raw_play_state,omitempty"`
	RawBans             string     `json:"raw_bans,omitempty"`
	RawBattlEye         string     `json:"raw_battleye,omitempty"`
	Steam64             string     `json:"steam64,omitempty"`
	SteamAvatar         string     `json:"steam_avatar,omitempty"`
	SteamPersona        string     `json:"steam_persona,omitempty"`
	SteamVacBans        int        `json:"steam_vac_bans,omitempty"`
	SteamGameBans       int        `json:"steam_game_bans,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	Nicknames           []string   `json:"nicknames,omitempty"`
	LinkedCftoolsIDs    []string   `json:"linked_cftools_ids,omitempty"`
	ServerIDs           []string   `json:"server_ids,omitempty"`
	LastServerIdentifier string    `json:"last_server_identifier,omitempty"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) UpsertPlayer(p *Player) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	res, err := r.db.Exec(`
		INSERT INTO players (cftools_id, display_name, avatar, is_bot, account_status, playtime_sec, sessions_count, bans_count, linked_accounts_count, last_activity_at, last_seen_at, online, raw_status, raw_overview, raw_structure, raw_play_state, raw_bans, raw_battleye, steam64, steam_avatar, steam_persona, steam_vac_bans, steam_game_bans, last_server_identifier, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(cftools_id) DO UPDATE SET
			display_name = excluded.display_name,
			avatar = COALESCE(NULLIF(excluded.avatar,''), avatar),
			is_bot = excluded.is_bot,
			account_status = excluded.account_status,
			playtime_sec = excluded.playtime_sec,
			sessions_count = excluded.sessions_count,
			bans_count = excluded.bans_count,
			linked_accounts_count = excluded.linked_accounts_count,
			last_activity_at = COALESCE(excluded.last_activity_at, last_activity_at),
			last_seen_at = excluded.last_seen_at,
			online = excluded.online,
			raw_status = excluded.raw_status,
			raw_overview = excluded.raw_overview,
			raw_structure = excluded.raw_structure,
			raw_play_state = excluded.raw_play_state,
			raw_bans = COALESCE(NULLIF(excluded.raw_bans,''), raw_bans),
			raw_battleye = COALESCE(NULLIF(excluded.raw_battleye,''), raw_battleye),
			steam64 = COALESCE(NULLIF(excluded.steam64,''), steam64),
			steam_avatar = COALESCE(NULLIF(excluded.steam_avatar,''), steam_avatar),
			steam_persona = COALESCE(NULLIF(excluded.steam_persona,''), steam_persona),
			steam_vac_bans = CASE WHEN excluded.steam_vac_bans > 0 OR excluded.steam_game_bans > 0 THEN excluded.steam_vac_bans ELSE steam_vac_bans END,
			steam_game_bans = CASE WHEN excluded.steam_vac_bans > 0 OR excluded.steam_game_bans > 0 THEN excluded.steam_game_bans ELSE steam_game_bans END,
			last_server_identifier = COALESCE(NULLIF(excluded.last_server_identifier,''), last_server_identifier),
			updated_at = excluded.updated_at
	`,
		p.CftoolsID, p.DisplayName, p.Avatar, boolToInt(p.IsBot), p.AccountStatus, p.PlaytimeSec, p.SessionsCount, p.BansCount, p.LinkedAccountsCount,
		timePtrToStr(p.LastActivityAt), timePtrToStr(p.LastSeenAt), boolToInt(p.Online),
		p.RawStatus, p.RawOverview, p.RawStructure, p.RawPlayState, p.RawBans, p.RawBattlEye,
		p.Steam64, p.SteamAvatar, p.SteamPersona, p.SteamVacBans, p.SteamGameBans, p.LastServerIdentifier, now,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	if id == 0 {
		var existingID int64
		_ = r.db.QueryRow("SELECT id FROM players WHERE cftools_id = ?", p.CftoolsID).Scan(&existingID)
		return existingID, nil
	}
	return id, nil
}

func (r *Repository) UpdatePlayerDisplayName(playerID int64, displayName string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.Exec(`UPDATE players SET display_name = ?, updated_at = ? WHERE id = ?`, displayName, now, playerID)
	return err
}

func (r *Repository) UpdatePlayerOnlineStatus(playerID int64, online bool, serverName string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if online && serverName != "" {
		_, err := r.db.Exec(`UPDATE players SET online = ?, last_seen_at = ?, last_server_identifier = ?, updated_at = ? WHERE id = ?`,
			boolToInt(online), now, serverName, now, playerID)
		return err
	}
	if !online {
		_, err := r.db.Exec(`UPDATE players SET online = ?, last_seen_at = ?, updated_at = ? WHERE id = ?`,
			boolToInt(online), now, now, playerID)
		return err
	}
	_, err := r.db.Exec(`UPDATE players SET online = ?, last_seen_at = ?, updated_at = ? WHERE id = ?`,
		boolToInt(online), now, now, playerID)
	return err
}

func (r *Repository) UpsertNickname(playerID int64, nickname, source string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.Exec(`
		INSERT INTO nicknames (player_id, nickname, source, last_seen_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(player_id, nickname) DO UPDATE SET last_seen_at = excluded.last_seen_at
	`, playerID, nickname, source, now)
	return err
}

func (r *Repository) UpsertPlayerLink(playerID int64, linkedCftoolsID string, confirmed, trusted bool) error {
	_, err := r.db.Exec(`
		INSERT INTO player_links (player_id, linked_cftools_id, confirmed, trusted)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(player_id, linked_cftools_id) DO UPDATE SET confirmed = excluded.confirmed, trusted = excluded.trusted
	`, playerID, linkedCftoolsID, boolToInt(confirmed), boolToInt(trusted))
	return err
}

func (r *Repository) DeletePlayerLinks(playerID int64) error {
	_, err := r.db.Exec("DELETE FROM player_links WHERE player_id = ?", playerID)
	return err
}

func (r *Repository) UpsertPlayerServer(playerID int64, serverID, identifier string, gameType int) error {
	_, err := r.db.Exec(`
		INSERT INTO player_servers (player_id, cftools_server_id, identifier, game_type)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(player_id, cftools_server_id) DO UPDATE SET identifier = excluded.identifier, game_type = excluded.game_type
	`, playerID, serverID, identifier, gameType)
	return err
}

func (r *Repository) DeletePlayerServers(playerID int64) error {
	_, err := r.db.Exec("DELETE FROM player_servers WHERE player_id = ?", playerID)
	return err
}

func (r *Repository) GetByCftoolsID(cftoolsID string) (*Player, error) {
	var p Player
	var avatar, rawStatus, rawOverview, rawStructure, rawPlayState, rawBans, rawBattleye sql.NullString
	var steam64, steamAvatar, steamPersona sql.NullString
	var lastActivityAt, lastSeenAt sql.NullString
	var createdAt, updatedAt string
	var lastServer string
	err := r.db.QueryRow(`
		SELECT id, cftools_id, display_name, avatar, is_bot, account_status, playtime_sec, sessions_count, bans_count, linked_accounts_count,
		       last_activity_at, last_seen_at, online, raw_status, raw_overview, raw_structure, raw_play_state,
		       raw_bans, raw_battleye, steam64, steam_avatar, steam_persona, steam_vac_bans, steam_game_bans,
		       COALESCE(last_server_identifier, ''), created_at, updated_at
		FROM players WHERE cftools_id = ?
	`, cftoolsID).Scan(
		&p.ID, &p.CftoolsID, &p.DisplayName, &avatar, &p.IsBot, &p.AccountStatus, &p.PlaytimeSec, &p.SessionsCount, &p.BansCount, &p.LinkedAccountsCount,
		&lastActivityAt, &lastSeenAt, &p.Online, &rawStatus, &rawOverview, &rawStructure, &rawPlayState,
		&rawBans, &rawBattleye, &steam64, &steamAvatar, &steamPersona, &p.SteamVacBans, &p.SteamGameBans,
		&lastServer, &createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	p.CreatedAt = parseTimeValue(createdAt)
	p.UpdatedAt = parseTimeValue(updatedAt)
	p.Avatar = avatar.String
	p.RawStatus = rawStatus.String
	p.RawOverview = rawOverview.String
	p.RawStructure = rawStructure.String
	p.RawPlayState = rawPlayState.String
	p.RawBans = rawBans.String
	p.RawBattlEye = rawBattleye.String
	p.Steam64 = steam64.String
	p.SteamAvatar = steamAvatar.String
	p.SteamPersona = steamPersona.String
	p.LastActivityAt = parseTime(lastActivityAt.String)
	p.LastSeenAt = parseTime(lastSeenAt.String)
	if lastServer != "" {
		p.LastServerIdentifier = lastServer
	}

	rows, _ := r.db.Query("SELECT nickname FROM nicknames WHERE player_id = ?", p.ID)
	for rows.Next() {
		var n string
		_ = rows.Scan(&n)
		p.Nicknames = append(p.Nicknames, n)
	}
	rows.Close()

	rows, _ = r.db.Query("SELECT linked_cftools_id FROM player_links WHERE player_id = ?", p.ID)
	for rows.Next() {
		var id string
		_ = rows.Scan(&id)
		p.LinkedCftoolsIDs = append(p.LinkedCftoolsIDs, id)
	}
	rows.Close()

	rows, _ = r.db.Query("SELECT cftools_server_id FROM player_servers WHERE player_id = ?", p.ID)
	for rows.Next() {
		var id string
		_ = rows.Scan(&id)
		p.ServerIDs = append(p.ServerIDs, id)
	}
	rows.Close()

	return &p, nil
}

type ListOptions struct {
	Limit      int
	Offset     int
	OnlyOnline bool
	OnlyBanned bool
	Sort       string // "online", "updated", "playtime", "bans"
}

func (r *Repository) ListAll(opts ListOptions) ([]*Player, error) {
	limit, offset := opts.Limit, opts.Offset
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	order := "ORDER BY updated_at DESC"
	switch opts.Sort {
	case "playtime":
		order = "ORDER BY playtime_sec DESC, updated_at DESC"
	case "bans":
		order = "ORDER BY bans_count DESC, playtime_sec DESC, updated_at DESC"
	case "online":
		order = "ORDER BY online DESC, COALESCE(last_seen_at,'') DESC, updated_at DESC"
	default:
		order = "ORDER BY updated_at DESC"
	}
	where := "1=1"
	if opts.OnlyOnline {
		where += " AND online = 1"
	}
	if opts.OnlyBanned {
		where += " AND bans_count > 0"
	}
	query := `
		SELECT id, cftools_id, display_name, avatar, is_bot, account_status, playtime_sec, sessions_count, bans_count, linked_accounts_count,
		       last_activity_at, last_seen_at, online, COALESCE(last_server_identifier,''), created_at, updated_at
		FROM players WHERE ` + where + " " + order + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Player
	for rows.Next() {
		var p Player
		var avatar sql.NullString
		var lastActivityAt, lastSeenAt sql.NullString
		var lastServer string
		var createdAt, updatedAt string
		_ = rows.Scan(&p.ID, &p.CftoolsID, &p.DisplayName, &avatar, &p.IsBot, &p.AccountStatus, &p.PlaytimeSec, &p.SessionsCount, &p.BansCount, &p.LinkedAccountsCount,
			&lastActivityAt, &lastSeenAt, &p.Online, &lastServer, &createdAt, &updatedAt)
		p.Avatar = avatar.String
		p.LastActivityAt = parseTime(lastActivityAt.String)
		p.LastSeenAt = parseTime(lastSeenAt.String)
		p.LastServerIdentifier = lastServer
		p.CreatedAt = parseTimeValue(createdAt)
		p.UpdatedAt = parseTimeValue(updatedAt)
		list = append(list, &p)
	}
	return list, nil
}

func (r *Repository) SearchByNickname(q string, limit int, opts *ListOptions) ([]*Player, error) {
	if limit <= 0 {
		limit = 5000
	}
	if limit > 10000 {
		limit = 10000
	}
	order := "ORDER BY p.updated_at DESC"
	if opts != nil {
		switch opts.Sort {
		case "playtime":
			order = "ORDER BY p.playtime_sec DESC, p.updated_at DESC"
		case "bans":
			order = "ORDER BY p.bans_count DESC, p.playtime_sec DESC, p.updated_at DESC"
		case "online":
			order = "ORDER BY p.online DESC, COALESCE(p.last_seen_at,'') DESC, p.updated_at DESC"
		}
	}
	where := "(LOWER(p.display_name) LIKE LOWER(?) OR LOWER(n.nickname) LIKE LOWER(?))"
	if opts != nil && opts.OnlyOnline {
		where += " AND p.online = 1"
	}
	if opts != nil && opts.OnlyBanned {
		where += " AND p.bans_count > 0"
	}
	rows, err := r.db.Query(`
		SELECT DISTINCT p.id, p.cftools_id, p.display_name, p.avatar, p.is_bot, p.account_status, p.playtime_sec, p.sessions_count, p.bans_count, p.linked_accounts_count,
		       p.last_activity_at, p.last_seen_at, p.online, COALESCE(p.last_server_identifier,''), p.created_at, p.updated_at
		FROM players p
		LEFT JOIN nicknames n ON n.player_id = p.id
		WHERE `+where+`
		`+order+` LIMIT ?
	`, "%"+q+"%", "%"+q+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Player
	for rows.Next() {
		var p Player
		var avatar sql.NullString
		var lastActivityAt, lastSeenAt sql.NullString
		var lastServer string
		var createdAt, updatedAt string
		_ = rows.Scan(&p.ID, &p.CftoolsID, &p.DisplayName, &avatar, &p.IsBot, &p.AccountStatus, &p.PlaytimeSec, &p.SessionsCount, &p.BansCount, &p.LinkedAccountsCount,
			&lastActivityAt, &lastSeenAt, &p.Online, &lastServer, &createdAt, &updatedAt)
		p.Avatar = avatar.String
		p.LastActivityAt = parseTime(lastActivityAt.String)
		p.LastSeenAt = parseTime(lastSeenAt.String)
		p.LastServerIdentifier = lastServer
		p.CreatedAt = parseTimeValue(createdAt)
		p.UpdatedAt = parseTimeValue(updatedAt)
		list = append(list, &p)
	}
	return list, nil
}

func (r *Repository) LogSync(playerID int64, cftoolsID, displayName string) error {
	_, err := r.db.Exec(`INSERT INTO sync_log (player_id, cftools_id, display_name) VALUES (?, ?, ?)`,
		playerID, cftoolsID, displayName)
	return err
}

func (r *Repository) Count(opts *ListOptions) (int, error) {
	where := "1=1"
	if opts != nil {
		if opts.OnlyOnline {
			where += " AND online = 1"
		}
		if opts.OnlyBanned {
			where += " AND bans_count > 0"
		}
	}
	var n int
	err := r.db.QueryRow("SELECT COUNT(*) FROM players WHERE " + where).Scan(&n)
	return n, err
}

// WipeAllData удаляет все данные приложения (игроки, группы, история, отслеживание). Таблица users не трогается.
func (r *Repository) WipeAllData() error {
	order := []string{
		"group_members", "groups", "tracked_players", "player_history", "sync_log",
		"nicknames", "player_links", "bans", "player_servers", "players",
	}
	for _, table := range order {
		if _, err := r.db.Exec("DELETE FROM " + table); err != nil {
			return err
		}
	}
	// Сброс автоинкремента
	_, _ = r.db.Exec("DELETE FROM sqlite_sequence WHERE name IN ('players','groups','group_members','player_history','tracked_players','sync_log','nicknames','player_links','bans','player_servers')")
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func timePtrToStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t := parseTimeValue(s)
	if t.IsZero() {
		return nil
	}
	return &t
}

// parseTimeValue парсит строку из SQLite TEXT в time.Time (driver возвращает string)
func parseTimeValue(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05.999999",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
