#!/usr/bin/env bash
set -euo pipefail

echo "==> Starting PostgreSQL"
if sudo pg_ctlcluster 16 main status >/dev/null 2>&1; then
  echo "PostgreSQL already running"
else
  sudo pg_ctlcluster 16 main start || sudo service postgresql start
fi

# Wait until Postgres accepts connections
for i in $(seq 1 30); do
  if pg_isready -h localhost -p 5432 >/dev/null 2>&1; then
    echo "PostgreSQL is ready"
    exit 0
  fi
  sleep 1
done

echo "PostgreSQL failed to become ready within 30s" >&2
exit 1
