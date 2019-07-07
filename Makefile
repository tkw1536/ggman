# Go itself
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Flags for versioning
GGMANVERSIONFLAGS=-X 'github.com/tkw1536/ggman/src/constants.BuildVersion=$(shell git describe --tags HEAD)' -X 'github.com/tkw1536/ggman/src/constants.buildTime=$(shell date +%s)'

# Binary paths
OUT_DIR=out
BINARY_NAME=ggman
BINARY_UNIX=$(BINARY_NAME)
BINARY_MACOS=$(BINARY_NAME)_mac
BINARY_WINDOWS=$(BINARY_NAME).exe

all: test build dist

build: build-local
build-local: deps
	$(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS)" -o $(OUT_DIR)/$(BINARY_NAME)

dist: mindist
build-linux: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(OUT_DIR)/dist/$(BINARY_UNIX)
build-macos: deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(OUT_DIR)/dist/$(BINARY_MACOS)
build-windows: deps
	-go get golang.org/x/sys/windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(GGMANVERSIONFLAGS) -s -w" -o $(OUT_DIR)/dist/$(BINARY_WINDOWS)
mindist: build-linux build-macos build-windows
	upx --brute $(OUT_DIR)/dist//$(BINARY_UNIX) $(OUT_DIR)/dist//$(BINARY_MACOS) $(OUT_DIR)/dist//$(BINARY_WINDOWS)


test: testdeps
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf $(OUT_DIR)
run: build-local
	./$(BINARY_NAME)
deps:
	$(GOGET) -v ./...
testdeps:
	$(GOGET) -v -t ./...