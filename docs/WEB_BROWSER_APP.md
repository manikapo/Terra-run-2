# Browser web app (Phase 1 — no native build)

Full GPS territory run in **Chrome on your phone**, hosted at **run.8me.in**.

Native Android/iOS app → **Phase 2**.

## Open the app

**https://run.8me.in/ota/index.html**

(If your host root is `public_html`, upload so `ota/index.html` is at that path.)

---

## One-time setup (15 min)

### 1. Edit `js/config.js` before upload

```javascript
window.TERRITORY_CONFIG = {
  API_BASE_URL: "https://YOUR-SERVICE.onrender.com",
  SUPABASE_URL: "https://xxxx.supabase.co",
  SUPABASE_ANON_KEY: "eyJ...",
  APP_URL: "https://run.8me.in/ota/index.html",
  H3_CAPTURE_RES: 10,
  H3_TILE_RES: 7,
};
```

### 2. Supabase → Authentication → URL Configuration

Add **Redirect URL**:

```
https://run.8me.in/ota/index.html
```

Enable **Google** provider (same as before).

### 3. Upload folder to server

Upload everything in `ota/public/` to `run.8me.in/ota/`:

```
ota/
  index.html      ← main app (map + GPS + claim)
  test.html       ← old bridge test (optional)
  css/app.css
  js/config.js    ← you edited this
  js/app.js
  js/config.example.js
  pages/profile.html
  manifest.json
```

### 4. Optional: root redirect

Upload `root-redirect.html` as `index.html` at site **root** so `run.8me.in` opens the app.

---

## How to use

1. Open **https://run.8me.in/ota/index.html** on your **phone** (Chrome).
2. **Sign in with Google**.
3. Allow **location** when asked.
4. Tap **Start run** — walk/jog outdoors.
5. See your route (teal line) and **cells** count increase.
6. Tap **Finish & claim** — territories sync to Render/Supabase.
7. Colored hexes on map: **teal = yours**, **orange = others**.

---

## What works in browser vs Phase 2 native app

| Feature | Browser (now) | Native app (Phase 2) |
|---------|---------------|----------------------|
| Map + territories | ✅ | ✅ |
| GPS run + claim | ✅ (browser GPS) | ✅ (better background GPS) |
| Google login | ✅ | ✅ |
| Offline / background | ❌ limited | ✅ |
| Anti-cheat | Basic | Stronger |
| Push notifications | ❌ | ✅ |

---

## Troubleshooting

| Issue | Fix |
|-------|-----|
| Login loops | Add redirect URL in Supabase |
| Map blank | Check browser console; allow mixed content if any |
| GPS not updating | Use phone, not desktop; allow location |
| API 401 | Sign out and sign in again |
| No colored territories | Finish a run first; tap Refresh map |
