#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Usage: $0 <backup.dump> [container] [database] [user]" >&2
  exit 1
fi

FILE="$1"
CONTAINER="${2:-weighing_postgres}"
DATABASE="${3:-weighing_db}"
DB_USER="${4:-weighing}"

if [ ! -f "$FILE" ]; then
  echo "Backup file not found: $FILE" >&2
  exit 1
fi

echo "Restoring $FILE into $DATABASE ($CONTAINER)..."
echo "WARNING: existing objects will be dropped (--clean --if-exists)."

docker exec -i "$CONTAINER" pg_restore -U "$DB_USER" -d "$DATABASE" --clean --if-exists --no-owner < "$FILE"

echo "Restore completed."
