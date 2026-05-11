#!/bin/sh

set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR"

export GOCACHE="${GOCACHE:-/tmp/go-build}"

(
  cd api
  go test ./...
)
