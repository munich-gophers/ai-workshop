package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// TODO: CHECKPOINT 1 - Add health check endpoint
	//
	// Goal: Return JSON indicating service is healthy
	//
	// Hint: Add this code here:
	//   mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	//       w.Header().Set("Content-Type", "application/json")
	//       w.WriteHeader(http.StatusOK)
	//       w.Write([]byte(`{"status":"healthy","service":"code-mentor","version":"1.0.0"}`))
	//   })
	//
	// Test: curl http://localhost:8080/health

	// TODO: CHECKPOINT 2 - Initialize analyzer
	//
	// Goal: Create an analyzer instance that connects to Gemini
	//
	// Hint: Add this code here:
	//   ctx := context.Background()
	//   analyzer, err := analyzer.New(ctx)
	//   if err != nil {
	//       log.Fatalf("Failed to initialize analyzer: %v", err)
	//   }
	//
	// Note: You'll also need to implement analyzer.New() in internal/analyzer/analyzer.go

	// TODO: CHECKPOINT 2 - Add /api/review endpoint
	//
	// Goal: Accept code diffs and return AI reviews
	//
	// Hint: Add this code here:
	//   mux.HandleFunc("/api/review", handler.HandleReview(analyzer))
	//
	// Note: You'll need to implement HandleReview() in internal/handler/review.go

	// TODO: CHECKPOINT 3 - Add /webhook/github endpoint
	//
	// Goal: Handle GitHub PR webhooks
	//
	// Hint: Add this code here:
	//   mux.HandleFunc("/webhook/github", webhook.HandleGitHub(analyzer))
	//
	// Note: You'll need to implement HandleGitHub() in internal/webhook/github.go

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("✅ Health check: http://localhost:%s/health", port)
	log.Printf("💡 Next: Implement health endpoint (Checkpoint 1)")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
