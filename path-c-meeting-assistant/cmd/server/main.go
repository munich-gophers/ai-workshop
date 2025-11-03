package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// TODO: CHECKPOINT 1 - Implement health check endpoint
	// - Add /health endpoint that returns service status
	// - Return JSON with status, service name, version, and features

	// TODO: CHECKPOINT 1 - Implement action item extraction endpoint
	// - Add POST /api/extract endpoint
	// - Accept meeting notes in request body
	// - Use pattern matching to extract action items
	// - Return list of action items with assignees and due dates

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
	log.Printf("💡 Next: Implement health check and extraction endpoints (Checkpoint 1)")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
