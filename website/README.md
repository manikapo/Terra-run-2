# 8me.in Landing Page

Marketing site for Infinite Me — INTVL-inspired design with 8me.in orange branding.

## Deploy to Hostinger

Upload to your **8me.in** site root:

```
index.html
css/styles.css
js/main.js
js/scroll.js
js/phone-screens.js
```

Keep existing `tutorial.html` and `privacy.html`.

## Preview locally

```bash
cd website && python3 -m http.server 8765
```

Open http://localhost:8765 (use a local server — scroll effects need HTTP, not file://).

## Scroll effects (INTVL-style)

- **Stack cards** — Hero, showcase, and stats stick and stack; previous cards scale down with shade
- **Hero parallax** — Map background drifts; headline fades on scroll
- **Wave dividers** — Curved SVG transitions between sections
- **Showcase phone** — Screen crossfades route tracking → territory map while scrolling
- **How it works** — 350vh scroll scrub: sticky phone + step highlights + 3 phone screen states
- **Smooth scroll** — Lenis (desktop)
- **Intro splash** — Orange bar zoom on load
- **GSAP ScrollTrigger** — Stats entrance, feature stagger, CTA scale

CDN deps: Lenis, GSAP, ScrollTrigger

Mobile: stack scaling off; HIW uses simple vertical steps.

Accent color: `#FF4D00` (8me.in orange, replacing INTVL's green).
