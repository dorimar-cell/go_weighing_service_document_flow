#!/usr/bin/env bash
set -euo pipefail

# Удаляет фото и записи weighting_photo старше DAYS дней.
# Запускать вручную или из cron; сервис не чистит архив сам.

CONTAINER="${CONTAINER:-weighing_postgres}"
DATABASE="${DATABASE:-weighing_db}"
DB_USER="${DB_USER:-weighing}"
DAYS="${DAYS:-90}"
DRY_RUN="${DRY_RUN:-0}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PHOTOS_DIR="${PHOTOS_DIR:-$SCRIPT_DIR/../data/photos}"

if [ ! -d "$PHOTOS_DIR" ]; then
  echo "Photos directory not found: $PHOTOS_DIR" >&2
  exit 1
fi

rows="$(docker exec "$CONTAINER" psql -U "$DB_USER" -d "$DATABASE" -t -A -F '|' \
  -c "SELECT id, id_weighing, file_name FROM weighting_photo
      WHERE created_date < now() - interval '$DAYS days' ORDER BY id")"

if [ -z "$rows" ]; then
  echo "No photos older than $DAYS days."
  exit 0
fi

count="$(printf '%s\n' "$rows" | wc -l | tr -d ' ')"
echo "Photos to remove: $count (older than $DAYS days)"

ids=""
deleted=0
missing=0

while IFS='|' read -r id session_id file_name; do
  [ -z "$id" ] && continue
  ids="${ids:+$ids,}$id"
  path="$PHOTOS_DIR/$session_id/$file_name"
  if [ -f "$path" ]; then
    if [ "$DRY_RUN" = "1" ]; then
      echo "Would delete: $path"
    else
      rm -f -- "$path"
    fi
    deleted=$((deleted + 1))
  else
    missing=$((missing + 1))
  fi
done <<< "$rows"

echo "Files: $deleted, already missing: $missing"

if [ "$DRY_RUN" = "1" ]; then
  echo "Dry run: database rows are kept."
  exit 0
fi

docker exec "$CONTAINER" psql -U "$DB_USER" -d "$DATABASE" \
  -c "DELETE FROM weighting_photo WHERE id IN ($ids)" >/dev/null
echo "Removed rows: $count"
