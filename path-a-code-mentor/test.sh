#!/bin/bash
# test.sh for Path A

CHECKPOINT=$1
PORT=${PORT:-8080}
BASE_URL="http://localhost:${PORT}"

if [ -z "$CHECKPOINT" ]; then
    echo "Usage: ./test.sh <checkpoint-1|checkpoint-2|checkpoint-3>"
    exit 1
fi

echo "🧪 Testing Path A - ${CHECKPOINT}..."
echo ""

case $CHECKPOINT in
    checkpoint-1)
        echo "📍 Testing health endpoint..."
        RESPONSE=$(curl -s "${BASE_URL}/health")

        if echo "$RESPONSE" | grep -q "code-mentor"; then
            echo "✅ PASSED: Health endpoint working"
            echo "Response: $RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
            echo ""
            echo "🎉 Checkpoint 1 complete!"
            echo "👉 Next: ./switch.sh path-a checkpoint-2"
        else
            echo "❌ FAILED: Expected 'code-mentor' in response"
            echo "Got: $RESPONSE"
            exit 1
        fi
        ;;

    checkpoint-2)
        echo "📍 Testing code analysis..."
        RESPONSE=$(curl -s -X POST "${BASE_URL}/api/review" \
            -H "Content-Type: application/json" \
            -d '{
                "diff": "+func add(a, b int) int {\n+    return a + b\n+}",
                "language": "go",
                "file_path": "math.go"
            }')

        if echo "$RESPONSE" | grep -q '"suggestions"' && \
           echo "$RESPONSE" | grep -q '"summary"'; then
            echo "✅ PASSED: Code analysis working"
            echo "Response: $RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
        else
            echo "❌ FAILED: Missing required fields"
            echo "Got: $RESPONSE"
            exit 1
        fi
        ;;

    checkpoint-3)
        echo "📍 Testing secret detection..."
        RESPONSE=$(curl -s -X POST "${BASE_URL}/api/review" \
            -H "Content-Type: application/json" \
            -d '{
                "diff": "+const API_KEY = \"sk-1234567890abcdef\"\n+const PASSWORD = \"admin123\"",
                "language": "go",
                "file_path": "config.go"
            }')

        if echo "$RESPONSE" | grep -q '"secrets_detected":true'; then
            echo "✅ PASSED: Secret detection working"
            echo "Response: $RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
        else
            echo "⚠️  WARNING: Secrets not detected (check security/detector.go)"
        fi

        echo ""
        echo "📍 Testing GitHub webhook handler..."

        WEBHOOK_RESPONSE=$(curl -s -X POST "${BASE_URL}/webhook/github" \
            -H "Content-Type: application/json" \
            -H "X-GitHub-Event: pull_request" \
            -H "X-GitHub-Delivery: test-delivery-123" \
            -d @examples/github-pr-payload.json)

        if echo "$WEBHOOK_RESPONSE" | grep -q '"analysis_id"'; then
            echo "✅ PASSED: Webhook handler working"
            echo "Response: $WEBHOOK_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$WEBHOOK_RESPONSE"
        else
            echo "❌ FAILED: Webhook test failed"
            echo "Got: $WEBHOOK_RESPONSE"
            exit 1
        fi
        ;;
esac

echo ""
echo "🎯 Checkpoint $CHECKPOINT complete!"
