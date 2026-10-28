#!/usr/bin/env bash
# Builds the files published to https://assets.invopop.com/popui/<version>/
# into the given directory: PopUI's browser assets, the jsonui runtime and the
# jsonui WebAssembly renderer with the wasm_exec.js from the Go that built it.
set -euo pipefail

out=${1:?usage: scripts/cdn-dist.sh <dir>}
cd "$(dirname "$0")/.."
mkdir -p "$out"

cp assets/* jsonui/assets/jsonui.js "$out/"
# The CSS loads its fonts from the default /_popui/assets/ path that Go apps
# serve them on; on the CDN they sit next to the CSS.
sed -i 's|/_popui/assets/||g' "$out/popui.css"

GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o "$out/popui.wasm" ./jsonui/wasm
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$out/"
