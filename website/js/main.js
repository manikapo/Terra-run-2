// Mobile menu
const navBtn = document.getElementById('navMenuBtn');
const mobileNav = document.getElementById('mobileNav');

navBtn?.addEventListener('click', () => {
  navBtn.classList.toggle('open');
  mobileNav?.classList.toggle('open');
  document.body.style.overflow = mobileNav?.classList.contains('open') ? 'hidden' : '';
});

mobileNav?.querySelectorAll('a').forEach((link) => {
  link.addEventListener('click', () => {
    navBtn?.classList.remove('open');
    mobileNav?.classList.remove('open');
    document.body.style.overflow = '';
  });
});

// Download modal
function openDownloadModal() {
  document.getElementById('downloadModal')?.classList.add('open');
  document.body.style.overflow = 'hidden';
}

function closeDownloadModal() {
  document.getElementById('downloadModal')?.classList.remove('open');
  document.body.style.overflow = '';
}

document.getElementById('downloadModal')?.addEventListener('click', (e) => {
  if (e.target.id === 'downloadModal') closeDownloadModal();
});

document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') closeDownloadModal();
});

window.openDownloadModal = openDownloadModal;
window.closeDownloadModal = closeDownloadModal;

// FAQ accordion
document.querySelectorAll('.faq-q').forEach((btn) => {
  btn.addEventListener('click', () => {
    const item = btn.parentElement;
    const wasOpen = item.classList.contains('open');
    document.querySelectorAll('.faq-item').forEach((i) => i.classList.remove('open'));
    if (!wasOpen) item.classList.add('open');
  });
});

// Fallback: populate phone screens if GSAP/scroll.js didn't run
if (!window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
  document.addEventListener('DOMContentLoaded', () => {
    setTimeout(() => {
      if (window.phoneScreens) {
        document.querySelectorAll('[data-screen="route"]').forEach((el) => {
          if (!el.innerHTML.trim()) el.innerHTML = window.phoneScreens.route();
        });
        document.querySelectorAll('[data-screen="territory"]').forEach((el) => {
          if (!el.innerHTML.trim()) el.innerHTML = window.phoneScreens.territory();
        });
        document.querySelectorAll('[data-screen="leaderboard"]').forEach((el) => {
          if (!el.innerHTML.trim()) el.innerHTML = window.phoneScreens.leaderboard();
        });
      }
    }, 100);
  });
} else {
  document.documentElement.classList.add('no-motion');
  if (window.phoneScreens) {
    document.querySelectorAll('[data-screen="route"]').forEach((el) => { el.innerHTML = window.phoneScreens.route(); });
    document.querySelectorAll('[data-screen="territory"]').forEach((el) => { el.innerHTML = window.phoneScreens.territory(); });
    document.querySelectorAll('[data-screen="leaderboard"]').forEach((el) => { el.innerHTML = window.phoneScreens.leaderboard(); });
  }
}
