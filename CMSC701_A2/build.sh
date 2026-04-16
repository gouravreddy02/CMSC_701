#!/usr/bin/env bash
set -e

/usr/local/go/bin/go build -o rsbuild ./cmd/rsbuild
/usr/local/go/bin/go build -o rsquery_rank ./cmd/rsquery_rank
/usr/local/go/bin/go build -o rsquery_select ./cmd/rsquery_select
/usr/local/go/bin/go build -o sarray ./cmd/sarray