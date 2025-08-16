#!/bin/bash
set -e

# Конфиг
CONTAINER_NAME="retry-postgres"
DUMPS_DIR="./dumps"
DATE_STR=$(date +%F)
FILENAME="${DUMPS_DIR}/dump_${DATE_STR}.sql"

# Создаём папку для дампов, если её нет
mkdir -p "$DUMPS_DIR"

# Делаем дамп
docker exec -t "$CONTAINER_NAME" pg_dump \
  -U retry_direct_user \
  -d retry \
  > "$FILENAME"

echo "✅ Дамп создан: $FILENAME"
