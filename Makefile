
APP_NAME := homelab-inventory
BUILD_DIR := bin
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GO_VERSION := $(shell go version | cut -d' ' -f3)
LDFLAGS := "-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME) -X main.GoVer=$(GO_VERSION)"

all: build-linux-amd64 build-linux-arm64 build-windows build-mac-arm64

build-linux-amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 main.go

build-linux-arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 main.go

build-windows:
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME).exe main.go

build-mac-arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-mac-arm64 main.go

clean:
	rm -rf $(BUILD_DIR)
