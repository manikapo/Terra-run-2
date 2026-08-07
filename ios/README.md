# iOS setup

1. Xcode → File → New → Project → App (SwiftUI, Swift)
2. Product name: `TerritoryRun`
3. Copy all files from `ios/TerritoryRun/` into the Xcode project
4. Add Swift Package dependencies:
   - `https://github.com/supabase/supabase-swift` (Auth)
   - `https://github.com/mapbox/mapbox-maps-ios` (optional Phase 1 map)
   - H3: add `CH3` or compute cells server-side only for minimal iOS build

5. Info.plist keys (or Build Settings → User-Defined):
   - `API_BASE_URL`
   - `SUPABASE_URL`
   - `SUPABASE_ANON_KEY`
   - `OTA_MANIFEST_URL`

6. Capabilities:
   - Background Modes → Location updates
   - Sign in with Apple / Google URL schemes per Supabase docs

7. `NSLocationWhenInUseUsageDescription` and `NSLocationAlwaysAndWhenInUseUsageDescription`

Note: `TerritoryEngine.swift` uses `CH3` — replace with server-only capture or add H3 Swift package before building.
