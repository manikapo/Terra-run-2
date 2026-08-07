(function () {
  function bridge(name) {
    if (window.TerritoryBridge && window.TerritoryBridge[name]) {
      return window.TerritoryBridge[name]();
    }
    return null;
  }

  async function loadProfile() {
    const token = bridge('getAccessToken');
    const apiBase = bridge('getApiBaseUrl') || '';
    const statusEl = document.getElementById('apiStatus');

    if (!token || !apiBase) {
      statusEl.textContent = 'Bridge unavailable — open from app WebView';
      return;
    }

    try {
      const res = await fetch(apiBase.replace(/\/$/, '') + '/api/v1/users/me', {
        headers: { Authorization: 'Bearer ' + token }
      });
      if (!res.ok) throw new Error('HTTP ' + res.status);
      const data = await res.json();

      document.getElementById('displayName').textContent = data.display_name || 'Runner';
      document.getElementById('username').textContent = '@' + (data.username || 'user');
      document.getElementById('cellsOwned').textContent = data.cells_owned ?? 0;
      document.getElementById('captureScore').textContent = data.total_capture_score ?? 0;
      document.getElementById('territories').textContent = data.territories_captured ?? 0;
      document.getElementById('activities').textContent = data.activity_count ?? 0;
      statusEl.textContent = 'Synced from API';
    } catch (e) {
      statusEl.textContent = 'Failed: ' + e.message;
    }
  }

  loadProfile();
})();
