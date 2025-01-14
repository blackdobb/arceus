#!/bin/bash

set -o errexit
set -u

readonly commitHash=$(git log -n1 --format=format:'%H')
go build \
  -ldflags="-s -w -X github.com/zc2638/arceus/pkg/version.version=$commitHash" \
  -installsuffix cgo \
  -o arceus \
  github.com/zc2638/arceus/cmd