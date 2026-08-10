# Infinite Me Website

Static site for [8me.in](https://8me.in), built from the INTVL export with 8me.in branding.

## Deploy to Hostinger

Upload the entire `website/` folder contents to your web root, including the `_next/` directory.

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
