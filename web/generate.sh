#!/usr/bin/env bash

set -exo pipefail

DIR="$(dirname -- "${BASH_SOURCE[0]}")"
DIR="$(realpath "$DIR")"
UI_DIR="$DIR/ui"

OPENAPI_SCHEMA="$DIR/openapi.yaml"

main() {
	pushd "$DIR"
	pushd api
	oapi-codegen --config server.cfg.yaml "$OPENAPI_SCHEMA"
	oapi-codegen --config types.cfg.yaml "$OPENAPI_SCHEMA"
	popd

	"$UI_DIR/node_modules/.bin/openapi-typescript" "$OPENAPI_SCHEMA" -o "$UI_DIR/src/api/mangal.ts"

	pushd ui
	npm run build
	popd

	popd
}

main
