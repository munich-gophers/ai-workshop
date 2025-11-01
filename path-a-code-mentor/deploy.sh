#!/bin/bash
# Deploy to Cloud Run

set -e

PROJECT_ID=$(gcloud config get-value project)
SERVICE_NAME="code-mentor"
REGION="us-central1"
IMAGE_NAME="gcr.io/${PROJECT_ID}/${SERVICE_NAME}"

echo "🚀 Deploying AI Code Mentor to Cloud Run..."
echo "Project: $PROJECT_ID"
echo "Region: $REGION"
echo ""

# Check for required secrets
if [ -z "$GEMINI_API_KEY" ]; then
    echo "❌ GEMINI_API_KEY not set in environment"
    exit 1
fi

# Build and push Docker image
echo "📦 Building Docker image..."
docker build -t ${IMAGE_NAME}:latest .

echo "⬆️  Pushing to Container Registry..."
docker push ${IMAGE_NAME}:latest

echo "🌐 Deploying to Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image ${IMAGE_NAME}:latest \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --set-env-vars GEMINI_API_KEY=$GEMINI_API_KEY \
  --memory 512Mi \
  --cpu 1 \
  --timeout 60s \
  --min-instances 0 \
  --max-instances 10 \
  --port 8080

# Get the URL
URL=$(gcloud run services describe $SERVICE_NAME \
  --region $REGION \
  --format 'value(status.url)')

echo ""
echo "✅ Deployment complete!"
echo "🌐 Service URL: $URL"
echo ""
echo "📝 Update your GitHub webhook settings:"
echo "   Payload URL: ${URL}/webhook/github"
echo "   Content type: application/json"
echo "   Events: Pull requests"
echo ""
echo "Test it:"
echo "curl ${URL}/health"
