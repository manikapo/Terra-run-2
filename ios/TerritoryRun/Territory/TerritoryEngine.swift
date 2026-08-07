import Foundation

/// Phase 1: client hints cells; server recomputes H3 from polyline.
/// Replace with H3 Swift (e.g. uber h3) before production map overlay.
final class TerritoryEngine {
    private var captured: Set<String> = []

    func reset() { captured.removeAll() }

    func process(lat: Double, lon: Double) -> String? {
        // ~66m grid hint at equator (server validates with real H3)
        let latKey = Int(lat * 1000)
        let lonKey = Int(lon * 1000)
        let hex = "\(latKey)_\(lonKey)"
        if captured.insert(hex).inserted { return hex }
        return nil
    }

    func cells() -> [String] { Array(captured) }
    var count: Int { captured.count }
}
