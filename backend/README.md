# Backend API

Go Fiber API for Territory Run Phase 1.

## Local run

```bash
cd backend
cp .env.example .env
# Edit .env with Supabase credentials

go mod tidy
go run ./cmd/server
```

## Docker

```bash
docker build -t territory-run-api .
docker run -p 8080:8080 --env-file .env territory-run-api
```

## Deploy

Push to GitHub; Render builds from `backend/Dockerfile`.

Health check: `GET /health` (no DB — safe for UptimeRobot).
