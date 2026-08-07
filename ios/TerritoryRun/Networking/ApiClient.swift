import Foundation

struct GpsPoint: Codable {
    let lat: Double
    let lon: Double
    let ts: Int64
    let acc: Double
    let speed: Double
}

struct UserProfile: Codable {
    let user_id: String
    let username: String
    let display_name: String
    let cells_owned: Int
    let territories_captured: Int
    let total_capture_score: Int
    let activity_count: Int
}

struct CompleteActivityResponse: Codable {
    let activity_id: String
    let status: String
    let capture_result: CaptureResult?
    let trust_score: Int
}

struct CaptureResult: Codable {
    let new_cells: Int
    let capture_score: Int
    let h3_cells: [String]
}

enum ApiClient {
    static var token: String?

    static func me() async throws -> UserProfile {
        try await request(path: "/api/v1/users/me", method: "GET")
    }

    static func createActivity(idempotencyKey: String) async throws -> String {
        struct Body: Codable { let started_at: String; let idempotency_key: String }
        struct Resp: Codable { let activity_id: String }
        let resp: Resp = try await request(
            path: "/api/v1/activities",
            method: "POST",
            body: Body(started_at: ISO8601DateFormatter().string(from: Date()), idempotency_key: idempotencyKey)
        )
        return resp.activity_id
    }

    static func uploadPoints(activityId: String, points: [GpsPoint]) async throws {
        struct Body: Codable { let points: [GpsPoint] }
        try await requestVoid(path: "/api/v1/activities/\(activityId)/points", method: "POST", body: Body(points: points))
    }

    static func completeActivity(activityId: String, distanceM: Double, durationS: Int, cells: [String]) async throws -> CompleteActivityResponse {
        struct Body: Codable {
            let ended_at: String
            let distance_m: Double
            let duration_s: Int
            let h3_cells_hint: [String]
        }
        return try await request(
            path: "/api/v1/activities/\(activityId)/complete",
            method: "POST",
            body: Body(
                ended_at: ISO8601DateFormatter().string(from: Date()),
                distance_m: distanceM,
                duration_s: durationS,
                h3_cells_hint: cells
            )
        )
    }

    private static func request<T: Decodable>(path: String, method: String, body: Encodable? = nil) async throws -> T {
        let url = URL(string: AppConfig.apiBaseURL + path)!
        var req = URLRequest(url: url)
        req.httpMethod = method
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let token { req.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization") }
        if let body {
            req.httpBody = try JSONEncoder().encode(AnyEncodable(body))
        }
        let (data, resp) = try await URLSession.shared.data(for: req)
        guard (resp as? HTTPURLResponse)?.statusCode ?? 500 < 300 else {
            throw URLError(.badServerResponse)
        }
        return try JSONDecoder().decode(T.self, from: data)
    }

    private static func requestVoid(path: String, method: String, body: Encodable) async throws {
        let _: [String: String] = try await request(path: path, method: method, body: body)
    }
}

private struct AnyEncodable: Encodable {
    let value: Encodable
    init(_ value: Encodable) { self.value = value }
    func encode(to encoder: Encoder) throws { try value.encode(to: encoder) }
}
