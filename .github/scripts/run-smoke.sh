#!/usr/bin/env bash
set -euo pipefail

# Assumptions:
#  - Frontend has already been built (vite build) if needed
#  - 'go run server.go' launches the application
#  - Success condition: HTTPS port responds with any content OR log line 'Starting API services...'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
LOG_DIR="$ROOT_DIR/smoke-logs"
mkdir -p "$LOG_DIR"
LOG_FILE="$LOG_DIR/server.log"
RESULT_FILE="$LOG_DIR/result.txt"
STATUS_JSON="$LOG_DIR/status.json"
PORT="8443"
TIMEOUT_SECONDS=${SMOKE_TIMEOUT_SECONDS:-45}
GRACE_KILL_SECONDS=5

# Start server in background
(
  cd "$ROOT_DIR"
  echo "[smoke] launching server..." | tee -a "$LOG_FILE"
  export LOG_LEVEL=20  # INFO
  # Run server
  nohup go run server.go > "$LOG_FILE" 2>&1 &
  echo $! > "$LOG_DIR/server.pid"
) || true

PID=$(cat "$LOG_DIR/server.pid" 2>/dev/null || echo "")
if [[ -z "$PID" ]] || ! kill -0 "$PID" 2>/dev/null; then
  echo "{\"status\":\"failed\",\"reason\":\"No valid server PID captured\"}" > "$STATUS_JSON"
  exit 0
fi

echo "[smoke] server pid: $PID"

# Wait for readiness
START_TIME=$(date +%s)
READY=0
while true; do
  NOW=$(date +%s)
  ELAPSED=$((NOW-START_TIME))
  if (( ELAPSED > TIMEOUT_SECONDS )); then
    echo "[smoke] timeout after $ELAPSED seconds" | tee -a "$LOG_FILE"
    break
  fi
  # Check log line
  if grep -q "Starting API services" "$LOG_FILE"; then
    READY=1
    echo "[smoke] readiness via log line" | tee -a "$LOG_FILE"
    break
  fi
  # Try HTTPS curl (ignore cert issues)
  if curl -sk --max-time 2 https://localhost:$PORT/ >/dev/null 2>&1; then
    READY=1
    echo "[smoke] readiness via https response" | tee -a "$LOG_FILE"
    break
  fi
  sleep 2
done

if [[ $READY -eq 1 ]]; then
  STATUS="passed"
else
  STATUS="failed"
fi

echo "[smoke] stopping server (status=$STATUS)" | tee -a "$LOG_FILE"
if kill "$PID" 2>/dev/null; then
  # allow graceful shutdown
  for i in $(seq 1 $GRACE_KILL_SECONDS); do
    if kill -0 "$PID" 2>/dev/null; then
      sleep 1
    else
      break
    fi
  done
  if kill -0 "$PID" 2>/dev/null; then
    echo "[smoke] force killing $PID" | tee -a "$LOG_FILE"
    kill -9 "$PID" || true
  fi
fi

echo "$STATUS" > "$RESULT_FILE"
cat > "$STATUS_JSON" <<EOF
{ "status": "$STATUS", "timeoutSeconds": $TIMEOUT_SECONDS }
EOF

echo "[smoke] done" | tee -a "$LOG_FILE"

# Always exit 0 so workflow can collect artifacts; status is in files.
exit 0
