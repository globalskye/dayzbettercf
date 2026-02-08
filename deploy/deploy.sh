#!/usr/bin/env bash
# Запуск DayZ Smart CF на Linux через Docker.
# Запускай из корня проекта: ./deploy/deploy.sh

set -e
cd "$(dirname "$0")/.."

if [ ! -f .env ]; then
  echo "Файл .env не найден. Скопируй и отредактируй:"
  echo "  cp .env.docker.example .env"
  echo "  nano .env   # вставь CFTOOLS_* и ADMIN_USER, ADMIN_PASS"
  exit 1
fi

echo "Сборка и запуск контейнеров..."
docker compose up -d --build

echo ""
echo "Готово. Приложение доступно по адресу:"
echo "  http://localhost:3000"
echo "  (или http://IP_СЕРВЕРА:3000)"
echo ""
echo "Логи: docker compose logs -f"
echo "Остановить: ./deploy/stop.sh"
