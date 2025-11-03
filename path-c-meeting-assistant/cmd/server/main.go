package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/munich-gophers/ai-workshop/meeting-assistant/internal/classifier"
	"github.com/munich-gophers/ai-workshop/meeting-assistant/internal/handler"
)

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
			"service":  "meeting-assistant",
			"version":  "1.0.0",
			"features": []string{"extract-actions", "analyze-meeting", "generate-email"},
		}

		json.NewEncoder(w).Encode(response)
	})

	// ✅ CHECKPOINT 1 - Action item extraction endpoint implemented
	mux.HandleFunc("/api/extract", handler.HandleExtract)

	// ✅ CHECKPOINT 2 - AI Classifier initialized
	ctx := context.Background()
	aiClassifier, err := classifier.New(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize classifier: %v", err)
	}
	log.Println("✅ AI Classifier initialized with Gemini 2.5-flash")

	// ✅ CHECKPOINT 2 - Meeting analysis endpoint implemented
	mux.HandleFunc("/api/analyze", handler.HandleAnalyze(aiClassifier))

	// ✅ CHECKPOINT 3 - Email generation endpoint implemented
	mux.HandleFunc("/api/generate-email", handler.HandleGenerateEmail(aiClassifier))

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("✅ Health check: http://localhost:%s/health", port)
	log.Printf("✅ Extract endpoint: http://localhost:%s/api/extract", port)
	log.Printf("✅ Analyze endpoint: http://localhost:%s/api/analyze", port)
	log.Printf("✅ Email generator: http://localhost:%s/api/generate-email", port)
	log.Printf("🎉 All checkpoints complete!")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
