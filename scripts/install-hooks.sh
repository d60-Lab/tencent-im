#!/bin/bash

# Install pre-commit hooks for this repository
# Run this script once after cloning the repository

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
HOOK_DIR="$REPO_ROOT/.git/hooks"
PRE_COMMIT_HOOK="$HOOK_DIR/pre-commit"
SCRIPT_PATH="$REPO_ROOT/scripts/pre-commit.sh"

echo "🔧 Installing pre-commit hooks..."

# Create hooks directory if it doesn't exist
mkdir -p "$HOOK_DIR"

# Check if pre-commit hook already exists
if [ -f "$PRE_COMMIT_HOOK" ] && [ ! -L "$PRE_COMMIT_HOOK" ]; then
    echo "⚠️  Existing pre-commit hook found"
    read -p "Do you want to overwrite it? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Installation cancelled"
        exit 1
    fi
    rm "$PRE_COMMIT_HOOK"
fi

# Make the script executable
chmod +x "$SCRIPT_PATH"

# Create symlink to the pre-commit script
ln -sf "../../scripts/pre-commit.sh" "$PRE_COMMIT_HOOK"

echo "✅ Pre-commit hook installed successfully!"
echo ""
echo "📝 Usage:"
echo "  - Hooks will run automatically on 'git commit'"
echo "  - Skip hooks with: git commit --no-verify"
echo "  - Skip tests with: SKIP_TESTS=1 git commit"
echo ""
echo "🔍 Optional tools to install for better checks:"
echo "  - goimports: go install golang.org/x/tools/cmd/goimports@latest"
echo "  - golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
echo ""
echo "🎉 You're all set!"
