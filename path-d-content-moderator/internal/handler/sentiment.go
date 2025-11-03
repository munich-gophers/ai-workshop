package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/munich-gophers/ai-workshop/content-moderator/internal/analyzer"
	"github.com/munich-gophers/ai-workshop/content-moderator/internal/models"
)

// HandleAnalyzeSentiment handles POST /api/analyze-sentiment requests
func HandleAnalyzeSentiment(w http.ResponseWriter, r *http.Request) {
	// Verify POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request
	var req models.ContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate content is not empty
	if req.Content == "" {
		http.Error(w, "Content field is required", http.StatusBadRequest)
		return
	}

	// Analyze sentiment using pattern-based analysis
	startTime := time.Now()
	sentiment := analyzer.AnalyzeSentiment(req.Content)
	processingTime := time.Since(startTime).Milliseconds()

	// Build response
	resp := models.SentimentResponse{
		ContentID:        req.ContentID,
		Sentiment:        sentiment,
		ProcessingTimeMs: processingTime,
		Method:           "pattern-based",
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
