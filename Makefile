BINARY := cix
VERSION := v0.1.0
BUILD_DIR := ./bin

.PHONY: build test install clean lint

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY) ./cmd/cix

test:
	go test ./... -v -count=1

test-short:
	go test ./... -short

install:
	go install ./cmd/cix

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)


release:
	GOOS=linux   GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-linux-amd64   ./cmd/cix
	GOOS=darwin  GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-darwin-amd64  ./cmd/cix
	GOOS=darwin  GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY)-darwin-arm64  ./cmd/cix
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe ./cmd/cix

# Quick dev test với sample fixture
dev:
	go run ./cmd/cix validate -f testdata/fixtures/sample.gitlab-ci.yml
	go run ./cmd/cix list -f testdata/fixtures/sample.gitlab-ci.yml