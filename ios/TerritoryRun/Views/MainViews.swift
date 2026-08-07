import SwiftUI

struct LoginView: View {
    @EnvironmentObject var auth: AuthViewModel

    var body: some View {
        VStack(spacing: 16) {
            Text("Territory Run").font(.largeTitle.bold())
            Text("Capture turf on every run")
            Button("Sign in with Google") {
                // TODO: GoogleSignIn → auth.signInWithGoogle(idToken:)
                Task { await auth.refreshSession() }
            }
        }
        .padding()
    }
}

struct HomeView: View {
    @State private var profile: UserProfile?
    @State private var showRun = false

    var body: some View {
        NavigationStack {
            VStack(alignment: .leading, spacing: 12) {
                if let p = profile {
                    Text("Cells: \(p.cells_owned)")
                    Text("Score: \(p.total_capture_score)")
                }
                Button("Start Run") { showRun = true }
                NavigationLink("Profile (OTA)") { OtaProfileView() }
            }
            .padding()
            .navigationTitle("Territory Run")
            .task { profile = try? await ApiClient.me() }
            .sheet(isPresented: $showRun) { RunTrackingView() }
        }
    }
}
