# Go itself
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOMOD=$(GOCMD) mod
GOGENERATE=$(GOCMD) generate

# Flags for versioning
GGMANVERSIONFLAGS=-X 'go.tkw01536.de/ggman.buildVersion=$(shell git describe --tags HEAD)' -X 'go.tkw01536.de/ggman.buildTime=$(shell date +%s)'

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
.PHONY: all install test lint testdeps clean deps dist generate

$(BINARY_NAME): deps
	CGO_ENABLED=0 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS)" -o $(BINARY_NAME) $(GGMAN_CMD_SRC)

$(BINARY_UNIX): deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(BINARY_UNIX) $(GGMAN_CMD_SRC)

$(BINARY_MACOS_INTEL): deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(BINARY_MACOS_INTEL) $(GGMAN_CMD_SRC)

$(BINARY_MACOS): deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(BINARY_MACOS) $(GGMAN_CMD_SRC)

$(BINARY_WINDOWS): deps
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(BINARY_WINDOWS) $(GGMAN_CMD_SRC)

dist: $(BINARY_UNIX) $(BINARY_MACOS) $(BINARY_MACOS_INTEL) $(BINARY_WINDOWS)

install:
	CGO_ENABLED=0 $(GOINSTALL) -ldflags="$(GGMANVERSIONFLAGS)" $(GGMAN_CMD_SRC)

test: testdeps
	$(GOTEST) -count=1 ./...
testdeps:
	$(GOMOD) download

lint:
	test -z $(shell gofmt -l .)
	$(GOTOOL) golangci-lint run ./...
	$(GOTOOL) govulncheck


generate:
	$(GOGENERATE) -v ./...

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(DIST_DIR)

deps:
	$(GOMOD) download
