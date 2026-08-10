# 8me.in Landing Page

Static HTML port of the **INTVL** (intvl.com.au) Next.js landing page, adapted for **Infinite Me** with orange (`#FF4D00`) branding.

Based on the exported INTVL source structure: stack cards, conversion column, scroll-driven "How it works", intro splash, viewport frame.

## Deploy to Hostinger

Upload the entire `website/` folder contents to your **8me.in** root:

```
index.html
css/styles.css
js/main.js
js/scroll.js
assets/          ← images from INTVL (globe, phone, stats, etc.)
```

Keep existing `tutorial.html` and `privacy.html`.

## Preview locally

```bash
cd website && python3 -m http.server 8765
```

Open http://localhost:8765 — **must use a local server** (not file://).

## INTVL features ported

| Feature | Implementation |
|---------|----------------|
| Intro glyph punch | GSAP scale animation on orange bar |
| Viewport frame | `assets/frame-orange.svg` + mobile frame |
| Stack cards (×3) | Sticky sections scale to **0.94**, shade **0.45** |
| Conversion column | Tall scroll area with phone hand + screen crossfade |
| Hero globe | `assets/heroGlobe.c6079214.webp` + parallax |
| How it works | `350vh+` scroll scrub, sticky phone, 3 screen states |
| Stats pills | White card on cyclist photo |
| Features grid | Runner photo background |
| Integrations bar | Orange marquee |
| FAQ accordion | Orange buttons |
| Smooth scroll | Lenis + GSAP ScrollTrigger |

## Color mapping

| INTVL token | 8me.in value |
|-------------|--------------|
| signal-400 | `#FF4D00` |
| ground-800 | `#070b06` |
| ground-600 | `#192616` |
