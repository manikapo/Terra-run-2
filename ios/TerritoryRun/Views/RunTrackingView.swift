import SwiftUI
import CoreLocation

struct RunTrackingView: View {
    @Environment(\.dismiss) private var dismiss
    @StateObject private var gps = GpsTracker()
    @State private var territory = TerritoryEngine()
    @State private var activityId = UUID().uuidString
    @State private var points: [GpsPoint] = []
    @State private var distanceM: Double = 0
    @State private var recording = false
    @State private var status = "Ready"

    var body: some View {
        VStack(spacing: 16) {
            Text("Distance: \(String(format: "%.2f", distanceM / 1000)) km")
            Text("Cells: \(territory.count)")
            Text(status)
            if !recording {
                Button("Start") {
                    gps.requestPermission()
                    territory.reset()
                    points = []
                    distanceM = 0
                    recording = true
                    status = "Recording"
                    Task {
                        do {
                            activityId = try await ApiClient.createActivity(idempotencyKey: activityId)
                            await trackLoop()
                        } catch {
                            status = error.localizedDescription
                        }
                    }
                }
            } else {
                Button("Finish & Sync") {
                    recording = false
                    gps.stop()
                    Task { await sync() }
                }
            }
        }
        .padding()
    }

    private func trackLoop() async {
        for await loc in gps.locations() {
            guard recording else { break }
            if let last = points.last {
                distanceM += haversine(lat1: last.lat, lon1: last.lon, lat2: loc.coordinate.latitude, lon2: loc.coordinate.longitude)
            }
            territory.process(lat: loc.coordinate.latitude, lon: loc.coordinate.longitude)
            points.append(GpsPoint(
                lat: loc.coordinate.latitude,
                lon: loc.coordinate.longitude,
                ts: Int64(loc.timestamp.timeIntervalSince1970),
                acc: loc.horizontalAccuracy,
                speed: max(loc.speed, 0)
            ))
        }
    }

    private func sync() async {
        status = "Syncing…"
        do {
            try await ApiClient.uploadPoints(activityId: activityId, points: points)
            let duration = points.count >= 2 ? Int(points.last!.ts - points.first!.ts) : 0
            let resp = try await ApiClient.completeActivity(
                activityId: activityId,
                distanceM: distanceM,
                durationS: duration,
                cells: territory.cells()
            )
            status = "Done — \(resp.capture_result?.new_cells ?? 0) cells"
            dismiss()
        } catch {
            status = error.localizedDescription
        }
    }

    private func haversine(lat1: Double, lon1: Double, lat2: Double, lon2: Double) -> Double {
        let r = 6371000.0
        let dLat = (lat2 - lat1) * .pi / 180
        let dLon = (lon2 - lon1) * .pi / 180
        let a = sin(dLat/2)*sin(dLat/2) + cos(lat1 * .pi/180)*cos(lat2 * .pi/180)*sin(dLon/2)*sin(dLon/2)
        return 2 * r * atan2(sqrt(a), sqrt(1-a))
    }
}
