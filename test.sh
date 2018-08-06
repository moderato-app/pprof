#!/usr/bin/env bash

set -e
set -x
MODE=atomic
echo "mode: $MODE" > coverage.txt

# All packages.
PKG=$(go list ./...)

staticcheck -ignore '
github.com/google/pprof/internal/binutils/binutils_test.go:SA6004
github.com/google/pprof/internal/driver/svg.go:SA6004
github.com/google/pprof/internal/report/source_test.go:SA6004
github.com/google/pprof/profile/filter_test.go:SA6004
' $PKG
unused $PKG

# Packages that have any tests.
PKG=$(go list -f '{{if .TestGoFiles}} {{.ImportPath}} {{end}}' ./...)

go test -v $PKG

for d in $PKG; do
  go test -race -coverprofile=profile.out -covermode=$MODE $d
  if [ -f profile.out ]; then
    cat profile.out | grep -v "^mode: " >> coverage.txt
    rm profile.out
  fi
done

