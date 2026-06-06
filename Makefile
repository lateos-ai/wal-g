.PHONY: build test lint security-check docker-build docker-push clean help compatibility-test

# Variables
BINARY_NAME := wal-g
DOCKER_REGISTRY := lateos
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(BINARY_NAME)
VERSION := $(shell git describe --tags --always)

help:
@echo "WAL-G Fork Makefile"
@echo ""
@echo "Targets:"
@echo "  build              - Build binary"
@echo "  test               - Run unit tests"
@echo "  compatibility-test - Test backward compatibility"
@echo "  security-check     - Run security checks"
@echo "  clean              - Remove build artifacts"

build:
@echo "Building $(BINARY_NAME)..."
go build -ldflags="-X main.Version=$(VERSION)" -o $(BINARY_NAME) .
@echo "Build complete"

test:
@echo "Running tests..."
go test -v -race -coverprofile=coverage.out ./...

compatibility-test:
@echo "Running compatibility tests..."
@echo "Placeholder for compatibility testing"

clean:
@echo "Cleaning..."
Remove-Item -Force $(BINARY_NAME) -ErrorAction SilentlyContinue
Remove-Item -Force coverage.* -ErrorAction SilentlyContinue
