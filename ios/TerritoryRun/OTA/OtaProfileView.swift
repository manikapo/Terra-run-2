import SwiftUI
import WebKit

struct OtaProfileView: View {
    var body: some View {
        OtaWebView()
            .ignoresSafeArea()
    }
}

struct OtaWebView: UIViewRepresentable {
    func makeUIView(context: Context) -> WKWebView {
        let config = WKWebViewConfiguration()
        config.userContentController.add(context.coordinator, name: "TerritoryBridge")
        let webView = WKWebView(frame: .zero, configuration: config)
        webView.loadFileURL(
            OtaManager.localProfileURL(),
            allowingReadAccessTo: OtaManager.otaDirectory()
        )
        return webView
    }

    func updateUIView(_ uiView: WKWebView, context: Context) {}

    func makeCoordinator() -> Coordinator { Coordinator() }

    class Coordinator: NSObject, WKScriptMessageHandler {
        func userContentController(_ userContentController: WKUserContentController, didReceive message: WKScriptMessage) {}
    }
}

enum OtaManager {
    static func otaDirectory() -> URL {
        FileManager.default.urls(for: .documentDirectory, in: .userDomainMask)[0]
            .appendingPathComponent("ota/current", isDirectory: true)
    }

    static func localProfileURL() -> URL {
        otaDirectory().appendingPathComponent("pages/profile.html")
    }

    static func syncManifest() async {
        // Mirror Android OtaManager — fetch manifest.json, verify sha256, unzip bundle
    }
}
