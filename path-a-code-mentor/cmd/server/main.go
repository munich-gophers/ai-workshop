package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/munich-gophers/ai-workshop/code-mentor/internal/analyzer"
	"github.com/munich-gophers/ai-workshop/code-mentor/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// ✅ CHECKPOINT 1: Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"status":              "healthy",
			"service":             "code-mentor",
			"version":             "1.0.0",
			"supported_platforms": []string{"github", "gitlab"},
		}

		json.NewEncoder(w).Encode(response)
	})

	// ✅ CHECKPOINT 2: Initialize analyzer
	ctx := context.Background()
	codeAnalyzer, err := analyzer.New(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to initialize analyzer: %v", err)
	}
	log.Println("✅ Analyzer initialized with Gemini")

	// ✅ CHECKPOINT 2: Add /api/review endpoint
	mux.HandleFunc("/api/review", handler.HandleReview(codeAnalyzer))

	// TODO: CHECKPOINT 3 - Add /webhook/github endpoint
	//
	// Goal: Handle GitHub PR webhooks
	//
	// Hint: Add this code here:
	//   // Configure GitHub webhook
	//   webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	//   if webhookSecret == "" {
	//       log.Println("⚠️  GITHUB_WEBHOOK_SECRET not set - webhook signature validation disabled")
	//   }
	//
	//   // Create GitHub client
	//   githubClient := github.NewClient(nil)
	//
	//   // Configure webhook handler
	//   webhookConfig := &webhook.Config{
	//       Secret:       webhookSecret,
	//       GitHubClient: githubClient,
	//   }
	//
	//   // Add webhook endpoint
	//   mux.HandleFunc("/webhook/github", webhook.HandleGitHub(codeAnalyzer, webhookConfig))
	//
	// Note: You'll need to implement HandleGitHub() in internal/webhook/github.go
	// Note: Don't forget to import "github.com/google/go-github/v76/github"

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("✅ Health check: http://localhost:%s/health", port)
	log.Printf("✅ Code review: http://localhost:%s/api/review", port)
	log.Printf("💡 Next: Implement webhook handler (Checkpoint 3)")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
