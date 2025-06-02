#!/bin/bash
set -e

# Print colorful messages
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}Project-Website Lint Script${NC}"
echo -e "${BLUE}========================================${NC}"

# Navigate to backend directory
cd "$(dirname "${BASH_SOURCE[0]}")/.."
BACKEND_DIR="$(pwd)"
echo -e "${GREEN}Backend directory:${NC} ${YELLOW}${BACKEND_DIR}${NC}"

# Check for required tools
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go first.${NC}"
    echo -e "  Visit ${YELLOW}https://go.dev/doc/install${NC} for installation instructions."
    exit 1
fi

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo -e "${YELLOW}golangci-lint not found. Installing...${NC}"
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    if ! command -v golangci-lint &> /dev/null; then
        echo -e "${RED}Failed to install golangci-lint. Please install it manually:${NC}"
        echo -e "  ${YELLOW}go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest${NC}"
        echo -e "  Or visit ${YELLOW}https://golangci-lint.run/usage/install/${NC} for other installation methods."
        exit 1
    fi
fi

echo -e "\n${GREEN}Step 1:${NC} Running go fmt..."
go_fmt_output=$(go fmt ./...)
if [ -n "$go_fmt_output" ]; then
    echo -e "${YELLOW}The following files were reformatted:${NC}"
    echo "$go_fmt_output"
else
    echo -e "${GREEN}All files are properly formatted.${NC}"
fi

echo -e "\n${GREEN}Step 2:${NC} Running go vet..."
go_vet_output=$(go vet ./... 2>&1)
if [ $? -ne 0 ]; then
    echo -e "${RED}go vet found issues:${NC}"
    echo "$go_vet_output"
    exit 1
else
    echo -e "${GREEN}go vet passed with no issues.${NC}"
fi

echo -e "\n${GREEN}Step 3:${NC} Running staticcheck..."
if ! command -v staticcheck &> /dev/null; then
    echo -e "${YELLOW}staticcheck not found. Installing...${NC}"
    go install honnef.co/go/tools/cmd/staticcheck@latest
    if ! command -v staticcheck &> /dev/null; then
        echo -e "${YELLOW}Failed to install staticcheck. Skipping staticcheck step.${NC}"
    fi
fi

if command -v staticcheck &> /dev/null; then
    staticcheck_output=$(staticcheck ./... 2>&1)
    if [ $? -ne 0 ]; then
        echo -e "${RED}staticcheck found issues:${NC}"
        echo "$staticcheck_output"
        exit 1
    else
        echo -e "${GREEN}staticcheck passed with no issues.${NC}"
    fi
fi

echo -e "\n${GREEN}Step 4:${NC} Running golangci-lint..."
golangci_output=$(golangci-lint run ./... 2>&1)
if [ $? -ne 0 ]; then
    echo -e "${RED}golangci-lint found issues:${NC}"
    echo "$golangci_output"
    exit 1
else
    echo -e "${GREEN}golangci-lint passed with no issues.${NC}"
fi

echo -e "\n${GREEN}All lint checks passed successfully!${NC}"
