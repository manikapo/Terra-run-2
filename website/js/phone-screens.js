// Phone screen SVGs for scroll-driven states
function phoneScreenRoute() {
  return `<svg viewBox="0 0 300 600" xmlns="http://www.w3.org/2000/svg">
    <rect width="300" height="600" fill="#0a0e14"/>
    <defs><pattern id="g1" width="24" height="24" patternUnits="userSpaceOnUse"><path d="M24 0H0V24" fill="none" stroke="rgba(255,255,255,0.04)" stroke-width="0.5"/></pattern></defs>
    <rect width="300" height="600" fill="url(#g1)"/>
    <path d="M80 450 Q120 380 150 320 Q180 260 170 200 Q160 150 140 120" stroke="#FF4D00" stroke-width="3" fill="none" stroke-linecap="round"/>
    <circle cx="140" cy="120" r="8" fill="#FF4D00"/><circle cx="140" cy="120" r="18" fill="#FF4D00" opacity="0.2"/>
    <rect x="16" y="460" width="268" height="72" rx="12" fill="rgba(20,26,20,0.95)" stroke="rgba(255,255,255,0.08)"/>
    <text x="40" y="492" fill="#fff" font-family="Inter,sans-serif" font-size="11" font-weight="600">5.2 km</text>
    <text x="120" y="492" fill="#fff" font-family="Inter,sans-serif" font-size="11" font-weight="600">28:14</text>
    <text x="200" y="492" fill="#fff" font-family="Inter,sans-serif" font-size="11" font-weight="600">5:25/km</text>
    <rect x="40" y="508" width="220" height="14" rx="7" fill="#FF4D00"/>
    <text x="150" y="519" fill="#fff" font-family="Inter,sans-serif" font-size="9" font-weight="700" text-anchor="middle">PAUSE RUN</text>
    <rect x="0" y="0" width="300" height="44" fill="rgba(0,0,0,0.6)"/>
    <text x="150" y="28" fill="white" font-family="Inter,sans-serif" font-weight="700" font-size="11" text-anchor="middle">Tracking run…</text>
  </svg>`;
}

function phoneScreenTerritory() {
  return `<svg viewBox="0 0 300 600" xmlns="http://www.w3.org/2000/svg">
    <rect width="300" height="600" fill="#0a0e14"/>
    <polygon points="40,100 140,60 200,140 160,220 60,180" fill="rgba(255,77,0,0.45)" stroke="#FF4D00" stroke-width="1.5"/>
    <polygon points="120,80 240,70 260,180 180,240 100,200" fill="rgba(0,197,102,0.35)" stroke="#00C566" stroke-width="1.5"/>
    <polygon points="60,260 160,240 180,360 100,400 40,320" fill="rgba(0,87,255,0.3)" stroke="#0057FF" stroke-width="1.5"/>
    <polygon points="180,280 260,260 280,380 200,420" fill="rgba(255,77,0,0.25)" stroke="#FF4D00" stroke-width="1"/>
    <text x="100" y="150" fill="#FF4D00" font-family="Inter,sans-serif" font-weight="800" font-size="8" letter-spacing="1">YOUR ZONE</text>
    <text x="200" y="130" fill="#00C566" font-family="Inter,sans-serif" font-weight="800" font-size="8" letter-spacing="1">RIVAL</text>
    <rect x="0" y="0" width="300" height="44" fill="rgba(0,0,0,0.6)"/>
    <text x="150" y="28" fill="white" font-family="Inter,sans-serif" font-weight="700" font-size="11" text-anchor="middle">Territory map</text>
    <rect x="16" y="500" width="268" height="56" rx="10" fill="rgba(255,77,0,0.9)"/>
    <text x="150" y="533" fill="white" font-family="Inter,sans-serif" font-weight="700" font-size="12" text-anchor="middle">+1.4 km² claimed</text>
  </svg>`;
}

function phoneScreenLeaderboard() {
  return `<svg viewBox="0 0 300 600" xmlns="http://www.w3.org/2000/svg">
    <rect width="300" height="600" fill="#0f1410"/>
    <rect x="0" y="0" width="300" height="44" fill="rgba(0,0,0,0.6)"/>
    <text x="150" y="28" fill="white" font-family="Inter,sans-serif" font-weight="700" font-size="11" text-anchor="middle">Leaderboard</text>
    <text x="24" y="72" fill="#FF4D00" font-family="Inter,sans-serif" font-weight="800" font-size="13">TOP RUNNERS</text>
    ${[['01','RV','Rahul V.','4.2 km²'],['02','PS','Priya S.','3.7 km²'],['03','AK','Arjun K.','3.1 km²']].map(([r,av,n,a],i)=>`
      <rect x="16" y="${88+i*72}" width="268" height="60" rx="10" fill="rgba(255,255,255,0.04)" stroke="rgba(255,255,255,0.06)"/>
      <text x="32" y="${124+i*72}" fill="${i===0?'#FF4D00':'rgba(255,255,255,0.3)'}" font-family="Inter,sans-serif" font-weight="900" font-size="16">${r}</text>
      <circle cx="68" cy="${118+i*72}" r="16" fill="${i===0?'#FF4D00':'#333'}"/>
      <text x="68" y="${122+i*72}" fill="white" font-family="Inter,sans-serif" font-weight="700" font-size="9" text-anchor="middle">${av}</text>
      <text x="96" y="${114+i*72}" fill="white" font-family="Inter,sans-serif" font-weight="600" font-size="12">${n}</text>
      <text x="96" y="${130+i*72}" fill="rgba(255,255,255,0.4)" font-family="Inter,sans-serif" font-size="10">${a}</text>
    `).join('')}
    <rect x="16" y="520" width="268" height="48" rx="10" fill="rgba(255,77,0,0.15)" stroke="rgba(255,77,0,0.4)" stroke-dasharray="4 3"/>
    <text x="150" y="549" fill="#FF4D00" font-family="Inter,sans-serif" font-weight="700" font-size="11" text-anchor="middle">Your rank: —</text>
  </svg>`;
}

window.phoneScreens = { route: phoneScreenRoute, territory: phoneScreenTerritory, leaderboard: phoneScreenLeaderboard };
