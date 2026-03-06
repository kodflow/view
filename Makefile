# view - Screen capture over LAN

MODULE := github.com/kodflow/view
SRC := ./src
BIN := ./bin

SERVER_BIN := $(BIN)/view-server
CLIENT_BIN := $(BIN)/view-client

GOFLAGS := -trimpath -ldflags="-s -w"

.PHONY: all build build-server build-client build-all test lint fmt clean deps help run/server run/client

all: help

help:
	@echo ""
	@echo "  view - Screen capture over LAN"
	@echo ""
	@echo "  Usage: make <target>"
	@echo ""
	@echo "  Build:"
	@echo "    build          Build server + client"
	@echo "    build-server   Build server only"
	@echo "    build-client   Build client only"
	@echo "    build-all      Cross-compile (macOS arm64 + Windows amd64)"
	@echo ""
	@echo "  Run:"
	@echo "    run/server     Build and start the server (port 22)"
	@echo "    run/client     Build and start the client (SERVER=ip)"
	@echo ""
	@echo "  Dev:"
	@echo "    deps           Download Go dependencies"
	@echo "    test           Run unit tests"
	@echo "    lint           Run golangci-lint"
	@echo "    fmt            Format Go source files"
	@echo ""
	@echo "  Clean:"
	@echo "    clean          Remove build artifacts"
	@echo ""

deps:
	cd $(SRC) && go mod download

build: build-server build-client

build-server:
	cd $(SRC) && go build $(GOFLAGS) -o ../$(SERVER_BIN) ./cmd/server

build-client:
	cd $(SRC) && go build $(GOFLAGS) -o ../$(CLIENT_BIN) ./cmd/client

build-all:
	cd $(SRC) && GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o ../$(SERVER_BIN)-darwin-arm64 ./cmd/server
	cd $(SRC) && GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o ../$(CLIENT_BIN)-darwin-arm64 ./cmd/client
	cd $(SRC) && GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o ../$(CLIENT_BIN)-windows-amd64.exe ./cmd/client

test:
	cd $(SRC) && go test -race -count=1 ./...

lint:
	cd $(SRC) && golangci-lint run ./...

fmt:
	cd $(SRC) && gofmt -s -w .

run/server: build-server
	$(SERVER_BIN)

run/client: build-client
	$(CLIENT_BIN)

clean:
	rm -rf $(BIN)
