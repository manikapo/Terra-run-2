/**
 * INTVL scroll system port for 8me.in
 * Matches: stack-card scale 0.94 / shade 0.45, conversion column, HIW scrub, intro glyph
 */
(function () {
  const reduced = matchMedia('(prefers-reduced-motion: reduce)').matches;
  const mobile = innerWidth < 900;

  if (reduced || typeof gsap === 'undefined' || typeof ScrollTrigger === 'undefined') {
    document.getElementById('introOverlay')?.classList.add('done');
    return;
  }

  gsap.registerPlugin(ScrollTrigger);

  // Lenis smooth scroll
  if (typeof Lenis !== 'undefined' && !mobile) {
    const lenis = new Lenis({ duration: 1.15, smoothWheel: true });
    lenis.on('scroll', ScrollTrigger.update);
    gsap.ticker.add((t) => lenis.raf(t * 1000));
    gsap.ticker.lagSmoothing(0);
  }

  // ── INTRO (intro-glyph scale punch) ──
  const intro = document.getElementById('introOverlay');
  if (intro) {
    gsap.set(['.intro-glyph', '.intro-hole'], { transformOrigin: 'center center', scale: 0 });
    gsap.timeline({
      delay: 0.2,
      onComplete: () => intro.classList.add('done'),
    })
      .to('.intro-glyph', { scale: 1, duration: 0.45, ease: 'power3.out' })
      .to('.intro-glyph', { scale: 28, duration: 0.75, ease: 'power2.in' }, 0.35)
      .to('.intro-hole', { scale: 28, duration: 0.75, ease: 'power2.in' }, 0.35);
  }

  // Nav scroll state
  ScrollTrigger.create({
    start: 80,
    end: 999999,
    onUpdate: (self) => {
      document.querySelector('.site-nav')?.classList.toggle('scrolled', self.scroll() > 60);
    },
  });

  // ── STACK CARDS: INTVL exact values scale 0.94, shade opacity 0.45 ──
  const stackCards = gsap.utils.toArray('.stack-card');
  stackCards.forEach((card, i) => {
    const next = stackCards[i + 1];
    if (!next) return;
    const inner = card.querySelector('.stack-card-inner');
    const shade = card.querySelector('.stack-card-shade');
    const tl = gsap.timeline({
      scrollTrigger: {
        trigger: next,
        start: 'top bottom',
        end: 'top top',
        scrub: true,
      },
    });
    if (inner) tl.to(inner, { scale: 0.94, ease: 'none' }, 0);
    if (shade) tl.to(shade, { opacity: 0.45, ease: 'none' }, 0);
  });

  // ── HERO parallax ──
  const globe = document.querySelector('.hero-globe');
  if (globe && stackCards[0]) {
    gsap.to(globe, {
      y: '20%',
      scale: 1.12,
      ease: 'none',
      scrollTrigger: {
        trigger: stackCards[0],
        start: 'top top',
        end: 'bottom top',
        scrub: true,
      },
    });
    gsap.to('.hero-content', {
      y: -60,
      opacity: 0,
      ease: 'none',
      scrollTrigger: {
        trigger: stackCards[0],
        start: 'top top',
        end: '70% top',
        scrub: true,
      },
    });
  }

  // ── CONVERSION COLUMN: phone screen crossfade + parallax ──
  const conversionWrap = document.querySelector('[data-conversion-wrap]');
  const conversionCol = document.querySelector('[data-conversion-column]');
  const convScreens = document.querySelectorAll('.phone-screens .phone-screen');

  if (conversionWrap && conversionCol) {
    // Column translate on scroll (parallax lift)
    gsap.to(conversionCol, {
      y: () => -(conversionCol.offsetHeight - innerHeight * 0.5),
      ease: 'none',
      scrollTrigger: {
        trigger: conversionWrap,
        start: 'top bottom',
        end: 'bottom top',
        scrub: true,
      },
    });

    if (convScreens.length >= 2) {
      ScrollTrigger.create({
        trigger: conversionWrap,
        start: 'top center',
        end: 'bottom center',
        scrub: 0.5,
        onUpdate: (self) => {
          const p = self.progress;
          convScreens[0].style.opacity = p < 0.5 ? '1' : String(1 - (p - 0.5) * 2);
          convScreens[1].style.opacity = p < 0.5 ? String(p * 2) : '1';
        },
      });
    }
  }

  // ── STATS pills entrance ──
  gsap.from('.stats-pills', {
    y: 50,
    opacity: 0,
    scale: 0.92,
    duration: 0.7,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: stackCards[2] || '.stats-scene',
      start: 'top 65%',
      toggleActions: 'play none none reverse',
    },
  });

  // ── HOW IT WORKS scroll scrub ──
  const hiwScrub = document.querySelector('.hiw-scrub');
  if (hiwScrub && !mobile) {
    const steps = hiwScrub.querySelectorAll('.hiw-step');
    const screens = hiwScrub.querySelectorAll('.hiw-screen');
    let cur = -1;

    const setStep = (n) => {
      if (n === cur) return;
      cur = n;
      steps.forEach((s, i) => {
        s.classList.toggle('active', i === n);
        s.classList.toggle('past', i < n);
      });
    };

    ScrollTrigger.create({
      trigger: hiwScrub,
      start: 'top top',
      end: 'bottom bottom',
      scrub: 0.35,
      onUpdate: (self) => {
        const p = self.progress;
        const step = p < 0.33 ? 0 : p < 0.66 ? 1 : 2;
        setStep(step);

        screens.forEach((screen, i) => {
          const segStart = i / 3;
          const segEnd = (i + 1) / 3;
          if (p < segStart) screen.style.opacity = '0';
          else if (p > segEnd) screen.style.opacity = i === screens.length - 1 ? '1' : '0';
          else {
            const local = (p - segStart) / (segEnd - segStart);
            screen.style.opacity = String(1 - local);
            if (screens[i + 1]) screens[i + 1].style.opacity = String(local);
          }
        });
      },
    });
    setStep(0);
    if (screens[0]) screens[0].style.opacity = '1';
  }

  // ── Feature tiles stagger ──
  gsap.from('.feature-tile', {
    y: 32,
    opacity: 0,
    duration: 0.55,
    stagger: 0.08,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: '.feature-grid',
      start: 'top 82%',
      toggleActions: 'play none none reverse',
    },
  });

  // ── CTA reveal ──
  gsap.from('.cta-inner', {
    scale: 0.92,
    opacity: 0,
    duration: 0.7,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: '.cta-scene',
      start: 'top 78%',
      toggleActions: 'play none none reverse',
    },
  });

})();
