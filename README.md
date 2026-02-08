# DayZ Smart CF

Полный контроль над игроками DayZ через CFtools: поиск, база, профили.

## Docker (рекомендуется)

```bash
# 1. Создай .env с CFtools cookies и логином админки
cp .env.docker.example .env
# Отредактируй .env: CFTOOLS_CDN_AUTH (и при необходимости SESSION, USER_INFO).
# Добавь ADMIN_USER=admin и ADMIN_PASS=admin — иначе при логине будет 502 (бэкенд не создаст админа).

# 2. Запуск
docker compose up -d

# 3. Открой http://localhost:8888
```

**Если на сервере уже есть сайт на 80:** приложение слушает свой порт (по умолчанию **8888**). Заходи по IP: `http://IP_СЕРВЕРА:8888`. Бэкенд наружу не открыт — все запросы к API идут через nginx во фронте. Сменить порт: в `docker-compose.yml` у сервиса `frontend` замени `"8888:80"` на нужный, например `"3000:80"`.

## Локальный запуск

### Backend (порт 8080)

```bash
cd backend
go run ./cmd/server
```

Переменные окружения (опционально):
- `PORT` — порт сервера (по умолчанию 8080)
- `ENV` — окружение (development/production)

### Frontend (порт 5173)

```bash
cd frontend
npm run dev
```

Откройте http://localhost:5173

### Одновременный запуск

1. В первом терминале: `cd backend && go run ./cmd/server`
2. Во втором терминале: `cd frontend && npm run dev`

## Авторизация

При первом запуске создаётся admin-пользователь, если заданы `ADMIN_USER` и `ADMIN_PASS` в `.env`. Роли: `admin`, `editor`, `viewer`. Страница настроек — только для admin.

**502 Bad Gateway при нажатии «Логин»:** в корневом `.env` (для Docker) должны быть `ADMIN_USER` и `ADMIN_PASS`. Если их нет, бэкенд не создаёт пользователя и прокси может отдавать 502. Добавь в `.env` строки `ADMIN_USER=admin` и `ADMIN_PASS=admin`, затем `docker compose up -d --force-recreate backend`. Если 502 остаётся — смотри логи: `docker compose logs backend`.

## API

- `GET /health` — проверка статуса (публично)
- `POST /api/v1/auth/login` — вход (username, password)
- `GET /api/v1/auth/me` — текущий пользователь (Bearer token)
- `GET /api/v1/players` — список игроков в БД
- `GET /api/v1/players/search?q=ник` — поиск по базе (локально)
- `GET /api/v1/players/cftools-search?q=ник` — поиск в CFtools API (только ответ, без сохранения)
- `POST /api/v1/players/sync-batch` — синхронизировать выбранных в базу (body: `{cftools_ids: [...]}`)
- `GET /api/v1/players/:id` — игрок по ID
- `POST /api/v1/players/:id/sync` — обновить данные игрока из CFtools

## CFtools

**Режим 1 — токен из браузера (рекомендуется):**
1. Залогинься на auth.cftools.cloud в браузере
2. DevTools → Application → Cookies → скопируй значения
3. В `.env`: `CFTOOLS_CDN_AUTH`, `CFTOOLS_SESSION`, `CFTOOLS_USER_INFO`, `CFTOOLS_CF_CLEARANCE`

**Режим 2 — автологин** (только локально, не в Docker):
- `CFTOOLS_IDENTIFIER` + `CFTOOLS_PASSWORD_HASH` (SHA256)
- Headless-браузер проходит Cloudflare (~15 сек)

## Технологии

**Backend:**
- Go 1.24
- Chi router
- CORS для фронтенда

**Frontend:**
- React 19
- TypeScript
- Vite 7
# dayzbettercf
