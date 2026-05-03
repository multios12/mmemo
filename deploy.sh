#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

LOCAL_BIN="${LOCAL_BIN:-$SCRIPT_DIR/dist/mmemo}"
REMOTE_HOST="${REMOTE_HOST:-modern-utopia.net}"
REMOTE_PATH="${REMOTE_PATH:-/home/docker/mmemo/mmemo}"
REMOTE_DIR=$(dirname -- "$REMOTE_PATH")
REMOTE_TMP="$REMOTE_PATH.tmp"
SSH_OPTS="${SSH_OPTS:-}"

if [ ! -f "$LOCAL_BIN" ]; then
  echo "Local binary not found: $LOCAL_BIN" >&2
  echo "Run ./build.sh first." >&2
  exit 1
fi

scp $SSH_OPTS "$LOCAL_BIN" "$REMOTE_HOST:$REMOTE_TMP"
ssh $SSH_OPTS "$REMOTE_HOST" "mkdir -p '$REMOTE_DIR' && chmod 755 '$REMOTE_TMP' && mv '$REMOTE_TMP' '$REMOTE_PATH'"

echo "Deployed $LOCAL_BIN to $REMOTE_HOST:$REMOTE_PATH"
