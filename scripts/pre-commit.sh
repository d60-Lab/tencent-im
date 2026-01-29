#!/bin/bash

# Pre-commit hook for Go projects
# This script runs before every commit to ensure code quality

set -e

echo "🔍 Running pre-commit checks..."

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get list of staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -z "$STAGED_GO_FILES" ]; then
    echo "✅ No Go files staged, skipping checks"
    exit 0
fi

echo "📝 Staged Go files:"
echo "$STAGED_GO_FILES"

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    exit 1
fi

# 1. Format check
echo ""
echo "🎨 Checking formatting..."
UNFORMATTED=$(gofmt -l $STAGED_GO_FILES 2>&1)
if [ -n "$UNFORMATTED" ]; then
    echo -e "${RED}❌ Some files are not formatted:${NC}"
    echo "$UNFORMATTED"
    echo ""
    echo "Run: gofmt -w $UNFORMATTED"
    exit 1
fi
echo -e "${GREEN}✅ All files are properly formatted${NC}"

# 2. Go imports check
echo ""
echo "📦 Checking imports..."
if command -v goimports &> /dev/null; then
    UNFORMATTED=$(goimports -l $STAGED_GO_FILES 2>&1)
    if [ -n "$UNFORMATTED" ]; then
        echo -e "${YELLOW}⚠️  Some files have unorganized imports:${NC}"
        echo "$UNFORMATTED"
        echo ""
        echo "Run: goimports -w $UNFORMATTED"
        # Don't fail, just warn
    else
        echo -e "${GREEN}✅ Imports are organized${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  goimports not found, skipping import check${NC}"
    echo "Install: go install golang.org/x/tools/cmd/goimports@latest"
fi

# 3. Go vet
echo ""
echo "🔬 Running go vet..."
if ! go vet ./... 2>&1; then
    echo -e "${RED}❌ go vet found issues${NC}"
    exit 1
fi
echo -e "${GREEN}✅ go vet passed${NC}"

# 4. Go mod tidy check
echo ""
echo "📋 Checking go.mod..."
if ! git diff --cached --name-only | grep -q "go.mod"; then
    # go.mod not staged, check if it needs tidying
    go mod tidy
    if ! git diff --exit-code go.mod go.sum > /dev/null 2>&1; then
        echo -e "${RED}❌ go.mod or go.sum needs tidying${NC}"
        echo "Run: go mod tidy"
        exit 1
    fi
fi
echo -e "${GREEN}✅ go.mod is tidy${NC}"

# 5. Build check
echo ""
echo "🔨 Building project..."
if ! go build -v ./... > /dev/null 2>&1; then
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Build successful${NC}"

# 6. Run tests (optional, can be slow)
if [ "${SKIP_TESTS}" != "1" ]; then
    echo ""
    echo "🧪 Running tests..."
    if ! go test -short ./... > /dev/null 2>&1; then
        echo -e "${RED}❌ Tests failed${NC}"
        echo "Run: go test ./..."
        exit 1
    fi
    echo -e "${GREEN}✅ Tests passed${NC}"
else
    echo -e "${YELLOW}⚠️  Skipping tests (SKIP_TESTS=1)${NC}"
fi

# 7. golangci-lint (if installed)
if command -v golangci-lint &> /dev/null; then
    echo ""
    echo "🔍 Running golangci-lint..."
    if ! golangci-lint run --timeout=5m ./... > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  golangci-lint found issues${NC}"
        echo "Run: golangci-lint run ./..."
        # Don't fail, just warn for now
    else
        echo -e "${GREEN}✅ golangci-lint passed${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  golangci-lint not found, skipping lint check${NC}"
    echo "Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
fi

echo ""
echo -e "${GREEN}🎉 All pre-commit checks passed!${NC}"
exit 0
