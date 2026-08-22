(function () {
  const cfg = window.TERRITORY_CONFIG || {};
  const H3_RES = cfg.H3_CAPTURE_RES || 10;
  const TILE_RES = cfg.H3_TILE_RES || 7;
  const GUEST_KEY = "territory_run_guest_id";

  let map, routeLayer, userMarker, territoryLayer;
  let guestId = null;
  let watchId = null;
  let recording = false;
  let paused = false;
  let activityId = null;
  let points = [];
  let capturedCells = new Set();
  let lastPoint = null;
  let distanceM = 0;
  let startTime = null;
  let loadedTiles = new Set();

  const $ = (id) => document.getElementById(id);

  function apiBase() {
    return (cfg.API_BASE_URL || "").replace(/\/$/, "");
  }

  function setStatus(el, text, type) {
    el.textContent = text;
    el.className = "status-bar" + (type ? " " + type : "");
  }

  function getOrCreateGuestId() {
    let id = localStorage.getItem(GUEST_KEY);
    if (!id) {
      id = crypto.randomUUID();
      localStorage.setItem(GUEST_KEY, id);
    }
    return id;
  }

  function newGuestId() {
    guestId = crypto.randomUUID();
    localStorage.setItem(GUEST_KEY, guestId);
    return guestId;
  }

  function haversineM(lat1, lon1, lat2, lon2) {
    const R = 6371000;
    const toRad = (d) => (d * Math.PI) / 180;
    const dLat = toRad(lat2 - lat1);
    const dLon = toRad(lon2 - lon1);
    const a =
      Math.sin(dLat / 2) ** 2 +
      Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLon / 2) ** 2;
    return 2 * R * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  }

  function formatPace(secPerKm) {
    if (!secPerKm || secPerKm > 3600) return "—";
    const m = Math.floor(secPerKm / 60);
    const s = Math.round(secPerKm % 60);
    return m + ":" + String(s).padStart(2, "0") + "/km";
  }

  async function api(path, options = {}) {
    if (!guestId) throw new Error("Not signed in as guest");
    const res = await fetch(apiBase() + path, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        "X-Guest-User": guestId,
        ...(options.headers || {}),
      },
    });
    if (!res.ok) {
      const err = await res.text();
      throw new Error(err || "HTTP " + res.status);
    }
    return res.json();
  }

  function initMap(lat, lon) {
    if (map) {
      map.setView([lat, lon], 15);
      return;
    }
    map = L.map("map", { zoomControl: true }).setView([lat, lon], 15);
    L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
      maxZoom: 19,
      attribution: "© OpenStreetMap",
    }).addTo(map);

    routeLayer = L.polyline([], { color: "#3dd6c3", weight: 4 }).addTo(map);
    userMarker = L.circleMarker([lat, lon], {
      radius: 8,
      color: "#fff",
      fillColor: "#3dd6c3",
      fillOpacity: 1,
      weight: 2,
    }).addTo(map);

    territoryLayer = L.geoJSON(null, {
      style: (feature) => {
        const owner = feature.properties?.owner || "";
        const isYou = owner && owner === guestId;
        return {
          color: isYou ? "#3dd6c3" : "#ff8c5a",
          weight: 1,
          fillColor: isYou ? "#3dd6c3" : "#ff8c5a",
          fillOpacity: isYou ? 0.35 : 0.25,
        };
      },
      onEachFeature: (feature, layer) => {
        const p = feature.properties || {};
        layer.bindPopup(
          "<div class='territory-popup'><strong>Territory</strong><br/>Owner: " +
            (p.owner || "none") +
            "</div>"
        );
      },
    }).addTo(map);

    map.on("moveend", () => debounceLoadTiles());
    loadTilesForCenter();
  }

  let tileTimer;
  function debounceLoadTiles() {
    clearTimeout(tileTimer);
    tileTimer = setTimeout(loadTilesForCenter, 400);
  }

  async function loadTilesForCenter() {
    if (!map || !guestId) return;
    const center = map.getCenter();
    const tileHex = h3.cellToString(h3.latLngToCell([center.lat, center.lng], TILE_RES));
    if (loadedTiles.has(tileHex)) return;
    loadedTiles.add(tileHex);
    try {
      const geo = await api("/api/v1/territories/tile/" + tileHex);
      territoryLayer.addData(geo);
    } catch (e) {
      console.warn("tile load", e);
    }
  }

  async function refreshProfile() {
    try {
      const p = await api("/api/v1/users/me");
      $("headerStats").textContent =
        "Cells: " + p.cells_owned +
        " · Score: " + p.total_capture_score +
        " · Stolen: " + (p.territories_stolen || 0);
    } catch (e) {
      console.warn("profile", e);
    }
  }

  function processGpsPosition(pos) {
    if (!recording || paused) return;
    const lat = pos.coords.latitude;
    const lon = pos.coords.longitude;
    const acc = pos.coords.accuracy;
    if (acc > 40) return;

    const ts = Math.floor(pos.timestamp / 1000);
    const point = { lat, lon, ts, acc, speed: pos.coords.speed || 0 };

    if (lastPoint) {
      const dt = ts - lastPoint.ts;
      const dist = haversineM(lastPoint.lat, lastPoint.lon, lat, lon);
      if (dt > 0 && dist / dt > 7) return;
      if (dist < 3) return;
      distanceM += dist;
    }

    lastPoint = point;
    points.push(point);

    const cell = h3.latLngToCell([lat, lon], H3_RES);
    capturedCells.add(h3.cellToString(cell));

    const latlngs = routeLayer.getLatLngs();
    latlngs.push([lat, lon]);
    routeLayer.setLatLngs(latlngs);
    userMarker.setLatLng([lat, lon]);

    $("metricDistance").textContent = (distanceM / 1000).toFixed(2);
    $("metricCells").textContent = capturedCells.size;
    if (startTime && distanceM > 100) {
      const elapsed = (Date.now() - startTime) / 1000;
      $("metricPace").textContent = formatPace(elapsed / (distanceM / 1000));
    }
  }

  function startGpsWatch() {
    if (!navigator.geolocation) {
      setStatus($("runStatus"), "GPS not supported in this browser.", "error");
      return;
    }
    watchId = navigator.geolocation.watchPosition(
      processGpsPosition,
      (err) => setStatus($("runStatus"), "GPS error: " + err.message, "error"),
      { enableHighAccuracy: true, maximumAge: 2000, timeout: 15000 }
    );
  }

  function stopGpsWatch() {
    if (watchId != null) {
      navigator.geolocation.clearWatch(watchId);
      watchId = null;
    }
  }

  async function startRun() {
    activityId = crypto.randomUUID();
    points = [];
    capturedCells = new Set();
    distanceM = 0;
    lastPoint = null;
    startTime = Date.now();
    recording = true;
    paused = false;
    routeLayer.setLatLngs([]);

    $("btnStart").classList.add("hidden");
    $("btnPause").classList.remove("hidden");
    $("btnFinish").classList.remove("hidden");
    setStatus($("runStatus"), "Recording… move outdoors to capture cells.", "ok");

    try {
      const resp = await api("/api/v1/activities", {
        method: "POST",
        body: JSON.stringify({
          started_at: new Date().toISOString(),
          idempotency_key: activityId,
          device_info: { platform: "web-guest", userAgent: navigator.userAgent },
        }),
      });
      if (resp.activity_id) activityId = resp.activity_id;
    } catch (e) {
      setStatus($("runStatus"), "API error: " + e.message, "error");
    }

    startGpsWatch();
  }

  async function finishRun() {
    recording = false;
    stopGpsWatch();
    $("btnPause").classList.add("hidden");
    $("btnFinish").classList.add("hidden");
    $("btnStart").classList.remove("hidden");
    setStatus($("runStatus"), "Syncing…");

    try {
      if (points.length > 0) {
        await api("/api/v1/activities/" + activityId + "/points", {
          method: "POST",
          body: JSON.stringify({
            points: points.map((p) => ({
              lat: p.lat, lon: p.lon, ts: p.ts, acc: p.acc, speed: p.speed,
            })),
          }),
        });
      }
      const duration = points.length >= 2
        ? points[points.length - 1].ts - points[0].ts
        : 0;
      const result = await api("/api/v1/activities/" + activityId + "/complete", {
        method: "POST",
        body: JSON.stringify({
          ended_at: new Date().toISOString(),
          distance_m: distanceM,
          duration_s: duration,
          h3_cells_hint: Array.from(capturedCells),
        }),
      });

      const cr = result.capture_result || {};
      const newCells = cr.new_cells ?? 0;
      const stolen = cr.stolen_cells ?? 0;
      const defended = cr.defended_cells ?? 0;
      const score = cr.capture_score ?? 0;
      setStatus(
        $("runStatus"),
        "Claimed! +" + newCells + " new · " + stolen + " stolen · " + result.status,
        "ok"
      );
      $("sumNew").textContent = String(newCells);
      $("sumStolen").textContent = String(stolen);
      $("sumDefended").textContent = String(defended);
      $("sumScore").textContent = String(score);
      $("summaryStatus").textContent = result.status === "verified"
        ? "Territory updated on the map."
        : (result.message || result.status);
      $("summarySheet").classList.remove("hidden");

      loadedTiles.clear();
      territoryLayer.clearLayers();
      await refreshProfile();
      await loadTilesForCenter();
    } catch (e) {
      setStatus($("runStatus"), "Sync failed: " + e.message, "error");
    }
  }

  async function enterApp() {
    guestId = getOrCreateGuestId();
    $("loginScreen").classList.add("hidden");
    $("loginScreen").style.display = "none";
    $("app").classList.remove("hidden");
    $("app").style.display = "flex";

    await refreshProfile();

    navigator.geolocation.getCurrentPosition(
      (pos) => initMap(pos.coords.latitude, pos.coords.longitude),
      () => initMap(28.6139, 77.2090),
      { enableHighAccuracy: true }
    );
  }

  function boot() {
    if (!cfg.API_BASE_URL || cfg.API_BASE_URL.includes("YOUR-SERVICE")) {
      setStatus($("loginStatus"), "Set API_BASE_URL in js/config.js", "error");
      return;
    }
    setStatus($("loginStatus"), "Tap Continue as Guest to open the map.", "ok");

    // Auto-enter if guest already exists
    if (localStorage.getItem(GUEST_KEY)) {
      enterApp();
    }
  }

  async function loadLeaderboard(kind) {
    const list = $("leaderboardList");
    list.innerHTML = "";
    setStatus($("boardStatus"), "Loading…");
    try {
      let path = "/api/v1/leaderboards/global?limit=25";
      if (kind === "local" && map) {
        const c = map.getCenter();
        const tileHex = h3.cellToString(h3.latLngToCell([c.lat, c.lng], TILE_RES));
        path = "/api/v1/leaderboards/local?h3_parent=" + encodeURIComponent(tileHex) + "&limit=25";
      }
      const data = await api(path);
      const entries = data.entries || [];
      if (!entries.length) {
        setStatus($("boardStatus"), "No ranks yet — finish a verified run.");
        return;
      }
      entries.forEach((e) => {
        const li = document.createElement("li");
        if (e.user_id === guestId) li.classList.add("you");
        li.innerHTML =
          "<span class='rank'>#" + e.rank + "</span>" +
          "<span>" + (e.display_name || e.username) + "</span>" +
          "<span>" + e.total_capture_score + " pts · " + e.cells_owned + " cells</span>";
        list.appendChild(li);
      });
      setStatus($("boardStatus"), data.scope || kind, "ok");
    } catch (e) {
      setStatus($("boardStatus"), "Could not load ranks: " + e.message, "error");
    }
  }

  $("btnGuestLogin").addEventListener("click", () => enterApp());
  $("btnNewGuest").addEventListener("click", () => {
    stopGpsWatch();
    newGuestId();
    location.reload();
  });
  $("btnStart").addEventListener("click", startRun);
  $("btnPause").addEventListener("click", () => {
    paused = !paused;
    $("btnPause").textContent = paused ? "Resume" : "Pause";
    setStatus($("runStatus"), paused ? "Paused" : "Recording…", paused ? "" : "ok");
  });
  $("btnFinish").addEventListener("click", finishRun);
  $("btnLeaderboard").addEventListener("click", () => {
    $("leaderboardSheet").classList.remove("hidden");
    $("tabGlobal").classList.add("active");
    $("tabLocal").classList.remove("active");
    loadLeaderboard("global");
  });
  $("btnCloseLeaderboard").addEventListener("click", () => {
    $("leaderboardSheet").classList.add("hidden");
  });
  $("tabGlobal").addEventListener("click", () => {
    $("tabGlobal").classList.add("active");
    $("tabLocal").classList.remove("active");
    loadLeaderboard("global");
  });
  $("tabLocal").addEventListener("click", () => {
    $("tabLocal").classList.add("active");
    $("tabGlobal").classList.remove("active");
    loadLeaderboard("local");
  });
  $("btnCloseSummary").addEventListener("click", () => {
    $("summarySheet").classList.add("hidden");
  });
  $("btnRefreshTiles").addEventListener("click", () => {
    loadedTiles.clear();
    territoryLayer.clearLayers();
    loadTilesForCenter();
  });

  boot();
})();
