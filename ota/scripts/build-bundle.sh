#!/usr/bin/env bash
# Build OTA bundle.zip and update manifest sha256
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BUNDLE_DIR="$ROOT/ota/bundle"
OUT_ZIP="$ROOT/ota/dist/v1/bundle.zip"
MANIFEST="$ROOT/ota/dist/v1/manifest.json"

mkdir -p "$(dirname "$OUT_ZIP")"

cd "$BUNDLE_DIR"
zip -r "$OUT_ZIP" . -x "*.DS_Store"

SHA=$(sha256sum "$OUT_ZIP" | awk '{print $1}')
echo "SHA256: $SHA"

# Update manifest (requires jq)
if command -v jq &>/dev/null; then
  jq --arg sha "$SHA" --arg url "https://cdn.yourdomain.com/ota/v1/bundle.zip" \
    '.sha256 = $sha | .bundle_url = $url' "$ROOT/ota/manifest.json" > "$MANIFEST"
  echo "Updated $MANIFEST"
else
  echo "Install jq to auto-update manifest, or paste SHA into manifest manually"
fi

echo "Upload:"
echo "  $OUT_ZIP → R2 bucket ota/v1/bundle.zip"
echo "  $MANIFEST → R2 bucket ota/v1/manifest.json (or ota/manifest.json)"
