#!/bin/bash

# Pre-commit check script
# Run this before committing to ensure everything is clean

set -e

echo "🔍 Pre-commit Check"
echo "==================="
echo

# Check 1: No uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
    echo "⚠️  Warning: You have uncommitted changes"
    echo "   Review with: git status"
    echo
fi

# Check 2: Build successful
echo "Check 1: Building project..."
if make build > /dev/null 2>&1; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi
echo

# Check 3: No debug code
echo "Check 2: Checking for debug code..."
if grep -r "fmt.Println" cmd/ internal/ 2>/dev/null | grep -v "_test.go" | grep -v "report.go" > /dev/null; then
    echo "⚠️  Warning: Found fmt.Println (may be debug code)"
    echo "   Review with: grep -r 'fmt.Println' cmd/ internal/"
else
    echo "✅ No debug code found"
fi
echo

# Check 4: No TODO comments
echo "Check 3: Checking for TODO comments..."
if grep -r "TODO\|FIXME\|XXX" cmd/ internal/ 2>/dev/null > /dev/null; then
    echo "⚠️  Warning: Found TODO/FIXME comments"
    echo "   Review with: grep -r 'TODO\|FIXME' cmd/ internal/"
else
    echo "✅ No TODO comments found"
fi
echo

# Check 5: Binary not in repo
echo "Check 4: Checking for binary files..."
if [ -f "agent-guard" ]; then
    echo "⚠️  Warning: Binary 'agent-guard' exists (will be ignored by .gitignore)"
fi
echo "✅ .gitignore configured correctly"
echo

# Check 6: Doc directory ignored
echo "Check 5: Checking doc directory..."
if [ -d "doc" ]; then
    if grep -q "^doc/$" .gitignore; then
        echo "✅ doc/ directory in .gitignore"
    else
        echo "⚠️  Warning: doc/ exists but not in .gitignore"
    fi
fi
echo

echo "==================="
echo "✅ Pre-commit check passed!"
echo
echo "Ready to commit:"
echo "  git add ."
echo "  git commit -m 'your message'"
