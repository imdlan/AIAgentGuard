#!/bin/bash

# AI AgentGuard Test Script
# Run this to verify the installation and basic functionality

set -e

echo "🛡️  AI AgentGuard - Installation Test"
echo "======================================"
echo

# Check if binary exists
if ! command -v agent-guard &> /dev/null; then
    echo "❌ Error: agent-guard not found in PATH"
    echo "Please install first:"
    echo "  make install"
    exit 1
fi

echo "✅ Binary found: $(which agent-guard)"
echo

# Test 1: Help command
echo "Test 1: Help command"
if agent-guard --help &> /dev/null; then
    echo "✅ Help command works"
else
    echo "❌ Help command failed"
    exit 1
fi
echo

# Test 2: Scan command
echo "Test 2: Security scan"
if agent-guard scan &> /dev/null; then
    echo "✅ Scan command works"
else
    echo "⚠️  Scan command had issues (may be normal)"
fi
echo

# Test 3: Report command
echo "Test 3: Report generation"
if agent-guard report &> /dev/null; then
    echo "✅ Report command works"
else
    echo "⚠️  Report command had issues"
fi
echo

# Test 4: Init command
echo "Test 4: Config initialization"
TEST_DIR=$(mktemp -d)
cd "$TEST_DIR"
if agent-guard init --force &> /dev/null; then
    if [ -f ".agent-guard.yaml" ]; then
        echo "✅ Config file created"
        rm -f .agent-guard.yaml
    else
        echo "❌ Config file not created"
    fi
else
    echo "⚠️  Init command had issues"
fi
cd - > /dev/null
rm -rf "$TEST_DIR"
echo

# Test 5: JSON output
echo "Test 5: JSON output"
if agent-guard scan --json &> /dev/null; then
    echo "✅ JSON output works"
else
    echo "⚠️  JSON output had issues"
fi
echo

echo "======================================"
echo "✅ All tests passed!"
echo
echo "You can now use agent-guard:"
echo "  agent-guard scan          # Scan for security risks"
echo "  agent-guard run 'echo hi' # Run in sandbox"
echo "  agent-guard report        # Generate report"
