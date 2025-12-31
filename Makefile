# spellchecker:disable

# Go itself
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOMOD=$(GOCMD) mod
GOGENERATE=$(GOCMD) generate
CSPELL=cspell
SHELLCHECK=shellcheck

# Flags for building and versioning
GGMANVERSIONFLAGS=-X 'go.tkw01536.de/ggman.buildVersion=$(shell git describe --tags HEAD)'
BUILDFLAGS=-mod=readonly -trimpath -buildvcs=false -ldflags="$(GGMANVERSIONFLAGS) -s -w -buildid="

# Binary paths
DIST_DIR=dist
BINARY_NAME=ggman
BINARY_UNIX=$(DIST_DIR)/$(BINARY_NAME)
BINARY_MACOS_INTEL=$(DIST_DIR)/$(BINARY_NAME)_mac_intel
BINARY_MACOS=$(DIST_DIR)/$(BINARY_NAME)_mac
BINARY_WINDOWS=$(DIST_DIR)/$(BINARY_NAME).exe

# the path to the ggman command sources
GGMAN_CMD_SRC=./cmd/ggman

# almost all the targets are phony

all: dist
.PHONY: all install test lint testdeps clean deps dist generate shellcheck spellcheck golint

$(BINARY_NAME): deps
	CGO_ENABLED=0 $(GOBUILD) $(BUILDFLAGS) -o $(BINARY_NAME) $(GGMAN_CMD_SRC)

$(BINARY_UNIX): deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILDFLAGS) -o $(BINARY_UNIX) $(GGMAN_CMD_SRC)

$(BINARY_MACOS_INTEL): deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILDFLAGS) -o $(BINARY_MACOS_INTEL) $(GGMAN_CMD_SRC)

$(BINARY_MACOS): deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILDFLAGS) -o $(BINARY_MACOS) $(GGMAN_CMD_SRC)

$(BINARY_WINDOWS): deps
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILDFLAGS)  -o $(BINARY_WINDOWS) $(GGMAN_CMD_SRC)

dist: $(BINARY_UNIX) $(BINARY_MACOS) $(BINARY_MACOS_INTEL) $(BINARY_WINDOWS)

install:
	CGO_ENABLED=0 $(GOINSTALL) $(BUILDFLAGS) $(GGMAN_CMD_SRC)

test: testdeps
	$(GOTEST) -count=1 ./...
testdeps:
	$(GOMOD) download

lint: golint shellcheck spellcheck 

golint:
	test -z $(shell $(GOTOOL) gofmt -l .)
	$(GOTOOL) golangci-lint run ./...
	$(GOTOOL) govulncheck

shellcheck:
	$(SHELLCHECK) --shell=sh internal/cmd/shellrc.sh
	$(SHELLCHECK) --shell=bash internal/cmd/shellrc.sh
	$(SHELLCHECK) --shell=dash internal/cmd/shellrc.sh
	$(SHELLCHECK) --shell=ksh internal/cmd/shellrc.sh

spellcheck:
	$(CSPELL) lint .

generate:
	$(GOGENERATE) -v ./...

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(DIST_DIR)

deps:
	$(GOMOD) download
