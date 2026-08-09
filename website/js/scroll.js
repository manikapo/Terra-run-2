/* Scroll animations — Lenis + GSAP ScrollTrigger */

(function initScroll() {
  // Populate phone screens first
  if (window.phoneScreens) {
    document.querySelectorAll('[data-screen="route"]').forEach((el) => { el.innerHTML = window.phoneScreens.route(); });
    document.querySelectorAll('[data-screen="territory"]').forEach((el) => { el.innerHTML = window.phoneScreens.territory(); });
    document.querySelectorAll('[data-screen="leaderboard"]').forEach((el) => { el.innerHTML = window.phoneScreens.leaderboard(); });
  }

  const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  const isMobile = window.innerWidth < 900;

  if (reduced || typeof gsap === 'undefined' || typeof ScrollTrigger === 'undefined') {
    document.documentElement.classList.add('no-motion');
    return;
  }

  gsap.registerPlugin(ScrollTrigger);

  // Lenis smooth scroll
  let lenis;
  if (typeof Lenis !== 'undefined' && !isMobile) {
    lenis = new Lenis({ duration: 1.2, easing: (t) => Math.min(1, 1.001 - Math.pow(2, -10 * t)), smoothWheel: true });
    lenis.on('scroll', ScrollTrigger.update);
    gsap.ticker.add((time) => lenis.raf(time * 1000));
    gsap.ticker.lagSmoothing(0);
  }

  // Intro splash
  const intro = document.getElementById('intro');
  if (intro) {
    gsap.timeline({ delay: 0.3 })
      .to('.intro-glyph', { scale: 1, duration: 0.5, ease: 'power3.out' })
      .to('.intro-glyph', { scale: 40, duration: 0.8, ease: 'power2.in' })
      .to(intro, { opacity: 0, duration: 0.4, onComplete: () => { intro.style.display = 'none'; } }, '-=0.2');
  }

  // Nav background on scroll
  const nav = document.querySelector('.nav');
  ScrollTrigger.create({
    start: 'top -80',
    end: 99999,
    onUpdate: (self) => {
      nav?.classList.toggle('nav-scrolled', self.scroll() > 60);
    },
  });

  // ── STACK CARDS: scale previous card as next stacks on top ──
  const stackCards = document.querySelectorAll('.stack-card');
  stackCards.forEach((card, i) => {
    if (i === stackCards.length - 1) return;
    const inner = card.querySelector('.stack-card-inner');
    const next = stackCards[i + 1];
    if (!inner || !next) return;

    gsap.to(inner, {
      scale: 0.9,
      borderRadius: '24px',
      ease: 'none',
      scrollTrigger: {
        trigger: next,
        start: 'top bottom',
        end: 'top top',
        scrub: 0.6,
      },
    });

    const shade = card.querySelector('.stack-card-shade');
    if (shade) {
      gsap.to(shade, {
        opacity: 0.55,
        ease: 'none',
        scrollTrigger: {
          trigger: next,
          start: 'top bottom',
          end: 'top top',
          scrub: 0.6,
        },
      });
    }
  });

  // ── HERO parallax ──
  const heroBg = document.querySelector('.hero-bg');
  if (heroBg) {
    gsap.to(heroBg, {
      y: '18%',
      scale: 1.1,
      ease: 'none',
      scrollTrigger: {
        trigger: '.stack-card-hero',
        start: 'top top',
        end: 'bottom top',
        scrub: true,
      },
    });
  }

  const heroContent = document.querySelector('.hero-content');
  if (heroContent) {
    gsap.to(heroContent, {
      y: -80,
      opacity: 0,
      ease: 'none',
      scrollTrigger: {
        trigger: '.stack-card-hero',
        start: 'top top',
        end: '80% top',
        scrub: true,
      },
    });
  }

  // ── SHOWCASE phone screen crossfade on scroll ──
  const showcaseScreens = document.querySelectorAll('.showcase-screen');
  if (showcaseScreens.length >= 2) {
    ScrollTrigger.create({
      trigger: '.stack-card-showcase',
      start: 'top center',
      end: 'bottom center',
      scrub: 0.5,
      onUpdate: (self) => {
        const p = self.progress;
        showcaseScreens[0].style.opacity = p < 0.5 ? 1 : 1 - (p - 0.5) * 2;
        showcaseScreens[1].style.opacity = p < 0.5 ? p * 2 : 1;
      },
    });
  }

  // ── STATS card entrance ──
  const statsCard = document.querySelector('.stats-card');
  if (statsCard) {
    gsap.from(statsCard, {
      y: 60,
      opacity: 0,
      scale: 0.92,
      duration: 0.8,
      ease: 'power3.out',
      scrollTrigger: {
        trigger: '.stack-card-stats',
        start: 'top 70%',
        toggleActions: 'play none none reverse',
      },
    });
  }

  // ── HOW IT WORKS: scroll-driven steps + phone crossfade ──
  const hiwScrub = document.querySelector('.hiw-scrub');
  if (hiwScrub && !isMobile) {
    const steps = hiwScrub.querySelectorAll('.hiw-step');
    const screens = hiwScrub.querySelectorAll('.hiw-screen');
    const progressBar = hiwScrub.querySelector('.hiw-progress-fill');

    let currentStep = -1;
    function setHiwStep(index) {
      if (index === currentStep) return;
      currentStep = index;
      steps.forEach((s, i) => {
        s.classList.toggle('active', i === index);
        s.classList.toggle('past', i < index);
      });
      screens.forEach((s, i) => s.classList.toggle('active', i === index));
    }

    ScrollTrigger.create({
      trigger: hiwScrub,
      start: 'top top',
      end: 'bottom bottom',
      scrub: 0.4,
      onUpdate: (self) => {
        const p = self.progress;
        if (progressBar) progressBar.style.width = `${p * 100}%`;
        const step = p < 0.33 ? 0 : p < 0.66 ? 1 : 2;

        // Smooth crossfade between phone screens
        const seg = p * 3;
        const i = Math.min(2, Math.floor(seg));
        const local = seg - i;
        screens.forEach((screen, idx) => {
          if (idx === i) screen.style.opacity = String(1 - local * 0.8);
          else if (idx === i + 1) screen.style.opacity = String(local * 0.8);
          else screen.style.opacity = '0';
        });

        setHiwStep(step);
      },
    });

    setHiwStep(0);
  }

  // ── Feature cards stagger ──
  gsap.from('.feature-card', {
    y: 40,
    opacity: 0,
    duration: 0.6,
    stagger: 0.1,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: '.feature-grid',
      start: 'top 80%',
      toggleActions: 'play none none reverse',
    },
  });

  // ── Section reveals ──
  gsap.utils.toArray('.reveal').forEach((el) => {
    gsap.from(el, {
      y: 36,
      opacity: 0,
      duration: 0.7,
      ease: 'power3.out',
      scrollTrigger: {
        trigger: el,
        start: 'top 88%',
        toggleActions: 'play none none reverse',
      },
    });
  });

  // ── CTA scale in ──
  gsap.from('.cta-inner', {
    scale: 0.94,
    opacity: 0,
    duration: 0.8,
    ease: 'power3.out',
    scrollTrigger: {
      trigger: '.cta',
      start: 'top 75%',
      toggleActions: 'play none none reverse',
    },
  });

})();
