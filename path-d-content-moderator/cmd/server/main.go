package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/munich-gophers/ai-workshop/content-moderator/internal/handler"
	"github.com/munich-gophers/ai-workshop/content-moderator/internal/moderator"
)

// TODO: CHECKPOINT 1 - Pattern-based sentiment analysis
// - Implement health check endpoint
// - Create handler for /api/analyze-sentiment using pattern-based analysis
// - Use internal/analyzer package for sentiment detection

// TODO: CHECKPOINT 2 - AI-powered content moderation
// - Initialize AI moderator with Genkit
// - Create handler for /api/moderate using AI for content safety
// - Detect categories: spam, harassment, hate-speech, violence, etc.

// TODO: CHECKPOINT 3 - Comprehensive analysis with recommendations
// - Create handler for /api/analyze-content combining sentiment + moderation
// - Implement action recommendations (approve/flag/reject/escalate)
// - Add confidence thresholds for auto-execution

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// ✅ CHECKPOINT 1 - Health check endpoint implemented
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"status":   "healthy",
			"service":  "content-moderator",
			"version":  "1.0.0",
			"features": []string{"sentiment-analysis", "content-moderation", "action-recommendations"},
		}

		json.NewEncoder(w).Encode(response)
	})

	// ✅ CHECKPOINT 1 - Sentiment analysis endpoint implemented
	mux.HandleFunc("/api/analyze-sentiment", handler.HandleAnalyzeSentiment)

	// ✅ CHECKPOINT 2 - AI Moderator initialized
	ctx := context.Background()
	aiModerator, err := moderator.New(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize moderator: %v", err)
	}
	log.Println("✅ AI Moderator initialized with Gemini 2.5-flash")

	// ✅ CHECKPOINT 2 - Content moderation endpoint implemented
	mux.HandleFunc("/api/moderate", handler.HandleModerate(aiModerator))

	// ✅ CHECKPOINT 3 - Comprehensive analysis endpoint implemented
	mux.HandleFunc("/api/analyze-content", handler.HandleAnalyzeContent(aiModerator))

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("✅ Health check: http://localhost:%s/health", port)
	log.Printf("✅ Sentiment analysis: http://localhost:%s/api/analyze-sentiment", port)
	log.Printf("✅ Content moderation: http://localhost:%s/api/moderate", port)
	log.Printf("✅ Comprehensive analysis: http://localhost:%s/api/analyze-content", port)
	log.Printf("🎉 All checkpoints complete!")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
