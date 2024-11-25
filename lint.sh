#!/bin/bash
set -e

echo "=> go vet"
go vet ./...

echo "=> staticcheck"
go run honnef.co/go/tools/cmd/staticcheck@latest ./...

echo "=> golangci-lint"
go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...

echo "=> govulncheck"
go run golang.org/x/vuln/cmd/govulncheck@latest

echo "=> gosec"
go run github.com/securego/gosec/v2/cmd/gosec@latest ./...