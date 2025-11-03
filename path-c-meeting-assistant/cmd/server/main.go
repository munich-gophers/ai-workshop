package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

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

	// TODO: CHECKPOINT 2 - Initialize AI Classifier
	// - Import and initialize Genkit with Gemini model
	// - Load prompts from internal/prompts/
	// - Handle GEMINI_API_KEY from environment

	// TODO: CHECKPOINT 2 - Implement full meeting analysis endpoint
	// - Add POST /api/analyze endpoint
	// - Use AI to extract action items, decisions, and participants
	// - Provide comprehensive meeting analysis with topics and summary

	// TODO: CHECKPOINT 3 - Implement follow-up email generation endpoint
	// - Add POST /api/generate-email endpoint
	// - Use AI to generate professional follow-up emails
	// - Support customizable tone (formal, casual, friendly)
	// - Include action items, decisions, and next steps

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("✅ Health check: http://localhost:%s/health", port)
	log.Printf("✅ Extract endpoint: http://localhost:%s/api/extract", port)
	log.Printf("💡 Next: Initialize AI classifier and implement analysis endpoint (Checkpoint 2)")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
