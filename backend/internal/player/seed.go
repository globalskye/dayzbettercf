package player

import (
	"log"
	"time"
)

const sampleCftoolsID = "sample-example-1"
const sampleDisplayName = "ExamplePlayer"

// SeedSample добавляет одного примерного игрока с историей онлайна для демонстрации.
// Вызывать при SEED_SAMPLE=1. Если игрок уже есть — не дублирует.
func (r *Repository) SeedSample() error {
	existing, _ := r.GetByCftoolsID(sampleCftoolsID)
	if existing != nil {
		log.Printf("SeedSample: игрок %s уже есть, пропуск", sampleCftoolsID)
		return nil
	}

	res, err := r.db.Exec(`
		INSERT INTO players (cftools_id, display_name, updated_at)
		VALUES (?, ?, datetime('now'))
	`, sampleCftoolsID, sampleDisplayName)
	if err != nil {
		return err
	}
	playerID, _ := res.LastInsertId()
	log.Printf("SeedSample: создан игрок id=%d %s", playerID, sampleDisplayName)

	now := time.Now().UTC()
	// Записываем историю: пары "зашёл онлайн" → "вышел оффлайн" с сервером и длительностью сессии.
	type session struct {
		server   string
		offAt    time.Time
		duration int64
	}
	sessions := []session{
		{"DayZ RU #1 | 1PP", now.Add(-2 * 24 * time.Hour).Add(-3*time.Hour), 2*3600 + 15*60},
		{"DE 1.25 [EXP] Official", now.Add(-2 * 24 * time.Hour).Add(-6*time.Hour), 3600 + 45*60},
		{"DayZ RU #1 | 1PP", now.Add(-3 * 24 * time.Hour).Add(-2*time.Hour), 4 * 3600},
		{"US West Coast | Vanilla", now.Add(-4 * 24 * time.Hour).Add(-5*time.Hour), 90 * 60},
		{"DayZ RU #1 | 1PP", now.Add(-4 * 24 * time.Hour).Add(-8*time.Hour), 2*3600 + 30*60},
		{"DE 1.25 [EXP] Official", now.Add(-5 * 24 * time.Hour).Add(-1*time.Hour), 3*3600 + 20*60},
		{"DayZ RU #1 | 1PP", now.Add(-6 * 24 * time.Hour).Add(-4*time.Hour), 3600},
		{"DayZ RU #2 | 3PP", now.Add(-7 * 24 * time.Hour).Add(-3*time.Hour), 5*3600 + 15*60},
		{"DayZ RU #1 | 1PP", now.Add(-8 * 24 * time.Hour).Add(-2*time.Hour), 2 * 3600},
		{"DE 1.25 [EXP] Official", now.Add(-9 * 24 * time.Hour).Add(-6*time.Hour), 4*3600 + 45*60},
		{"DayZ RU #1 | 1PP", now.Add(-10 * 24 * time.Hour).Add(-2*time.Hour), 1*3600 + 10*60},
		{"US West Coast | Vanilla", now.Add(-11 * 24 * time.Hour).Add(-5*time.Hour), 2*3600 + 30*60},
		{"DayZ RU #1 | 1PP", now.Add(-12 * 24 * time.Hour).Add(-1*time.Hour), 3 * 3600},
		{"DE 1.25 [EXP] Official", now.Add(-14 * 24 * time.Hour).Add(-4*time.Hour), 90 * 60},
		{"DayZ RU #2 | 3PP", now.Add(-15 * 24 * time.Hour).Add(-2*time.Hour), 2*3600 + 15*60},
	}

	for _, s := range sessions {
		startAt := s.offAt.Add(-time.Duration(s.duration) * time.Second)
		tsOff := s.offAt.Format(time.RFC3339)
		tsOn := startAt.Format(time.RFC3339)
		_, err = r.db.Exec(`
			INSERT INTO player_history (player_id, ts, online, server_name, playtime_sec, sessions_count, display_name, session_duration_sec, offline_duration_sec)
			VALUES (?, ?, 1, ?, 0, 0, ?, 0, 0)
		`, playerID, tsOn, s.server, sampleDisplayName)
		if err != nil {
			return err
		}
		_, err = r.db.Exec(`
			INSERT INTO player_history (player_id, ts, online, server_name, playtime_sec, sessions_count, display_name, session_duration_sec, offline_duration_sec)
			VALUES (?, ?, 0, ?, 0, 0, ?, ?, 0)
		`, playerID, tsOff, s.server, sampleDisplayName, s.duration)
		if err != nil {
			return err
		}
	}
	log.Printf("SeedSample: добавлено записей истории: %d сессий", len(sessions))

	_, err = r.db.Exec(`INSERT OR IGNORE INTO tracked_players (player_id, added_at) VALUES (?, datetime('now'))`, playerID)
	if err != nil {
		return err
	}
	log.Printf("SeedSample: игрок добавлен в отслеживание")
	return nil
}
