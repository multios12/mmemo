#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR"

export COREPACK_HOME="${COREPACK_HOME:-$SCRIPT_DIR/.cache/corepack}"
mkdir -p "$COREPACK_HOME" "$SCRIPT_DIR/cmd/mmemo/static" "$SCRIPT_DIR/dist"

node "$SCRIPT_DIR/front/.yarn/releases/yarn-4.13.0.cjs" --cwd "$SCRIPT_DIR/front" build
cp -R "$SCRIPT_DIR/front/dist/." "$SCRIPT_DIR/cmd/mmemo/static/"

cd "$SCRIPT_DIR/cmd/mmemo"
export GOOS=linux
export GOARCH=amd64
go build -ldflags="-s -w" -trimpath -o "$SCRIPT_DIR/dist/"
