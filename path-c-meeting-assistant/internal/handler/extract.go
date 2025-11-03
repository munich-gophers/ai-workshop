package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/munich-gophers/ai-workshop/meeting-assistant/internal/analyzer"
	"github.com/munich-gophers/ai-workshop/meeting-assistant/internal/models"
)

// HandleExtract handles POST /api/extract requests
func HandleExtract(w http.ResponseWriter, r *http.Request) {
	// Verify POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request
	var req models.MeetingNotesRequest
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

	// Extract action items using pattern matching
	startTime := time.Now()
	actionItems := analyzer.ExtractActionItems(req.Notes)
	processingTime := time.Since(startTime).Milliseconds()

	// Generate brief summary (first 200 characters of notes)
	summary := req.Notes
	if len(summary) > 200 {
		summary = summary[:200] + "..."
	}

	// Build response
	resp := models.ExtractResponse{
		ActionItems:      actionItems,
		TotalActions:     len(actionItems),
		ProcessingTimeMs: processingTime,
		Summary:          summary,
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
