## ── HELPERS ──

.PHONY: help
help: ## show this help
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## ── DEVELOPMENT ──

.PHONY: run/latticed
run/latticed: ## run the daemon
	go run ./cmd/latticed

.PHONY: run/lattice
run/lattice: ## run the CLI
	go run ./cmd/lattice

.PHONY: proto/generate
proto/generate: ## regenerate protobuf code
	protoc --go_out=. --go-grpc_out=. udf/proto/udf.proto
	protoc --go_out=. --go-grpc_out=. server/proto/lattice.proto

## ── QUALITY ──

.PHONY: audit
audit: ## run vet, staticcheck, and tests
	go vet ./...
	staticcheck ./...
	go test -race -vet=off ./...

.PHONY: test
test: ## run unit tests
	go test ./...

.PHONY: test/verbose
test/verbose: ## run unit tests with verbose output
	go test -v ./...

.PHONY: test/cover
test/cover: ## run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: test/race
test/race: ## run tests with race detector
	go test -race ./...

.PHONY: test/integration
test/integration: ## run integration tests (requires running services)
	go test -tags=integration ./store/... ./blobstore/... ./vectorstore/... ./dispatch/... ./auth/...

.PHONY: lint
lint: ## run golangci-lint
	golangci-lint run

.PHONY: fmt
fmt: ## format code
	gofmt -s -w .

.PHONY: fmt/check
fmt/check: ## check formatting
	@test -z "$$(gofmt -l .)" || (echo "Files need formatting:" && gofmt -l . && exit 1)

## ── BUILD ──

.PHONY: build/latticed
build/latticed: ## build the daemon
	go build -o bin/latticed ./cmd/latticed

.PHONY: build/lattice
build/lattice: ## build the CLI
	go build -o bin/lattice ./cmd/lattice

.PHONY: build/all
build/all: build/latticed build/lattice ## build everything

## ── CLEAN ──

.PHONY: clean
clean: ## remove build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
