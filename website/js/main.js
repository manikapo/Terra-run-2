// Mobile menu
const navToggle = document.getElementById('navToggle');
const mobileMenu = document.getElementById('mobileMenu');

if (navToggle && mobileMenu) {
  navToggle.addEventListener('click', () => {
    navToggle.classList.toggle('open');
    mobileMenu.classList.toggle('open');
  });
  mobileMenu.querySelectorAll('a').forEach((link) => {
    link.addEventListener('click', () => {
      navToggle.classList.remove('open');
      mobileMenu.classList.remove('open');
    });
  });
}

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

// FAQ accordion
document.querySelectorAll('.faq-q').forEach((q) => {
  q.addEventListener('click', () => {
    const item = q.parentElement;
    const wasOpen = item.classList.contains('open');
    document.querySelectorAll('.faq-item').forEach((i) => i.classList.remove('open'));
    if (!wasOpen) item.classList.add('open');
  });
});

// Scroll reveal
const revealEls = document.querySelectorAll('.reveal');
const revealObserver = new IntersectionObserver(
  (entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        entry.target.classList.add('in');
        revealObserver.unobserve(entry.target);
      }
    });
  },
  { threshold: 0.12, rootMargin: '0px 0px -40px 0px' }
);
revealEls.forEach((el) => revealObserver.observe(el));

// Expose modal fn for inline onclick
window.openDownloadModal = openDownloadModal;
window.closeDownloadModal = closeDownloadModal;
