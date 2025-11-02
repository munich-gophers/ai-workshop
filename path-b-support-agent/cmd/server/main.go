package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/munich-gophers/ai-workshop/support-agent/internal/classifier"
	"github.com/munich-gophers/ai-workshop/support-agent/internal/handler"
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
			"service":  "support-agent",
			"version":  "1.0.0",
			"features": []string{"triage", "pii-redaction"},
		}

		json.NewEncoder(w).Encode(response)
	})

	// ✅ CHECKPOINT 2 - AI Classifier initialized
	ctx := context.Background()
	aiClassifier, err := classifier.New(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize classifier: %v", err)
	}
	log.Println("✅ AI Classifier initialized with Gemini 2.5-flash")

	// ✅ CHECKPOINT 2 - Triage endpoint implemented
	mux.HandleFunc("/api/triage", handler.HandleTriage(aiClassifier))

	// ✅ CHECKPOINT 3 - PII redaction integrated
	// PII detection and redaction is now active within the classifier
	// - Detects PII before AI processing
	// - Redacts PII from messages sent to Gemini
	// - Includes detected PII in response for logging/compliance

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("✅ Health check: http://localhost:%s/health", port)
	log.Printf("✅ Triage endpoint: http://localhost:%s/api/triage", port)
	log.Printf("✅ PII Detection: Active - protects customer data")
	log.Printf("🎉 All checkpoints complete!")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
