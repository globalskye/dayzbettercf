# DayZ Smart CF — Roadmap

На основе старого проекта (`old/DayzProfileFront`, `old/DayzProfilesBack`) и текущего dayzsmartcf.

## Логика из старого проекта

- **Поиск по identifier** → global-query → список профилей → добавить в группу
- **Группы** — создание, добавление/удаление игроков, refresh ников
- **Профиль** — overview (aliases, playtime, sessions), steam (avatar, VAC/EAC), bans (CFtools + BattlEye), alternate accounts
- **Авторизация** — Bearer токен (сейчас используем cookies cdn-auth)

## Что уже есть (dayzsmartcf)

- ✅ Поиск по нику (global-query)
- ✅ База игроков SQLite (players, nicknames, player_links, bans, player_servers)
- ✅ Профиль: статус, playState, overview, structure
- ✅ Обновление auth через UI (Настройки)
- ✅ Темная тема

## Что добавить

### Backend

1. **Steam** — endpoint `/app/v1/profile/{id}/steam`, сохранять steam64, avatar, VAC/Game bans
2. **Bans** — endpoint `/app/v1/profile/{id}/bans`, парсить и сохранять в таблицу bans
3. **BattlEye** — endpoint `/app/v1/profile/{id}/publisher-services/battleye/ban-status`
4. **Группы** — таблицы `groups`, `group_members`, CRUD API
5. **Расширить sync** — загружать steam, bans, BattlEye при sync
6. **Расширить схему** — steam64, steam_avatar, steam_vac_bans, steam_game_bans (или отдельная таблица steam_profiles)

### Frontend

1. **UI** — более современный и красивый (как в старом, но лучше)
2. **Страница групп** — список групп, создание, добавление/удаление игроков, refresh
3. **Профиль игрока** — блоки: Steam (аватар, ссылка, VAC/EAC), Bans, Alternate accounts, Nicknames
4. **Поиск** — добавить результат в группу одним кликом
5. **Настройки** — вынести auth в sidebar/header

### Приоритеты

1. Steam + Bans в sync и профиле
2. Группы (CRUD)
3. UI upgrade (группы, расширенный профиль)

---

## Планы на будущее

1. **Жёсткий поиск по банам**  
   Отдельный режим поиска: находить людей именно по банам (по причине, по серверу, по дате). Фильтры и строгий поиск по базе банов.

2. **Steam-аккаунты и парсинг друзей**  
   Поиск и привязка Steam-аккаунтов к игрокам, парсинг списков друзей, глубокая аналитика и вычисление связей между аккаунтами (альты, общие друзья, граф связей).
