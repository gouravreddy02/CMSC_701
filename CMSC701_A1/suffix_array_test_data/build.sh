#!/usr/bin/env bash

set -e

/usr/local/go/bin/go build -o buildsa ./cmd/buildsa
/usr/local/go/bin/go build -o inspectsa ./cmd/inspectsa
/usr/local/go/bin/go build -o querysa ./cmd/querysa