# Infinite Me Website

Static site for [8me.in](https://8me.in), built from the INTVL export with 8me.in branding.

## Deploy to Hostinger

Upload the **entire** `website/` folder contents to your web root:

```
index.html
css/
js/
_next/          ← must include BOTH subfolders below
  *.jpg/png     ← root-level image aliases (heroGlobee458.jpg, etc.)
  static/
    media/      ← all images, fonts, SVGs
    chunks/     ← JavaScript bundles
    css/
```

**Common mistake:** uploading only `index.html` and `_next/static/` but missing the `_next/*.jpg` files at the root of `_next/`. The site needs the full `_next/` directory exactly as in the zip.

**Do not** open `index.html` directly via `file://` — use an HTTP server or upload to Hostinger.

## Local preview

```bash
cd website
python3 -m http.server 8080
```

Open http://localhost:8080 (must use HTTP server, not `file://`).

## Files

- `index.html` — main page (INTVL scroll animations preserved)
- `_next/` — CSS, JS, images, fonts from Next.js export
- `css/8me-theme.css` — orange color overrides for 8me.in brand
- `js/8me.js` — Android download modal + CTA link fixes

## Brand colors

| Token | Value |
|-------|-------|
| Orange | `#FF4D00` |
| Black | `#0D0D0D` |
| Off-white | `#F5F3EE` |

## Links

- Play: https://play.8me.in
- Android testers: https://groups.google.com/g/infiniteme-testers
- Play Store testing: https://play.google.com/apps/testing/in.me8.infiniteme
