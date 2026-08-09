// Mobile nav
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

// How it works tabs
const howTabs = document.querySelectorAll('.how-tab');
const howPanels = document.querySelectorAll('.how-panel p');

howTabs.forEach((tab) => {
  tab.addEventListener('click', () => {
    const step = tab.dataset.step;
    howTabs.forEach((t) => t.classList.remove('active'));
    tab.classList.add('active');
    howPanels.forEach((p) => {
      p.classList.toggle('active', p.dataset.panel === step);
    });
  });
});

// FAQ accordion
document.querySelectorAll('.faq-q').forEach((btn) => {
  btn.addEventListener('click', () => {
    const item = btn.parentElement;
    const wasOpen = item.classList.contains('open');
    document.querySelectorAll('.faq-item').forEach((i) => i.classList.remove('open'));
    if (!wasOpen) item.classList.add('open');
  });
});

// Scroll reveal
const revealEls = document.querySelectorAll('.reveal');
const observer = new IntersectionObserver(
  (entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        entry.target.classList.add('in');
        observer.unobserve(entry.target);
      }
    });
  },
  { threshold: 0.1, rootMargin: '0px 0px -40px 0px' }
);
revealEls.forEach((el) => observer.observe(el));
