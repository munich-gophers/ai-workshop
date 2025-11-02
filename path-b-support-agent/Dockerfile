# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies (cached layer)
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/server \
    ./cmd/server

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/server /app/server

# Copy prompts directory (needed at runtime)
COPY --from=builder /build/prompts /app/prompts

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/server"]
```

---

## 🚫 **.dockerignore** (Same for both paths)
```
# path-a-code-mentor/.dockerignore
# path-b-support-agent/.dockerignore

# Git
.git
.gitignore
.gitattributes

# Documentation
README.md
PROGRESS.md
*.md

# IDE
.vscode
.idea
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Build artifacts
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test files
*_test.go
testdata/

# Examples (not needed in container)
examples/

# Scripts
*.sh
Makefile

# Env files (secrets!)
.env
.env.*
!.env.example

# Docker
Dockerfile
.dockerignore
docker-compose.yml
