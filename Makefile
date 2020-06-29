# Go itself
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Flags for versioning
GGMANVERSIONFLAGS=-X 'github.com/tkw1536/ggman/src/constants.BuildVersion=$(shell git describe --tags HEAD)' -X 'github.com/tkw1536/ggman/src/constants.buildTime=$(shell date +%s)'

# Binary paths
DIST_DIR=dist
BINARY_NAME=ggman
BINARY_UNIX=$(DIST_DIR)/$(BINARY_NAME)
BINARY_MACOS=$(DIST_DIR)/$(BINARY_NAME)_mac
BINARY_WINDOWS=$(DIST_DIR)/$(BINARY_NAME).exe

# almost all the targets are phone


all: $(BINARY_UNIX) $(BINARY_MACOS) $(BINARY_WINDOWS)
.PHONY: all install test testdeps clean deps dist

$(BINARY_NAME): deps
	$(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS)" -o $(BINARY_NAME)

$(BINARY_UNIX): deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(BINARY_UNIX)

$(BINARY_MACOS): deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(BINARY_MACOS)

$(BINARY_WINDOWS): deps
	-go get golang.org/x/sys/windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(BINARY_WINDOWS)

dist: $(BINARY_UNIX) $(BINARY_MACOS) $(BINARY_WINDOWS)
	upx --brute $(BINARY_UNIX) $(BINARY_MACOS) $(BINARY_WINDOWS)

install:
	$(GOINSTALL) -ldflags="$(GGMANVERSIONFLAGS)"

test: testdeps
	$(GOTEST) -v ./...
testdeps:
	$(GOGET) -v ./...

clean: 
	$(GOCLEAN)
	rm $(BINARY_NAME)
	rm -rf $(OUT_DIR)

deps:
	$(GOGET) -v ./...
