#! /usr/bin/env bash
set -euo pipefail

# Always work from the root of the repo.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR"/.. && pwd)"
cd "$ROOT_DIR"

cd "mocks/git/github.com"
for d in *
do
    ( cd "$d" && git init && git add . && git commit -m 'Initial Commit' )
done
