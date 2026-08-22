# Territory Run — Phase 1 MVP

GPS territory capture running app (Strava / INTVL-inspired). Phase 1: capture-only territories. Phase 2: steal, leaderboards, decay.

Full layout: **[docs/INTVL_LAYOUT.md](docs/INTVL_LAYOUT.md)**.

## Repository layout

```
territory-run/
├── backend/          # Go API → deploy to Render
├── supabase/         # SQL migrations → run in Supabase dashboard
├── android/          # Kotlin app → Android Studio
├── ios/              # Swift app → Xcode
├── ota/              # HTML UI bundle → Cloudflare R2
├── render.yaml       # Render Blueprint (optional)
└── docs/DEPLOYMENT.md  # Step-by-step upload guide
```

## Quick start

1. Read **[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)** for where to upload each piece.
2. Run Supabase migration: `supabase/migrations/001_initial_schema.sql`
3. Deploy API: connect `backend/` to Render (see DEPLOYMENT.md)
4. Upload OTA bundle to R2
5. Open `android/` or `ios/` in IDE, set env keys, run on device

## Phase 1 scope

- [x] Google login (Supabase Auth)
- [x] GPS tracking (foreground service / Core Location)
- [x] H3 territory capture (resolution 10)
- [x] Activity upload + server validation
- [x] Basic territory map tiles (GeoJSON API)
- [x] User profile stats
- [x] OTA profile page (WebView)
- [x] Territory steal (Phase 2 — H3 cell takeover)
- [x] Global / local leaderboards
- [x] Territory decay cron
- [ ] FCM push (later)
- [ ] Strava login (later)

## Environment variables

See `backend/.env.example` and mobile `local.properties` / `Config.xcconfig` templates.
