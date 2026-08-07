# WebView test guide (beginner) — domain run.8me.in

Test the OTA profile UI in the Android app WebView using your domain **https://run.8me.in**.

You need **3 pieces** working:

| Piece | Where | Purpose |
|-------|--------|---------|
| **OTA HTML** | `run.8me.in/ota/` | Profile page the WebView loads |
| **API** | Render (free) | Profile stats (`/api/v1/users/me`) |
| **Auth** | Supabase (free) | JWT token for the API |

---

## Part A — Upload files to run.8me.in (15 min)

### What to upload

Copy everything from this folder on your PC:

```
C:\Users\VCL-ManishK\Projects\territory-run\ota\public\
```

Also copy the **bundle contents** so CSS/JS work:

```
C:\Users\VCL-ManishK\Projects\territory-run\ota\bundle\pages\  →  run.8me.in/ota/pages/
C:\Users\VCL-ManishK\Projects\territory-run\ota\bundle\css\    →  run.8me.in/ota/css/
C:\Users\VCL-ManishK\Projects\territory-run\ota\bundle\js\     →  run.8me.in/ota/js/
```

Plus `ota\public\test.html` and `ota\public\manifest.json` → `run.8me.in/ota/`

### Final structure on your server

```
https://run.8me.in/ota/
  manifest.json
  test.html
  pages/profile.html
  css/app.css
  js/profile.js
  v1/bundle.zip          ← optional (see Part A.2)
```

### How to upload (pick one)

**cPanel / hosting file manager**

1. Log in to your host for `8me.in`.
2. Open `public_html` or `www` (or subdomain folder for `run`).
3. Create folder `ota`.
4. Upload files keeping the folder structure above.

**FTP (FileZilla)**

1. Host: your server IP or `ftp.8me.in`
2. Upload local `ota\public` + `ota\bundle` folders into `/ota/` on server.

### A.1 — Quick browser check (not WebView, but confirms hosting)

Open in Chrome:

- https://run.8me.in/ota/test.html  
  → Should show “No bridge” (normal in browser).

- https://run.8me.in/ota/pages/profile.html  
  → Dark profile UI; status “Bridge unavailable” in browser.

If these URLs **404**, fix hosting before testing the app.

### A.2 — Optional: bundle.zip (full OTA download)

In PowerShell:

```powershell
cd C:\Users\VCL-ManishK\Projects\territory-run\ota\scripts
.\build-bundle.ps1
```

Upload `ota\public\v1\bundle.zip` to `run.8me.in/ota/v1/bundle.zip` and re-upload `ota\public\manifest.json` (script updates SHA256).

For **first test**, skip zip — the app loads `https://run.8me.in/ota/pages/profile.html` directly.

---

## Part B — Supabase (auth + database) (~20 min)

