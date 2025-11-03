package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/munich-gophers/ai-workshop/content-moderator/internal/models"
	"github.com/munich-gophers/ai-workshop/content-moderator/internal/moderator"
)

// HandleModerate handles POST /api/moderate requests
func HandleModerate(m *moderator.Moderator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Moderate content using AI
		resp, err := m.Moderate(r.Context(), req)
		if err != nil {
			log.Printf("Error moderating content: %v", err)
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
