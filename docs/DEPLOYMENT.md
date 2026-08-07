# Territory Run — Phase 1 Deployment Guide

Where to upload each piece and in what order.

---

## Overview

| Component | Upload to | Cost |
|-----------|-----------|------|
| Database schema | **Supabase** SQL Editor | Free |
| Auth (Google) | **Supabase** Auth settings | Free |
| REST API | **Render** (Docker) | Free |
| Uptime keep-alive | **UptimeRobot** → `/health` | Free |
| OTA HTML/CSS/JS | **Cloudflare R2** + CDN | Free tier |
| Android app | **Google Play** (later) / APK sideload | — |
| iOS app | **App Store** (later) | $99/year Apple |

---

## Step 1 — Supabase (database + auth)

### 1.1 Create project

1. Go to [https://supabase.com](https://supabase.com) → New project (free).
2. Save **Project URL** and **anon public** key (Settings → API).
3. Save **JWT Secret** (Settings → API → JWT Settings) — needed for Render API.

### 1.2 Run migration

1. Dashboard → **SQL Editor** → New query.
2. Paste contents of `supabase/migrations/001_initial_schema.sql`.
3. Run.

### 1.3 Google OAuth

1. [Google Cloud Console](https://console.cloud.google.com) → APIs → Credentials → OAuth 2.0 Client.
2. Create **Web client** (for Supabase) + **Android** / **iOS** clients for native apps.
3. Supabase → Authentication → Providers → **Google** → enable, paste Web client ID + secret.
4. Add redirect URL from Supabase to Google authorized redirects.

### 1.4 Database connection string

Settings → Database → **Connection string** → URI (Transaction pooler, port 6543).

Use this as `DATABASE_URL` on Render.

---

## Step 2 — Render (API)

### 2.1 Push code to GitHub

```bash
cd C:\Users\VCL-ManishK\Projects\territory-run
git init
git add .
git commit -m "Phase 1 MVP scaffold"
git remote add origin https://github.com/YOUR_USER/territory-run.git
git push -u origin main
```

### 2.2 Create Web Service

1. [https://render.com](https://render.com) → New → **Web Service**.
2. Connect GitHub repo `Terra-run-2`.
3. Settings:
   - **Root Directory**: leave empty (repo root)
   - **Runtime**: Docker
   - **Dockerfile Path**: `Dockerfile` (at repo root — **not** `backend/Dockerfile`)
   - **Docker Context**: `.` (repo root)
   - **Plan**: Free
   - **Health Check Path**: `/health`

Or use Blueprint: Dashboard → New → Blueprint → connect repo (`render.yaml` included).

### 2.3 Environment variables (Render → Environment)

| Key | Value |
|-----|--------|
| `DATABASE_URL` | Supabase pooler URI |
| `SUPABASE_URL` | `https://xxxx.supabase.co` |
| `SUPABASE_JWT_SECRET` | From Supabase JWT settings |
| `INTERNAL_JOB_SECRET` | Random string (openssl rand -hex 32) |
| `ENV` | `production` |
| `H3_CAPTURE_RESOLUTION` | `10` |

### 2.4 Deploy

Render builds Docker image and deploys. Your API URL:

`https://territory-run-api.onrender.com` (or your service name).

Test:

```bash
curl https://YOUR-SERVICE.onrender.com/health
```

---

## Step 3 — UptimeRobot (keep Render free tier awake)

1. [https://uptimerobot.com](https://uptimerobot.com) → Add monitor.
2. Type: **HTTP(s)**.
3. URL: `https://YOUR-SERVICE.onrender.com/health`
4. Interval: **5 minutes**.
5. Alert when down (email).

This hits lightweight `/health` (no database query).

**Note:** One always-on free service ≈ 720 hrs/month (within Render 750 hr limit).

---

## Step 4 — Cloudflare R2 (OTA bundles)

### 4.1 Create bucket

1. Cloudflare Dashboard → **R2** → Create bucket `territory-run`.
2. Enable public access via **R2.dev subdomain** or custom domain on Cloudflare (free CDN).

### 4.2 Build OTA bundle

On Mac/Linux or Git Bash on Windows:

```bash
cd ota
bash scripts/build-bundle.sh
```

Outputs:

- `ota/dist/v1/bundle.zip`
- `ota/dist/v1/manifest.json` (with sha256)

Update `bundle_url` in manifest to your real CDN URL before upload.

### 4.3 Upload files

Upload to R2 (Dashboard or `aws s3 cp` with R2 credentials):

| File | R2 path |
|------|---------|
| `bundle.zip` | `ota/v1/bundle.zip` |
| `manifest.json` | `ota/manifest.json` |

Public URLs (example):

- `https://pub-xxxx.r2.dev/ota/manifest.json`
- `https://pub-xxxx.r2.dev/ota/v1/bundle.zip`

Set `OTA_MANIFEST_URL` in Android `local.properties` and iOS Info.plist.

### 4.4 (Optional) Signature

Phase 1 uses SHA256 only. Phase 2: sign manifest with RSA private key in CI; embed public key in native app.

---

## Step 5 — Android app

### 5.1 Open in Android Studio

1. Open folder `android/`.
2. Copy `local.properties.example` → `local.properties`.
3. Set:

```properties
sdk.dir=C\:\\Users\\YOUR_NAME\\AppData\\Local\\Android\\Sdk
API_BASE_URL=https://YOUR-SERVICE.onrender.com
SUPABASE_URL=https://xxxx.supabase.co
SUPABASE_ANON_KEY=eyJ...
OTA_MANIFEST_URL=https://pub-xxxx.r2.dev/ota/manifest.json
MAPBOX_ACCESS_TOKEN=pk....
```

4. Sync Gradle → Run on physical device (GPS required).

### 5.2 Google Sign-In

Wire Credential Manager / Google Sign-In in `LoginScreen.kt` → `SupabaseAuthManager.signInWithGoogleIdToken()`.

Until then, use Supabase dashboard test user or temporary dev bypass.

---

## Step 6 — iOS app

See `ios/README.md`.

1. Create Xcode project, copy `ios/TerritoryRun/` sources.
2. Add Supabase Swift package.
3. Set Info.plist config keys (API, Supabase, OTA URLs).
4. Enable Background Location capability.
5. Run on device.

---

## Step 7 — Verify end-to-end

1. `curl /health` → `{"status":"ok"}`
2. Sign in on mobile → get Supabase JWT.
3. Start run → GPS points stored locally.
4. Finish run → API `POST /activities` → `points` → `complete`.
5. Check Supabase Table Editor → `territory_cells`, `activities`.
6. Open Profile (OTA) → stats load from API.

---

## API quick reference

```http
GET  /health
GET  /api/v1/users/me              Authorization: Bearer <supabase_jwt>
POST /api/v1/activities            { started_at, idempotency_key }
POST /api/v1/activities/{id}/points { points: [{lat, lon, ts, acc, speed}] }
POST /api/v1/activities/{id}/complete { ended_at, distance_m, duration_s, h3_cells_hint }
GET  /api/v1/territories/tile/{h3_7}  → GeoJSON
```

---

## Troubleshooting

| Issue | Fix |
|-------|-----|
| Render cold start 30–60s | UptimeRobot 5 min ping; or upgrade Render $7/mo |
| 401 from API | JWT expired; refresh Supabase session |
| DB connection fails | Use Supabase **pooler** URI; check password |
| OTA page blank | Run `OtaManager.sync`; check R2 public URL + sha256 |
| No territory capture | Activity `quarantined`/`rejected` — check `trust_score` rules |
| Render suspended (750 hrs) | Second always-on service — use GitHub Actions cron instead |

---

## What to upload where (cheat sheet)

```
territory-run/
  supabase/migrations/*.sql     → Supabase SQL Editor (run once)
  backend/                      → GitHub → Render auto-deploy
  ota/dist/*                    → Cloudflare R2 (manual or CI)
  android/                      → Build locally → Play Store / APK
  ios/                          → Build locally → App Store
```

---

## Optional — GitHub Actions (decay job stub)

Create `.github/workflows/decay.yml` to POST internal job weekly (uses `INTERNAL_JOB_SECRET`).

No payment required for Phase 1 beyond optional Apple Developer Program for iOS distribution.