1. Go to [supabase.com](https://supabase.com) → **New project** (free).
2. **SQL Editor** → paste and run:
   `supabase/migrations/001_initial_schema.sql`
3. **Settings → API** — copy:
   - Project URL → `https://xxxx.supabase.co`
   - `anon` public key
   - JWT Secret (under JWT Settings)
4. **Authentication → Providers → Google** — enable when ready (optional for first WebView-only test).

**Quick test without Google:** Authentication → Users → **Add user** → create email/password user, then sign in via Supabase in app later.

---

## Part C — Render API (~15 min)

1. Push code to GitHub (you already linked `Terra-run-2`).
2. [render.com](https://render.com) → **New → Web Service** → repo `Terra-run-2`.
3. Settings:
   - **Root Directory:** leave empty
   - **Runtime:** Docker
   - **Dockerfile Path:** `backend/Dockerfile`
   - **Docker Context:** `backend`
   - **Health Check Path:** `/health`
4. **Environment variables:**

| Key | Value |
|-----|--------|
| `DATABASE_URL` | Supabase connection string (pooler, port 6543) |
| `SUPABASE_URL` | `https://xxxx.supabase.co` |
| `SUPABASE_JWT_SECRET` | From Supabase JWT settings |
| `ENV` | `production` |

5. Deploy → copy URL e.g. `https://terra-run-2.onrender.com`

6. Test:

```text
https://YOUR-SERVICE.onrender.com/health
```

Should return `{"status":"ok",...}`

7. **UptimeRobot** (optional): ping `/health` every 5 minutes so free tier stays warm.

---

## Part D — Configure Android app (~10 min)

1. Open `android/` in **Android Studio**.
2. Create `android/local.properties` (copy from `local.properties.example`):

```properties
sdk.dir=C\:\\Users\\VCL-ManishK\\AppData\\Local\\Android\\Sdk
API_BASE_URL=https://YOUR-SERVICE.onrender.com
SUPABASE_URL=https://xxxx.supabase.co
SUPABASE_ANON_KEY=eyJhbG...your-anon-key
OTA_MANIFEST_URL=https://run.8me.in/ota/manifest.json
OTA_PROFILE_URL=https://run.8me.in/ota/pages/profile.html
MAPBOX_ACCESS_TOKEN=pk.test
```

3. **Sync Gradle** → run on a **real phone** (emulator OK for WebView-only test).

---

## Part E — Test WebView in the app

1. Install and open **Territory Run**.
2. Tap through login (dev: login button may skip auth — see note below).
3. Tap **Profile (OTA)**.
4. WebView should load **https://run.8me.in/ota/pages/profile.html**.

### What you should see

| Screen | Meaning |
|--------|---------|
| Profile UI loads, styled dark theme | OTA hosting on `run.8me.in` works |
| “Bridge unavailable” | Opened in Chrome, not app WebView |
| “Token: empty” on test page | Sign in with Supabase not wired yet |
| “Synced from API” on profile | Full stack works |
| “Failed: HTTP 401” | API URL wrong or no valid JWT |
| Blank / connection error | Wrong URL or SSL issue on domain |

### Test page inside app

Temporarily change `OTA_PROFILE_URL` to:

```properties
OTA_PROFILE_URL=https://run.8me.in/ota/test.html
```

Rebuild → Profile (OTA) → should show **Bridge OK** and API URL.

---

## Part F — Get profile data showing (needs login)

The profile page calls your API with a Supabase JWT. Without login:

- UI loads from `run.8me.in` ✅  
- Stats stay 0 / “Token empty” ❌

To see real stats:

1. Wire **Google Sign-In** → `SupabaseAuthManager.signInWithGoogleIdToken()`  
   **OR** use Supabase email login in a dev build.
2. After login, open Profile (OTA) again → should show cells/score from API.

---

## Troubleshooting

| Problem | Fix |
|---------|-----|
| 404 on run.8me.in | Re-check upload paths; `profile.html` must be at `/ota/pages/profile.html` |
| Mixed content / blocked | Use `https://` everywhere (not `http://`) |
| WebView blank | Enable internet permission; check phone can open URL in Chrome |
| CORS errors | API is called from WebView with Bearer token — CORS less critical; 401 = auth issue |
| Render slow first request | Cold start; use UptimeRobot or wait 30–60s |
| CSS missing | Ensure `/ota/css/app.css` exists (paths in HTML use `../css/`) |

---

## Optional — point API at your domain

In your DNS for `8me.in`:

- `api.run.8me.in` → CNAME → `your-service.onrender.com`

Then set `API_BASE_URL=https://api.run.8me.in` in the app.

---

## Checklist (minimum WebView test)

- [ ] https://run.8me.in/ota/pages/profile.html opens in browser  
- [ ] https://run.8me.in/ota/test.html opens in browser  
- [ ] Render `/health` returns OK  
- [ ] Supabase SQL migration ran  
- [ ] `local.properties` has correct URLs  
- [ ] App → Profile (OTA) shows the profile page  

That’s enough to confirm **WebView + remote HTML on run.8me.in** works. Add Supabase login to see live API data.
