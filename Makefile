# Go itself
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGENERATE=$(GOCMD) generate

# Flags for versioning
GGMANVERSIONFLAGS=-X 'github.com/tkw1536/ggman/constants.buildVersion=$(shell git describe --tags HEAD)' -X 'github.com/tkw1536/ggman/constants.buildTime=$(shell date +%s)'

# Binary paths
DIST_DIR=dist
BINARY_NAME=ggman
BINARY_UNIX=$(DIST_DIR)/$(BINARY_NAME)
BINARY_MACOS_INTEL=$(DIST_DIR)/$(BINARY_NAME)_mac_intel
BINARY_MACOS_SILICON=$(DIST_DIR)/$(BINARY_NAME)_mac_silicon
BINARY_MACOS_UNIVERSAL=$(DIST_DIR)/$(BINARY_NAME)_mac
BINARY_WINDOWS=$(DIST_DIR)/$(BINARY_NAME).exe


# the path to the ggman command sources
GGMAN_CMD_SRC=./cmd/ggman

# almost all the targets are phony

all: dist
.PHONY: all install test lint testdeps clean deps dist generate

$(BINARY_NAME): deps
	$(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS)" -o $(BINARY_NAME) $(GGMAN_CMD_SRC)

$(BINARY_UNIX): deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(BINARY_UNIX) $(GGMAN_CMD_SRC)

$(BINARY_MACOS_INTEL): deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(BINARY_MACOS_INTEL) $(GGMAN_CMD_SRC)

$(BINARY_MACOS_SILICON): deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(BINARY_MACOS_SILICON) $(GGMAN_CMD_SRC)

$(BINARY_MACOS_UNIVERSAL): $(BINARY_MACOS_INTEL) $(BINARY_MACOS_SILICON)
	lipo -output $(BINARY_MACOS_UNIVERSAL) -create $(BINARY_MACOS_SILICON) $(BINARY_MACOS_INTEL)

$(BINARY_WINDOWS): deps
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(BINARY_WINDOWS) $(GGMAN_CMD_SRC)

dist: $(BINARY_UNIX) $(BINARY_MACOS_UNIVERSAL) $(BINARY_WINDOWS)

install:
	$(GOINSTALL) -ldflags="$(GGMANVERSIONFLAGS)" $(GGMAN_CMD_SRC)

test: testdeps
	$(GOTEST) -tags doccheck ./...
testdeps:
	$(GOGET) -v ./...

lint:
	test -z $(shell gofmt -l .)

generate:
	$(GOGENERATE) -v ./...

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(DIST_DIR)

deps:
	$(GOGET) -v ./...
