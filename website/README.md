# INTVL — Static HTML

Exact static HTML conversion of the INTVL Next.js landing page export.

## Deploy

Upload the **entire** `website/` folder contents to your web root:

```
index.html
favicon.ico
iconb9d2.svg
apple-iconfb83.png
_next/          ← full folder required
  *.jpg/png     ← root-level image aliases
  static/
    media/      ← images, fonts, SVGs
    chunks/     ← JavaScript
    css/        ← styles
```

Do **not** open via `file://` — use an HTTP server.

## Local preview

```bash
cd website
python3 -m http.server 8080
```

Open http://localhost:8080

## Notes

- Next.js `/_next/image` optimizer URLs were rewritten to direct static media paths.
- Absolute `/_next/` asset paths in JS were made relative so the site works from any host root.
- Visual design, copy, animations, and assets match the original INTVL export.
