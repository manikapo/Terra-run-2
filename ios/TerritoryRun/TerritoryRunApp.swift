// Territory Run — iOS Phase 1
// Open TerritoryRun.xcodeproj in Xcode (create project from these sources)
// Or: File → New → App → SwiftUI, then replace Sources with files below

import SwiftUI

@main
struct TerritoryRunApp: App {
    @StateObject private var auth = AuthViewModel()

    var body: some Scene {
        WindowGroup {
            if auth.isLoggedIn {
                HomeView()
                    .environmentObject(auth)
            } else {
                LoginView()
                    .environmentObject(auth)
            }
        }
    }
}
