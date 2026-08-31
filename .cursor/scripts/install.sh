#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$REPO_ROOT"

echo "==> Installing Go backend dependencies"
cd backend
go mod download
go mod tidy
cd "$REPO_ROOT"

echo "==> Ensuring PostgreSQL is running"
if ! sudo pg_ctlcluster 16 main status >/dev/null 2>&1; then
  sudo pg_ctlcluster 16 main start || sudo service postgresql start
fi

echo "==> Setting up local database"
sudo -u postgres psql -tc "SELECT 1 FROM pg_roles WHERE rolname='territory_run'" | grep -q 1 \
  || sudo -u postgres psql -c "CREATE USER territory_run WITH PASSWORD 'territory_run_dev';"

sudo -u postgres psql -tc "SELECT 1 FROM pg_database WHERE datname='territory_run'" | grep -q 1 \
  || sudo -u postgres psql -c "CREATE DATABASE territory_run OWNER territory_run;"

sudo -u postgres psql -d territory_run -c "GRANT ALL ON SCHEMA public TO territory_run;" 2>/dev/null || true

sudo -u postgres psql -d territory_run -f "$REPO_ROOT/supabase/migrations/001_initial_schema.sql" 2>/dev/null || true

sudo -u postgres psql -d territory_run -c "
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO territory_run;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO territory_run;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO territory_run;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO territory_run;
ALTER TABLE users DISABLE ROW LEVEL SECURITY;
ALTER TABLE activities DISABLE ROW LEVEL SECURITY;
ALTER TABLE territory_cells DISABLE ROW LEVEL SECURITY;
" 2>/dev/null || true

echo "==> Writing local development config"
cat > "$REPO_ROOT/backend/.env" <<'EOF'
PORT=8080
ENV=development
DATABASE_URL=postgresql://territory_run:territory_run_dev@localhost:5432/territory_run
SUPABASE_URL=http://localhost:54321
SUPABASE_JWT_SECRET=local-dev-jwt-secret-not-for-production
INTERNAL_JOB_SECRET=local-dev-internal-secret
ALLOW_GUEST_AUTH=true
H3_CAPTURE_RESOLUTION=10
H3_TILE_RESOLUTION=7
EOF

cat > "$REPO_ROOT/ota/public/js/config.js" <<'EOF'
window.TERRITORY_CONFIG = {
  API_BASE_URL: "http://localhost:8080",
  GUEST_MODE: true,
  H3_CAPTURE_RES: 10,
  H3_TILE_RES: 7,
};
EOF

echo "==> Install complete"
