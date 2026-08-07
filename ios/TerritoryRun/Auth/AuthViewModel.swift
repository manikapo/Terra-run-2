import Foundation
import Supabase

enum AppConfig {
    static let apiBaseURL = Bundle.main.object(forInfoDictionaryKey: "API_BASE_URL") as? String
        ?? "https://territory-run-api.onrender.com"
    static let supabaseURL = Bundle.main.object(forInfoDictionaryKey: "SUPABASE_URL") as? String
        ?? "https://YOUR_REF.supabase.co"
    static let supabaseAnonKey = Bundle.main.object(forInfoDictionaryKey: "SUPABASE_ANON_KEY") as? String
        ?? "YOUR_ANON_KEY"
    static let otaManifestURL = Bundle.main.object(forInfoDictionaryKey: "OTA_MANIFEST_URL") as? String
        ?? "https://cdn.yourdomain.com/ota/manifest.json"
}

let supabase = SupabaseClient(
    supabaseURL: URL(string: AppConfig.supabaseURL)!,
    supabaseKey: AppConfig.supabaseAnonKey
)

@MainActor
final class AuthViewModel: ObservableObject {
    @Published var isLoggedIn = false
    @Published var accessToken: String?

    func refreshSession() async {
        if let session = try? await supabase.auth.session {
            accessToken = session.accessToken
            isLoggedIn = true
            ApiClient.token = session.accessToken
        }
    }

    func signInWithGoogle(idToken: String) async throws {
        try await supabase.auth.signInWithIdToken(
            credentials: .init(provider: .google, idToken: idToken)
        )
        await refreshSession()
    }

    func signOut() async {
        try? await supabase.auth.signOut()
        isLoggedIn = false
        accessToken = nil
        ApiClient.token = nil
    }
}
