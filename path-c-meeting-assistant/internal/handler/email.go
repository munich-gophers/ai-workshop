package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/munich-gophers/ai-workshop/meeting-assistant/internal/classifier"
	"github.com/munich-gophers/ai-workshop/meeting-assistant/internal/models"
)

// HandleGenerateEmail handles POST /api/generate-email requests
func HandleGenerateEmail(c *classifier.Classifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode request
		var req models.EmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding request: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate notes are not empty
		if req.Notes == "" {
			http.Error(w, "Notes field is required", http.StatusBadRequest)
			return
		}

		// Generate email using AI
		resp, err := c.GenerateEmail(r.Context(), req)
		if err != nil {
			log.Printf("Error generating email: %v", err)
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
