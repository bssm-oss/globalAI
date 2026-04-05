#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d)"
SERVER_LOG="$TMP_DIR/server.log"

cleanup() {
  if [[ -n "${SERVER_PID:-}" ]] && kill -0 "$SERVER_PID" 2>/dev/null; then
    kill "$SERVER_PID" 2>/dev/null || true
    wait "$SERVER_PID" 2>/dev/null || true
  fi
  rm -rf "$TMP_DIR"
}

trap cleanup EXIT

mkdir -p "$TMP_DIR/root/.claude" "$TMP_DIR/home/.claude"
printf '# Project agent\n' > "$TMP_DIR/root/AGENTS.md"
printf '# Local Claude\n' > "$TMP_DIR/root/.claude/project.md"
printf '# Global Claude\n' > "$TMP_DIR/home/.claude/global.md"

go build -o "$TMP_DIR/globalai" ./cmd/globalai

(
  cd "$TMP_DIR/root"
  HOME="$TMP_DIR/home" "$TMP_DIR/globalai" web --no-open --addr 127.0.0.1:0 > "$SERVER_LOG" 2>&1
) &
SERVER_PID=$!

for _ in {1..50}; do
  if grep -q 'globalai viewer ready at http://' "$SERVER_LOG"; then
    break
  fi
  sleep 0.2
done

if ! grep -q 'globalai viewer ready at http://' "$SERVER_LOG"; then
  echo 'viewer did not start in time' >&2
  cat "$SERVER_LOG" >&2 || true
  exit 1
fi

URL="$(grep -o 'http://[^ ]*' "$SERVER_LOG" | head -n 1)"
BODY="$(curl --fail --silent "$URL/api/sources")"

echo "$BODY" | grep '"totalSources":3' >/dev/null
echo "$BODY" | grep '"AGENTS.md"' >/dev/null
echo "$BODY" | grep 'project' >/dev/null
echo "$BODY" | grep 'global' >/dev/null
