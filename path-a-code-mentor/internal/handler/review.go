package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/munich-gophers/ai-workshop/code-mentor/internal/analyzer"
	"github.com/munich-gophers/ai-workshop/code-mentor/internal/models"
)

// HandleReview processes code review requests
func HandleReview(a *analyzer.Analyzer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only accept POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse request
		var req models.ReviewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate
		if req.Diff == "" {
			http.Error(w, "Diff is required", http.StatusBadRequest)
			return
		}
		if req.FilePath == "" {
			http.Error(w, "File path is required", http.StatusBadRequest)
			return
		}

		// Track processing time
		start := time.Now()

		// Analyze code
		result, err := a.Review(r.Context(), req)
		if err != nil {
			log.Printf("Analysis error: %v", err)
			http.Error(w, "Analysis failed", http.StatusInternalServerError)
			return
		}

		// Add processing time
		result.ProcessingTimeMs = int(time.Since(start).Milliseconds())

		// Return result
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
