#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR"

if [ ! -d "$SCRIPT_DIR/front" ]; then
  echo "front directory not found" >&2
  exit 1
fi

export COREPACK_HOME="${COREPACK_HOME:-$SCRIPT_DIR/.cache/corepack}"
mkdir -p "$COREPACK_HOME"

if [ -f "$SCRIPT_DIR/front/package-lock.json" ]; then
  npm --prefix "$SCRIPT_DIR/front" run build
elif [ -f "$SCRIPT_DIR/front/.yarn/releases/yarn-4.13.0.cjs" ]; then
  node "$SCRIPT_DIR/front/.yarn/releases/yarn-4.13.0.cjs" --cwd "$SCRIPT_DIR/front" build
else
  echo "No supported frontend package manager configuration found." >&2
  exit 1
fi

mkdir -p "$SCRIPT_DIR/cmd/mmemo/static"
cp -R "$SCRIPT_DIR/front/dist/." "$SCRIPT_DIR/cmd/mmemo/static/"

if git diff --quiet -- cmd/mmemo/static; then
  echo "No frontend asset changes."
  exit 0
fi

git add cmd/mmemo/static
