#!/usr/bin/env bash
set -euo pipefail

CONTAINER="${CONTAINER:-weighing_postgres}"
DATABASE="${DATABASE:-weighing_db}"
DB_USER="${DB_USER:-weighing}"
KEEP_DAYS="${KEEP_DAYS:-14}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUT_DIR="${OUT_DIR:-$SCRIPT_DIR/../backups}"

mkdir -p "$OUT_DIR"
STAMP="$(date +%Y%m%d-%H%M%S)"
FILE="$OUT_DIR/$DATABASE-$STAMP.dump"

docker exec "$CONTAINER" pg_dump -U "$DB_USER" -d "$DATABASE" -Fc > "$FILE"

if [ ! -s "$FILE" ]; then
  echo "Backup file is empty: $FILE" >&2
  exit 1
fi

echo "Backup created: $FILE ($(du -h "$FILE" | cut -f1))"

find "$OUT_DIR" -name "$DATABASE-*.dump" -mtime +"$KEEP_DAYS" -print -delete
