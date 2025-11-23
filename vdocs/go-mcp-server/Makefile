.PHONY: all build test clean install run help

# Variables
BINARY_NAME=go-mcp-server
BUILD_DIR=bin
CMD_DIR=cmd/server
VERSION?=1.0.0
GO=go

# Build flags
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

all: test build

# Build the binary
build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	${GO} build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ./${CMD_DIR}
	@echo "Build complete: ${BUILD_DIR}/${BINARY_NAME}"

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	@mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 ${GO} build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 ./${CMD_DIR}
	GOOS=linux GOARCH=arm64 ${GO} build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-arm64 ./${CMD_DIR}

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p ${BUILD_DIR}
	GOOS=darwin GOARCH=amd64 ${GO} build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 ./${CMD_DIR}
	GOOS=darwin GOARCH=arm64 ${GO} build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 ./${CMD_DIR}

build-windows:
	@echo "Building for Windows..."
	@mkdir -p ${BUILD_DIR}
	GOOS=windows GOARCH=amd64 ${GO} build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe ./${CMD_DIR}

# Run tests
test:
	@echo "Running tests..."
	${GO} test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run tests with coverage
test-coverage: test
	@echo "Generating coverage report..."
	${GO} tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linters
lint:
	@echo "Running linters..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed, installing..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	${GO} fmt ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	${GO} mod tidy

# Install binary
install: build
	@echo "Installing ${BINARY_NAME}..."
	cp ${BUILD_DIR}/${BINARY_NAME} ${GOPATH}/bin/
	@echo "Installed to ${GOPATH}/bin/${BINARY_NAME}"

# Run the server
run: build
	@echo "Running ${BINARY_NAME}..."
	./${BUILD_DIR}/${BINARY_NAME}

# Run with environment variables
run-dev:
	@echo "Running ${BINARY_NAME} in development mode..."
	ALIBABA_CLOUD_ACCESS_KEY_ID=${ALIBABA_CLOUD_ACCESS_KEY_ID} \
	ALIBABA_CLOUD_ACCESS_KEY_SECRET=${ALIBABA_CLOUD_ACCESS_KEY_SECRET} \
	./${BUILD_DIR}/${BINARY_NAME} --transport stdio

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf ${BUILD_DIR}
	rm -f coverage.txt coverage.html
	@echo "Clean complete"

# Show help
help:
	@echo "Available targets:"
	@echo "  all           - Run tests and build"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for all platforms"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-darwin  - Build for macOS"
	@echo "  build-windows - Build for Windows"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Run linters"
	@echo "  fmt           - Format code"
	@echo "  tidy          - Tidy dependencies"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  run           - Build and run"
	@echo "  run-dev       - Run with development settings"
	@echo "  clean         - Clean build artifacts"
	@echo "  help          - Show this help message"

