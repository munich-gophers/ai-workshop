package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/munich-gophers/ai-workshop/support-agent/internal/classifier"
	"github.com/munich-gophers/ai-workshop/support-agent/internal/models"
)

// HandleTriage handles POST /api/triage requests
func HandleTriage(c *classifier.Classifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode request
		var req models.TriageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding request: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate message is not empty
		if req.Message == "" {
			http.Error(w, "Message field is required", http.StatusBadRequest)
			return
		}

		// Call classifier
		resp, err := c.Triage(r.Context(), req)
		if err != nil {
			log.Printf("Error processing triage: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
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
}
