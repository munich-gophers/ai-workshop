#!/bin/bash
# Quick test requests for local development

BASE_URL="http://localhost:8080"

echo "Testing AI Code Mentor API"
echo "=========================="
echo ""

# Test 1: Health check
echo "1. Health Check:"
curl -s "${BASE_URL}/health" | python3 -m json.tool
echo ""
echo ""

# Test 2: Simple code review
echo "2. Simple Code Review:"
curl -s -X POST "${BASE_URL}/api/review" \
  -H "Content-Type: application/json" \
  -d '{
    "diff": "+func add(a, b int) int {\n+    return a + b\n+}",
    "language": "go",
    "file_path": "math.go"
  }' | python3 -m json.tool
echo ""
echo ""

# Test 3: Code with secrets
echo "3. Code with Secrets:"
curl -s -X POST "${BASE_URL}/api/review" \
  -H "Content-Type: application/json" \
  -d '{
    "diff": "+const API_KEY = \"sk-1234567890abcdef\"\n+const PASSWORD = \"admin123\"",
    "language": "go",
    "file_path": "config.go"
  }' | python3 -m json.tool
echo ""
echo ""

# Test 4: GitHub webhook
echo "4. GitHub Webhook:"
curl -s -X POST "${BASE_URL}/webhook/github" \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: pull_request" \
  -d @github-pr-payload.json | python3 -m json.tool
echo ""
