#!/usr/bin/env bash
set -euo pipefail

# Проверка резервной копии: разворачивает дамп во временную БД,
# считает ключевые таблицы и удаляет временную БД. Рабочую БД не трогает.
#
# Использование: ./verify_backup.sh backups/weighing_db-20260101-020000.dump

CONTAINER="${CONTAINER:-weighing_postgres}"
DB_USER="${DB_USER:-weighing}"
KEEP="${KEEP:-0}"
FILE="${1:-}"

if [ -z "$FILE" ] || [ ! -f "$FILE" ]; then
  echo "Backup file not found: ${FILE:-<empty>}" >&2
  exit 1
fi

TEMP_DB="weighing_verify_$(date +%Y%m%d_%H%M%S)"
CREATED=0

cleanup() {
  if [ "$CREATED" = "1" ] && [ "$KEEP" != "1" ]; then
    echo "Dropping temporary database $TEMP_DB..."
    docker exec "$CONTAINER" dropdb -U "$DB_USER" --if-exists "$TEMP_DB" >/dev/null || true
  fi
}
trap cleanup EXIT

echo "Creating temporary database $TEMP_DB..."
docker exec "$CONTAINER" createdb -U "$DB_USER" "$TEMP_DB"
CREATED=1

echo "Restoring $FILE into $TEMP_DB..."
docker exec -i "$CONTAINER" pg_restore -U "$DB_USER" -d "$TEMP_DB" --no-owner < "$FILE"

counts="$(docker exec "$CONTAINER" psql -U "$DB_USER" -d "$TEMP_DB" -t -A -F '|' \
  -c "SELECT
        (SELECT count(*) FROM weighing),
        (SELECT count(*) FROM weighing_measurement),
        (SELECT count(*) FROM dir_users),
        (SELECT count(*) FROM schema_migrations)")"

migrations="$(echo "$counts" | cut -d'|' -f4)"
if [ -z "$counts" ] || [ "${migrations:-0}" -le 0 ]; then
  echo "Restore verification failed: $counts" >&2
  exit 1
fi

echo "Restore verified:"
echo "  sessions|measurements|users|migrations = $counts"
if [ "$KEEP" = "1" ]; then
  echo "Temporary database kept: $TEMP_DB"
fi
echo "Backup restore check passed."
