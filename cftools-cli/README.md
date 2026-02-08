# CFtools CLI

Отдельный Go-проект для выполнения запросов к CFtools API с твоими параметрами.

## Установка

```bash
cd cftools-cli
go mod tidy
```

## Конфигурация

Создай `.env` в папке `cftools-cli` или скопируй `.env` из корня проекта:

```
CFTOOLS_CDN_AUTH=eyJ...твой_токен...
```

Либо передай токен флагом: `-token "eyJ..."`

## Использование

**Поиск по нику (global-query):**
```bash
go run . -query "Nickname123"
```

**Профиль игрока (status, overview, playState и т.д.):**
```bash
go run . -action profile-status -profile 637d4e4290d48f5870f81294
go run . -action profile-overview -profile 637d4e4290d48f5870f81294
go run . -action profile-playState -profile 637d4e4290d48f5870f81294
go run . -action profile-structure -profile 637d4e4290d48f5870f81294
go run . -action profile-activities -profile 637d4e4290d48f5870f81294
```

**Флаги:**
- `-query` — ник или identifier для поиска
- `-profile` — CFtools ID для запроса профиля
- `-action` — действие: `global-query` | `profile-status` | `profile-overview` | `profile-playState` | `profile-structure` | `profile-activities`
- `-token` — CFTOOLS_CDN_AUTH (если не в .env)

## Примеры

```bash
# Поиск по нику
go run . -query "Steam_76561198..."

# Профиль: статус
go run . -action profile-status -profile 637d4e4290d48f5870f81294

# Профиль: обзор
go run . -action profile-overview -profile 637d4e4290d48f5870f81294
```
