# INTVL-style GPS running app — layout

How [INTVL](https://www.intvl.com.au/) maps onto **Territory Run** (this repo), and the phased path to feature parity.

This is the product layout for the Go + H3 stack in this repository. Infinite Me (`play.8me.in`, Python/Railway) is a separate live product with polygon steal; do not merge those backends.

---

## What INTVL is

A gamified GPS fitness app: runs and rides become territory conquest.

| Layer | INTVL | Territory Run |
|-------|--------|----------------|
| Tracking | GPS run/ride | Android / iOS / browser GPS |
| Territory | Route → map cells; loops claim area | H3 resolution 10 along the route |
| Competition | Steal, leaderboards, battles | Phase 2: steal + leaderboards + decay |
| Social | Clubs, lobbies | Phase 3 |
| Rewards | Seasonal prizes | Phase 3 |
| Analytics | Splits, elevation | Distance / duration / cells today |
| Integrations | Strava, watches | Phase 2+ |

---

## Architecture

```
Clients (Android / iOS / Leaflet PWA)
    ↓ JWT or X-Guest-User
Go API on Render
    ├── activities (record + complete)
    ├── territories (H3 tiles + steal on complete)
    ├── leaderboards (global + local tile)
    └── internal jobs (decay cron)
    ↓
Supabase Postgres
    users, activities, territory_cells, territory_events,
    user_territory_stats, leaderboard_snapshots, battles
```

---

## Screen / UX map

```
ONBOARDING     Splash → guest or Google → map
HOME           Full-screen territory map, Start run, rank, notifications
RUN TRACKING   Live polyline, time / distance / pace / cells, Pause / Finish
POST-RUN       Capture animation, new vs stolen cells, score
ACTIVITY       Route + overlay, share
LEADERBOARDS   Global and local (H3 tile ~5 km)
PROFILE        Cells owned, score, streak, history
BATTLES        Phase 2 table exists; UI later
CLUBS          Phase 3
```

---

## Game loop

1. **Start** — `POST /api/v1/activities`
2. **Record** — buffer GPS locally; `POST /api/v1/activities/:id/points`
3. **Finish** — `POST /api/v1/activities/:id/complete`
   - Anti-cheat (speed, accuracy, min duration)
   - Route → unique H3 cells
   - Neutral cell → capture (10 pts)
   - Own cell → defend (2 pts, refresh `last_defended_at`)
   - Rival cell → **steal** (15 pts, `STOLEN` event, victim `territories_lost`)
4. **Map** — `GET /api/v1/territories/tile/:h3_tile` (resolution 7 parent)
5. **Compete** — `GET /api/v1/leaderboards/global` and `.../local?h3_parent=`
6. **Decay** — daily `POST /internal/jobs/decay` (undefended cells → neutral)

---

## H3 grid

| Resolution | Role | Approx. size |
|------------|------|----------------|
| 10 | Capture unit | ~66 m edge |
| 7 | Map tile + local leaderboard | ~5 km edge |

---

## Phases

| Phase | Goal | Status in this repo |
|-------|------|---------------------|
| 1 | Run → capture → see map | Done |
| 2 | Steal, decay, leaderboards | This change |
| 3 | Clubs, feed, prizes, ride mode | Not started |
| 4 | Watches, Pro, PostGIS heatmaps | Not started |

---

## Deploy order (Phase 2)

1. Run `supabase/migrations/002_phase2_steal_leaderboards.sql` in the SQL editor.
2. Redeploy the Render API.
3. Upload `ota/public/` (cache-bust `?v=4` on CSS/JS).
4. Confirm GitHub Action `decay.yml` has `API_BASE_URL` and `INTERNAL_JOB_SECRET`.

See [DEPLOYMENT.md](DEPLOYMENT.md) for Phase 1 hosting.
