#!/usr/bin/env bash
# Остановка контейнеров DayZ Smart CF.

set -e
cd "$(dirname "$0")/.."

docker compose down
echo "Контейнеры остановлены."
