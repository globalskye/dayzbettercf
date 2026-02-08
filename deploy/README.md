# Деплой на Linux

## Вариант 1: Docker (рекомендуется)

### Требования
- Docker и Docker Compose (v2)
- На Ubuntu/Debian: `sudo apt install docker.io docker-compose-plugin`

### Шаги

1. **Клонируй репозиторий** (или скопируй проект на сервер):
   ```bash
   git clone <url> dayzsmartcf
   cd dayzsmartcf
   ```

2. **Создай `.env`** из примера и заполни:
   ```bash
   cp .env.docker.example .env
   nano .env
   ```
   Обязательно укажи:
   - `CFTOOLS_CDN_AUTH`, `CFTOOLS_SESSION`, `CFTOOLS_USER_INFO`, `CFTOOLS_CF_CLEARANCE` (из cookies браузера на auth.cftools.cloud)
   - `ADMIN_USER` и `ADMIN_PASS` — логин/пароль админки
   - `JWT_SECRET` — смени на случайную строку в продакшене

3. **Запуск**:
   ```bash
   chmod +x deploy/deploy.sh deploy/stop.sh
   ./deploy/deploy.sh
   ```

4. Открой в браузере: **http://IP_СЕРВЕРА:3000** (порт задаётся в `docker-compose.yml`, по умолчанию 3000).

### Полезные команды
- Логи: `docker compose logs -f`
- Остановить: `./deploy/stop.sh`
- Перезапуск после изменений: `docker compose up -d --build`

### Автозапуск после перезагрузки
Контейнеры уже с `restart: unless-stopped` — поднимутся сами, если включён Docker. Убедись, что Docker стартует при загрузке:
```bash
sudo systemctl enable docker
sudo systemctl start docker
```

---

## Вариант 2: Без Docker (бинарник + nginx)

Если хочешь запускать бэкенд бинарником и раздавать фронт через nginx.

### Бэкенд
```bash
cd backend
go build -o server ./cmd/server
export PORT=8080
export DATABASE_URL=file:./dayzsmartcf.db
# Скопируй .env из корня или задай CFTOOLS_*, ADMIN_USER, ADMIN_PASS
source ../.env 2>/dev/null || true
./server
```
Или положи бинарник и `migrations/` в `/opt/dayzsmartcf`, создай systemd-юнит (см. `deploy/systemd/dayzsmartcf-backend.service.example`).

### Фронтенд
```bash
cd frontend
npm ci
npm run build
```
Скопируй содержимое `dist/` в каталог nginx (например `/var/www/dayzsmartcf`) и настрой проксирование `/api/` и `/health` на `http://127.0.0.1:8080`.

В папке `deploy/` есть примеры:
- **systemd/dayzsmartcf-backend.service.example** — юнит для автозапуска бэкенда.
- **nginx-bare.conf.example** — конфиг nginx для раздачи статики и прокси на бэкенд.
