package player

import (
	"database/sql"
	"time"
)

const maxTracked = 10

type HistoryRecord struct {
	Ts                 string `json:"ts"`
	Online             bool   `json:"online"`
	ServerName         string `json:"server_name,omitempty"`
	PlaytimeSec        int64  `json:"playtime_sec"`
	SessionsCount      int    `json:"sessions_count"`
	DisplayName        string `json:"display_name,omitempty"`
	SessionDurationSec int64  `json:"session_duration_sec,omitempty"` // при уходе оффлайн — длительность сессии
	OfflineDurationSec int64  `json:"offline_duration_sec,omitempty"` // при возврате онлайн — сколько был оффлайн
}

func (r *Repository) GetLastPlayerHistory(playerID int64) (*HistoryRecord, error) {
	var h HistoryRecord
	var onlineInt int
	err := r.db.QueryRow(`SELECT ts, online, COALESCE(server_name,''), playtime_sec, sessions_count, COALESCE(display_name,''), COALESCE(session_duration_sec,0), COALESCE(offline_duration_sec,0) FROM player_history WHERE player_id = ? ORDER BY ts DESC LIMIT 1`,
		playerID).Scan(&h.Ts, &onlineInt, &h.ServerName, &h.PlaytimeSec, &h.SessionsCount, &h.DisplayName, &h.SessionDurationSec, &h.OfflineDurationSec)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	h.Online = onlineInt != 0
	return &h, nil
}

func (r *Repository) AppendPlayerHistory(playerID int64, online bool, serverName string, playtimeSec int64, sessionsCount int, displayName string, sessionDurationSec, offlineDurationSec int64) error {
	ts := time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.Exec(`INSERT INTO player_history (player_id, ts, online, server_name, playtime_sec, sessions_count, display_name, session_duration_sec, offline_duration_sec) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		playerID, ts, boolToInt(online), serverName, playtimeSec, sessionsCount, displayName, sessionDurationSec, offlineDurationSec)
	return err
}

func (r *Repository) GetPlayerHistory(playerID int64, limit int) ([]HistoryRecord, error) {
	if limit <= 0 {
		limit = 500
	}
	rows, err := r.db.Query(`SELECT ts, online, COALESCE(server_name,''), playtime_sec, sessions_count, COALESCE(display_name,''), COALESCE(session_duration_sec,0), COALESCE(offline_duration_sec,0) FROM player_history WHERE player_id = ? ORDER BY ts DESC LIMIT ?`,
		playerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []HistoryRecord
	for rows.Next() {
		var h HistoryRecord
		var onlineInt int
		_ = rows.Scan(&h.Ts, &onlineInt, &h.ServerName, &h.PlaytimeSec, &h.SessionsCount, &h.DisplayName, &h.SessionDurationSec, &h.OfflineDurationSec)
		h.Online = onlineInt != 0
		list = append(list, h)
	}
	return list, nil
}

func (r *Repository) AddTracked(playerID int64) error {
	var count int
	_ = r.db.QueryRow("SELECT COUNT(*) FROM tracked_players").Scan(&count)
	if count >= maxTracked {
		return ErrTrackedLimit
	}
	_, err := r.db.Exec(`INSERT OR IGNORE INTO tracked_players (player_id, added_at) VALUES (?, datetime('now'))`, playerID)
	return err
}

var ErrTrackedLimit = &trackedLimitError{}

type trackedLimitError struct{}

func (e *trackedLimitError) Error() string {
	return "tracked players limit reached (max 10)"
}

func (r *Repository) RemoveTracked(playerID int64) error {
	_, err := r.db.Exec("DELETE FROM tracked_players WHERE player_id = ?", playerID)
	return err
}

func (r *Repository) ListTracked(sort string) ([]*Player, error) {
	ids, err := r.ListTrackedCftoolsIDs()
	if err != nil {
		return nil, err
	}
	var list []*Player
	for _, cftoolsID := range ids {
		p, _ := r.GetByCftoolsID(cftoolsID)
		if p != nil {
			list = append(list, p)
		}
	}
	return list, nil
}

// ListTrackedCftoolsIDs возвращает cftools_id отслеживаемых в порядке added_at.
// Для списка отслеживаемых данные подтягиваются из CF в хендлере.
// Данные по игрокам потом подтягиваются из CF в хендлере.
func (r *Repository) ListTrackedCftoolsIDs() ([]string, error) {
	rows, err := r.db.Query(`SELECT p.cftools_id FROM players p JOIN tracked_players tp ON p.id = tp.player_id ORDER BY tp.added_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *Repository) IsTracked(playerID int64) (bool, error) {
	var n int
	err := r.db.QueryRow("SELECT 1 FROM tracked_players WHERE player_id = ?", playerID).Scan(&n)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}
