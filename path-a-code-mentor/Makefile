.PHONY: help build run stop clean test deploy dev

IMAGE_NAME := code-mentor
CONTAINER_NAME := code-mentor
PORT := 8080

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Docker image
	docker build -t $(IMAGE_NAME):latest .

run: ## Run container
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p $(PORT):8080 \
		-e GEMINI_API_KEY="${GEMINI_API_KEY}" \
		-e GITHUB_WEBHOOK_SECRET="${GITHUB_WEBHOOK_SECRET}" \
		$(IMAGE_NAME):latest

stop: ## Stop and remove container
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

logs: ## View container logs
	docker logs -f $(CONTAINER_NAME)

shell: ## Open shell in container
	docker exec -it $(CONTAINER_NAME) /bin/sh

clean: stop ## Clean up images and containers
	docker rmi $(IMAGE_NAME):latest || true

rebuild: clean build ## Rebuild from scratch

test: ## Test the service
	@echo "Testing health endpoint..."
	@curl -s http://localhost:$(PORT)/health | python3 -m json.tool || echo "Server not running or health check not implemented"

deploy: ## Deploy to Cloud Run
	./deploy.sh

dev: ## Run locally without Docker
	go run cmd/server/main.go
