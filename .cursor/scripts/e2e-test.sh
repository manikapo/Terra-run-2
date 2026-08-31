#!/usr/bin/env bash
set -euo pipefail

echo "=== Territory Run E2E Validation ==="
echo ""

echo "1. Health check"
curl -sf http://localhost:8080/health | python3 -m json.tool
echo ""

echo "2. OTA web app"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/index.html)
echo "index.html HTTP status: $HTTP_CODE"
test "$HTTP_CODE" = "200"
echo ""

echo "3. Guest user flow"
GUEST_ID=$(python3 -c "import uuid; print(uuid.uuid4())")
echo "Guest ID: $GUEST_ID"

PROFILE=$(curl -sf -H "X-Guest-User: $GUEST_ID" http://localhost:8080/api/v1/users/me)
echo "Profile: $PROFILE"

ACTIVITY=$(curl -sf -X POST -H "Content-Type: application/json" -H "X-Guest-User: $GUEST_ID" \
  -d "{\"idempotency_key\":\"e2e-$(date +%s)\",\"device_info\":{\"platform\":\"e2e-test\"}}" \
  http://localhost:8080/api/v1/activities)
echo "Activity created: $ACTIVITY"
ACTIVITY_ID=$(echo "$ACTIVITY" | python3 -c "import sys,json; print(json.load(sys.stdin)['activity_id'])")

curl -sf -X POST -H "Content-Type: application/json" -H "X-Guest-User: $GUEST_ID" \
  -d '{"points":[{"lat":37.7749,"lon":-122.4194,"ts":1756617600,"acc":10},{"lat":37.7750,"lon":-122.4195,"ts":1756617660,"acc":10}]}' \
  "http://localhost:8080/api/v1/activities/$ACTIVITY_ID/points" | python3 -m json.tool

COMPLETE=$(curl -sf -X POST -H "Content-Type: application/json" -H "X-Guest-User: $GUEST_ID" \
  -d '{"ended_at":"2026-08-31T06:05:00Z","distance_m":150,"duration_s":120}' \
  "http://localhost:8080/api/v1/activities/$ACTIVITY_ID/complete")
echo "Complete: $COMPLETE"

UPDATED=$(curl -sf -H "X-Guest-User: $GUEST_ID" http://localhost:8080/api/v1/users/me)
echo "Updated profile: $UPDATED"
echo ""

echo "=== All E2E checks passed ==="
