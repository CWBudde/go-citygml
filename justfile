set shell := ["bash", "-uc"]

# Default recipe - show available commands
default:
    @just --list

# Format all code using treefmt
fmt:
    treefmt --allow-missing-formatter

# Check if code is formatted correctly
check-formatted:
    treefmt --allow-missing-formatter --fail-on-change

# Run linters
lint:
    GOCACHE="${GOCACHE:-/tmp/gocache}" GOMODCACHE="${GOMODCACHE:-/tmp/gomodcache}" GOLANGCI_LINT_CACHE="${GOLANGCI_LINT_CACHE:-/tmp/golangci-lint-cache}" golangci-lint run --timeout=2m ./...

# Run linters with auto-fix
lint-fix:
    GOCACHE="${GOCACHE:-/tmp/gocache}" GOMODCACHE="${GOMODCACHE:-/tmp/gomodcache}" GOLANGCI_LINT_CACHE="${GOLANGCI_LINT_CACHE:-/tmp/golangci-lint-cache}" golangci-lint run --fix --timeout=2m ./...

# Ensure go.mod is tidy
check-tidy:
    go mod tidy
    git diff --exit-code go.mod go.sum

# Run all tests
test:
    go test -v ./...

# Run tests with race detector
test-race:
    go test -race ./...

# Run tests with coverage
test-coverage:
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

# Run all checks (formatting, linting, tests, tidiness)
ci: check-formatted test lint check-tidy

# Build the CLI tool
build:
    go build -o citygml ./cmd/citygml

# Build the WASM module for the web demo
build-wasm:
    GOOS=js GOARCH=wasm go build -o web/citygml.wasm ./cmd/citygmlwasm
    @GOROOT=$(go env GOROOT) && \
    if [ -f "$$GOROOT/misc/wasm/wasm_exec.js" ]; then \
        cp "$$GOROOT/misc/wasm/wasm_exec.js" web/; \
    elif [ -f "$$GOROOT/lib/wasm/wasm_exec.js" ]; then \
        cp "$$GOROOT/lib/wasm/wasm_exec.js" web/; \
    fi
    @echo "WASM build complete: web/citygml.wasm"

# Clean build artifacts
clean:
    rm -f coverage.out coverage.html citygml web/citygml.wasm web/wasm_exec.js
