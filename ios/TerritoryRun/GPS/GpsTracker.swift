import CoreLocation
import Combine

@MainActor
final class GpsTracker: NSObject, ObservableObject, CLLocationManagerDelegate {
    @Published var lastPoint: CLLocation?
    private let manager = CLLocationManager()
    private var continuation: AsyncStream<CLLocation>.Continuation?

    override init() {
        super.init()
        manager.delegate = self
        manager.desiredAccuracy = kCLLocationAccuracyBest
        manager.distanceFilter = 3
        manager.allowsBackgroundLocationUpdates = true
        manager.pausesLocationUpdatesAutomatically = true
        manager.activityType = .fitness
    }

    func requestPermission() {
        manager.requestWhenInUseAuthorization()
        manager.requestAlwaysAuthorization()
    }

    func locations() -> AsyncStream<CLLocation> {
        AsyncStream { cont in
            self.continuation = cont
            manager.startUpdatingLocation()
        }
    }

    func stop() {
        manager.stopUpdatingLocation()
        continuation?.finish()
        continuation = nil
    }

    nonisolated func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {
        guard let loc = locations.last else { return }
        if loc.horizontalAccuracy > 40 { return }
        Task { @MainActor in
            self.lastPoint = loc
            self.continuation?.yield(loc)
        }
    }
}
