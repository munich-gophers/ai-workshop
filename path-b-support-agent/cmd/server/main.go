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

	// TODO: CHECKPOINT 3 - Add PII redaction to triage flow
	//
	// Goal: Detect and redact PII before AI processing
	//
	// Note: This will be implemented within the triage handler
	// The handler should:
	// 1. Detect PII using redactor.DetectPII()
	// 2. Redact PII using redactor.RedactPII() before sending to AI
	// 3. Include detected PII in the response
	//
	// The redactor package is already implemented in internal/redactor/pii.go

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("✅ Health check: http://localhost:%s/health", port)
	log.Printf("✅ Triage endpoint: http://localhost:%s/api/triage", port)
	log.Printf("💡 Next: Add PII redaction (Checkpoint 3)")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
