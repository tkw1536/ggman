# Go itself
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary paths
OUT_DIR=out
BINARY_NAME=ggman
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_MACOS=$(BINARY_NAME)_mac
BINARY_WINDOWS=$(BINARY_NAME)_windows.exe

all: test build

build: build-local build-linux build-macos build-windows
build-local: deps
	$(GOBUILD) -o $(OUT_DIR)/$(BINARY_NAME)
build-linux: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o $(OUT_DIR)/$(BINARY_UNIX)
build-macos: deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o $(OUT_DIR)/$(BINARY_MACOS)
build-windows: deps
	go get golang.org/x/sys/windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o $(OUT_DIR)/$(BINARY_WINDOWS)
build-upx: build-linux build-macos build-windows
	upx --brute $(OUT_DIR)/$(BINARY_UNIX)
	upx --brute $(OUT_DIR)/$(BINARY_MACOS)
	upx --brute $(OUT_DIR)/$(BINARY_WINDOWS)


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