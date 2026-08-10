(function () {
  function openDownloadModal() {
    var modal = document.getElementById('downloadModal');
    if (modal) modal.classList.add('open');
  }
  function closeDownloadModal() {
    var modal = document.getElementById('downloadModal');
    if (modal) modal.classList.remove('open');
  }
  window.openDownloadModal = openDownloadModal;
  window.closeDownloadModal = closeDownloadModal;

  document.addEventListener('click', function (e) {
    var modal = document.getElementById('downloadModal');
    if (modal && e.target === modal) closeDownloadModal();
  });
  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') closeDownloadModal();
  });

  // Wire Android CTAs to modal instead of Play Store
  document.querySelectorAll('.hero-cta-play a, a[href*="play.google.com/store/apps/details?id=com.intvl"]').forEach(function (el) {
    el.addEventListener('click', function (e) {
      e.preventDefault();
      openDownloadModal();
    });
  });

  // Update iOS CTAs to play in browser
  document.querySelectorAll('.hero-cta-ios a, a[href*="apps.apple.com/au/app/intvl"]').forEach(function (el) {
    el.href = 'https://play.8me.in';
    el.target = '_blank';
    el.rel = 'noopener noreferrer';
    var label = el.querySelector('span') || el;
    if (label.textContent && label.textContent.indexOf('iOS') !== -1) {
      label.textContent = 'Play in browser';
    }
  });
})();
