# Android Studio

Open the `android/` folder.

1. Copy `local.properties.example` → `local.properties`
2. Set API + Supabase + OTA URLs (see `docs/DEPLOYMENT.md`)
3. Sync Gradle, run on a **physical device** with GPS

## Modules (Phase 1)

- `gps` — Fused Location Provider, Kalman filter, outlier rejection
- `territory` — H3 cell capture on device
- `data/local` — Room (GPS points + activities)
- `data/remote` — Retrofit API client
- `service` — Foreground tracking notification
- `ota` — WebView + bundle download from R2

## TODO before production

- [ ] Google Credential Manager in `LoginScreen`
- [ ] Map screen with Mapbox territory layer
- [ ] Manifest signature verification for OTA
- [ ] Play Integrity for anti-cheat trust
