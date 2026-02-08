-- Длительность сессии (при уходе оффлайн) и оффлайн (при возврате онлайн)
ALTER TABLE player_history ADD COLUMN session_duration_sec INTEGER DEFAULT 0;
ALTER TABLE player_history ADD COLUMN offline_duration_sec INTEGER DEFAULT 0;
