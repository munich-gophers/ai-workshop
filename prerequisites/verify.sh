#!/bin/bash
# Verify workshop prerequisites

set +e  # Don't exit on errors, we want to check everything

echo "🔍 Verifying Workshop Prerequisites"
echo "===================================="
echo ""

ERRORS=0

# Check Go
echo -n "Checking Go installation... "
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo "✅ $GO_VERSION"
else
    echo "❌ Go not found"
    echo "   Install from: https://go.dev/doc/install"
    ERRORS=$((ERRORS + 1))
fi

# Check Go version
echo -n "Checking Go version... "
GO_VERSION_NUM=$(go version 2>/dev/null | grep -oE 'go[0-9]+\.[0-9]+' | grep -oE '[0-9]+\.[0-9]+')
if [ ! -z "$GO_VERSION_NUM" ]; then
    MAJOR=$(echo $GO_VERSION_NUM | cut -d. -f1)
    MINOR=$(echo $GO_VERSION_NUM | cut -d. -f2)
    if [ "$MAJOR" -ge 1 ] && [ "$MINOR" -ge 23 ]; then
        echo "✅ Go $GO_VERSION_NUM (>= 1.23)"
    else
        echo "⚠️  Go $GO_VERSION_NUM (need >= 1.23)"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "❌ Could not determine Go version"
    ERRORS=$((ERRORS + 1))
fi

# Check gcloud
echo -n "Checking gcloud CLI... "
if command -v gcloud &> /dev/null; then
    echo "✅ Installed"
else
    echo "❌ gcloud not found"
    echo "   Install from: https://cloud.google.com/sdk/docs/install"
    ERRORS=$((ERRORS + 1))
fi

# Check gcloud authentication
echo -n "Checking gcloud authentication... "
GCLOUD_ACCOUNT=$(gcloud config get-value account 2>/dev/null)
if [ ! -z "$GCLOUD_ACCOUNT" ]; then
    echo "✅ $GCLOUD_ACCOUNT"
else
    echo "❌ Not authenticated"
    echo "   Run: gcloud auth login"
    ERRORS=$((ERRORS + 1))
fi

# Check Gemini API key
echo -n "Checking GEMINI_API_KEY... "
if [ ! -z "$GEMINI_API_KEY" ]; then
    # Show first 10 chars only
    KEY_PREVIEW="${GEMINI_API_KEY:0:10}..."
    echo "✅ Set ($KEY_PREVIEW)"
else
    echo "❌ Not set"
    echo "   Get key from: https://aistudio.google.com/app/apikey"
    echo "   Set with: export GEMINI_API_KEY=your-key-here"
    ERRORS=$((ERRORS + 1))
fi

# Check Git
echo -n "Checking Git... "
if command -v git &> /dev/null; then
    GIT_VERSION=$(git --version | awk '{print $3}')
    echo "✅ $GIT_VERSION"
else
    echo "❌ Git not found"
    echo "   Install from: https://git-scm.com/downloads"
    ERRORS=$((ERRORS + 1))
fi

# Optional: Check Docker
echo -n "Checking Docker (optional)... "
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | tr -d ',')
    echo "✅ $DOCKER_VERSION"
else
    echo "⚠️  Not installed (optional for workshop)"
fi

echo ""
echo "===================================="
if [ $ERRORS -eq 0 ]; then
    echo "✅ All prerequisites met!"
    echo "🚀 You're ready for the workshop!"
else
    echo "❌ Found $ERRORS issue(s)"
    echo "⚠️  Please fix the issues above before the workshop"
    exit 1
fi
