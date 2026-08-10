// Nav drawer
const navBtn = document.getElementById('navMenuBtn');
const mobileNav = document.getElementById('mobileNav');
navBtn?.addEventListener('click', () => {
  mobileNav?.classList.toggle('open');
  document.body.style.overflow = mobileNav?.classList.contains('open') ? 'hidden' : '';
});
mobileNav?.querySelectorAll('a').forEach((a) => {
  a.addEventListener('click', () => {
    mobileNav?.classList.remove('open');
    document.body.style.overflow = '';
  });
});

// Modal
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
document.addEventListener('keydown', (e) => { if (e.key === 'Escape') closeDownloadModal(); });
window.openDownloadModal = openDownloadModal;
window.closeDownloadModal = closeDownloadModal;

// FAQ
document.querySelectorAll('.faq-btn').forEach((btn) => {
  btn.addEventListener('click', () => {
    const item = btn.parentElement;
    const open = item.classList.contains('open');
    document.querySelectorAll('.faq-item').forEach((i) => i.classList.remove('open'));
    if (!open) item.classList.add('open');
  });
});

// Tint mobile frame SVG to orange via CSS filter on frame-mobile
const frameMobile = document.querySelector('.frame-mobile');
if (frameMobile) {
  frameMobile.style.filter = 'hue-rotate(85deg) saturate(1.8) brightness(1.1)';
}
