#!/bin/bash

# Test script for running different types of tests

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
    print_error "Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Load environment variables for testing
if [ -f .env.test ]; then
    export $(cat .env.test | xargs)
fi

case "$1" in
    "unit")
        print_status "Running unit tests..."
        go test -v ./tests/unit/...
        ;;
    "integration")
        print_status "Running integration tests..."
        go test -v ./tests/integration/...
        ;;
    "all")
        print_status "Running all tests..."
        go test -v ./...
        ;;
    "coverage")
        print_status "Running tests with coverage..."
        go test -v -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out -o coverage.html
        print_status "Coverage report generated: coverage.html"
        ;;
    "benchmark")
        print_status "Running benchmark tests..."
        go test -bench=. -benchmem ./...
        ;;
    "clean")
        print_status "Cleaning test artifacts..."
        rm -f coverage.out coverage.html
        print_status "Clean completed"
        ;;
    *)
        echo "Usage: $0 {unit|integration|all|coverage|benchmark|clean}"
        echo "  unit        - Run unit tests only"
        echo "  integration - Run integration tests only"
        echo "  all         - Run all tests"
        echo "  coverage    - Run tests with coverage report"
        echo "  benchmark   - Run benchmark tests"
        echo "  clean       - Clean test artifacts"
        exit 1
        ;;
esac