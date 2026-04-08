#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "$TMP_DIR"
}

trap cleanup EXIT

assert_not_on_path() {
  local target_dir="$1"
  case ":$PATH:" in
    *":$target_dir:"*)
      echo "expected $target_dir to be absent from PATH" >&2
      exit 1
      ;;
  esac
}

assert_globalai_missing() {
  if command -v globalai >/dev/null 2>&1; then
    echo "expected globalai to be missing from PATH" >&2
    exit 1
  fi
}

assert_globalai_present() {
  if ! command -v globalai >/dev/null 2>&1; then
    echo "expected globalai to be available on PATH" >&2
    exit 1
  fi
}

assert_help_works() {
  "$1" --help | grep 'globalai helps inspect AI prompt and instruction sources.' >/dev/null
}

ORIGINAL_PATH="$PATH"
GO_BIN_DIR="$(dirname "$(command -v go)")"
BASE_PATH="$GO_BIN_DIR:/usr/bin:/bin:/usr/sbin:/sbin"

# Scenario 1: explicit GOBIN not on PATH
EXPLICIT_GOBIN="$TMP_DIR/explicit-bin"
mkdir -p "$EXPLICIT_GOBIN"
PATH="$BASE_PATH"
assert_not_on_path "$EXPLICIT_GOBIN"
assert_globalai_missing

(
  cd "$ROOT_DIR"
  GOBIN="$EXPLICIT_GOBIN" PATH="$PATH" go install ./cmd/globalai
)

assert_help_works "$EXPLICIT_GOBIN/globalai"
PATH="$BASE_PATH"
assert_globalai_missing

PATH="$EXPLICIT_GOBIN:$BASE_PATH"
assert_globalai_present
assert_help_works globalai

# Scenario 2: fallback to GOPATH/bin when GOBIN is unset
TEMP_GOPATH="$TMP_DIR/gopath"
mkdir -p "$TEMP_GOPATH"
FALLBACK_BIN="$TEMP_GOPATH/bin"
PATH="$BASE_PATH"
assert_not_on_path "$FALLBACK_BIN"
assert_globalai_missing

(
  cd "$ROOT_DIR"
  GOPATH="$TEMP_GOPATH" PATH="$PATH" go install ./cmd/globalai
)

assert_help_works "$FALLBACK_BIN/globalai"
PATH="$BASE_PATH"
assert_globalai_missing

PATH="$FALLBACK_BIN:$BASE_PATH"
assert_globalai_present
assert_help_works globalai

# Scenario 3: go-run installer chooses a writable PATH directory under HOME
INSTALL_HOME="$TMP_DIR/installer-home"
INSTALL_PATH_DIR="$INSTALL_HOME/bin"
mkdir -p "$INSTALL_PATH_DIR"
PATH="$INSTALL_PATH_DIR:$BASE_PATH"
assert_globalai_missing

(
  cd "$ROOT_DIR"
  HOME="$INSTALL_HOME" SHELL="/bin/zsh" PATH="$PATH" go run ./cmd/globalai install > "$TMP_DIR/install-output.log"
)

assert_globalai_present
assert_help_works globalai
grep 'globalai is ready on PATH' "$TMP_DIR/install-output.log" >/dev/null
