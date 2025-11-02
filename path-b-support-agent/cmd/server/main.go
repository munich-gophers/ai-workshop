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

	// TODO: CHECKPOINT 1 - Add health check endpoint
	//
	// Goal: Return JSON indicating service is healthy
	//
	// Hint: Add this code here:
	//   mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	//       w.Header().Set("Content-Type", "application/json")
	//       w.WriteHeader(http.StatusOK)
	//
	//       response := map[string]interface{}{
	//           "status":    "healthy",
	//           "service":   "support-agent",
	//           "version":   "1.0.0",
	//           "features":  []string{"triage", "pii-redaction"},
	//       }
	//
	//       json.NewEncoder(w).Encode(response)
	//   })
	//
	// Test: curl http://localhost:8080/health

	// TODO: CHECKPOINT 2 - Initialize AI classifier
	//
	// Goal: Create a classifier instance that connects to Gemini
	//
	// Hint: Add this code here:
	//   ctx := context.Background()
	//   aiClassifier, err := classifier.New(ctx)
	//   if err != nil {
	//       log.Fatalf("Failed to initialize classifier: %v", err)
	//   }
	//   log.Println("✅ AI Classifier initialized with Gemini")
	//
	// Note: You'll also need to implement classifier.New() in internal/classifier/classifier.go

	// TODO: CHECKPOINT 2 - Add /api/triage endpoint
	//
	// Goal: Accept customer messages and return AI triage analysis
	//
	// Hint: Add this code here:
	//   mux.HandleFunc("/api/triage", handler.HandleTriage(aiClassifier))
	//
	// Note: You'll need to implement HandleTriage() in internal/handler/triage.go

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
	log.Printf("💡 Next: Implement health endpoint (Checkpoint 1)")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
