package handler

import (
	"net/http"

	"github.com/munich-gophers/ai-workshop/support-agent/internal/classifier"
)

// HandleTriage handles POST /api/triage requests
func HandleTriage(c *classifier.Classifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: CHECKPOINT 2 - Implement triage endpoint
		//
		// Steps:
		// 1. Check request method is POST
		// 2. Decode JSON request body into models.TriageRequest
		// 3. Call c.Triage(r.Context(), request)
		// 4. Handle any errors appropriately
		// 5. Return JSON response with models.TriageResponse
		//
		// Example:
		//   var req models.TriageRequest
		//   if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//       http.Error(w, "Invalid request", http.StatusBadRequest)
		//       return
		//   }
		//
		//   startTime := time.Now()
		//   resp, err := c.Triage(r.Context(), req)
		//   if err != nil {
		//       http.Error(w, err.Error(), http.StatusInternalServerError)
		//       return
		//   }
		//
		//   resp.ProcessingTimeMs = time.Since(startTime).Milliseconds()
		//
		//   w.Header().Set("Content-Type", "application/json")
		//   json.NewEncoder(w).Encode(resp)

		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}
}
