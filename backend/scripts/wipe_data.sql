-- Очистка всех данных приложения (игроки, группы, история, отслеживание).
-- Пользователи (users) НЕ удаляются — логин остаётся.
-- Выполнить: sqlite3 path/to/dayzsmartcf.db < scripts/wipe_data.sql
-- Или в PowerShell: Get-Content backend\scripts\wipe_data.sql | sqlite3 backend\dayzsmartcf.db

PRAGMA foreign_keys = ON;

DELETE FROM group_members;
DELETE FROM groups;
DELETE FROM tracked_players;
DELETE FROM player_history;
DELETE FROM sync_log;
DELETE FROM nicknames;
DELETE FROM player_links;
DELETE FROM bans;
DELETE FROM player_servers;
DELETE FROM players;

-- Сброс автоинкремента (опционально, чтобы id снова начинались с 1)
DELETE FROM sqlite_sequence WHERE name IN ('players', 'groups', 'group_members', 'player_history', 'tracked_players', 'sync_log', 'nicknames', 'player_links', 'bans', 'player_servers');
