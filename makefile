.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build   - Build the application"
	@echo "  run     - Run the application"
	@echo "  test    - Run tests"
	@echo "  clean   - Clean up build artifacts"

.PHONY: build/api
build/api:
	@echo "Building API server..."
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build \
		-ldflags="-w -s -extldflags '-static' -X main.version=$(shell git describe --tags --always --dirty)" \
		-a -installsuffix cgo \
		-o ./bin/api \
		./cmd/api
	@echo "API server built successfully to ./bin/api"

.PHONY: build/local
build/local:
	@echo "Building API server for local development..."
	go build \
		-ldflags="-X main.version=$(shell git describe --tags --always --dirty)" \
		-o ./bin/api \
		./cmd/api
	@echo "Local API server built successfully to ./bin/api"

.PHONY: run/api
run/api:
	@echo "Running API server..."
	go run ./cmd/api
