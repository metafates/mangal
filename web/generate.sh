#!/usr/bin/env bash

set -exo pipefail

DIR="$(dirname -- "${BASH_SOURCE[0]}")"
DIR="$(realpath "$DIR")"

OPENAPI_SCHEMA="$DIR/mangal.yaml"

pushd "$DIR/api"
oapi-codegen --config server.cfg.yaml "$OPENAPI_SCHEMA"
oapi-codegen --config types.cfg.yaml "$OPENAPI_SCHEMA"
popd

npx openapi-typescript "$OPENAPI_SCHEMA" -o "$DIR"/ui/src/lib/mangal.d.ts
