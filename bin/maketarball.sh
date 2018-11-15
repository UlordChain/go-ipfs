#!/usr/bin/env bash
# vim: set expandtab sw=2 ts=2:

# bash safe mode
set -euo pipefail
IFS=$'\n\t'


OUTPUT=$(realpath ${1:-go-udfs-source.tar.gz})

TMPDIR="$(mktemp -d)"
NEWUDFS="$TMPDIR/github.com/udfs/go-udfs"
mkdir -p "$NEWUDFS"
cp -r . "$NEWUDFS"
( cd "$NEWUDFS" &&
  echo $PWD &&
  GOPATH="$TMPDIR" gx install --local &&
  (git rev-parse --short HEAD || true) > .tarball &&
  tar -czf "$OUTPUT" --exclude="./.git" .
)

rm -rf "$TMPDIR"
