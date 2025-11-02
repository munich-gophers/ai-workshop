package classifier

import (
	"context"

	"github.com/munich-gophers/ai-workshop/support-agent/internal/models"
)

// Classifier performs AI-powered customer message triage using Genkit v1.1.0 API
type Classifier struct {
	// TODO: CHECKPOINT 2 - Add fields
	//
	// You'll need:
	// - genkit *genkit.Genkit
	// - model ai.Model
	// - basePrompt string
	// - intentPrompt string
	// - urgencyPrompt string
}

// New creates a new Classifier instance using Genkit v1.1.0 API
func New(ctx context.Context) (*Classifier, error) {
	// TODO: CHECKPOINT 2 - Initialize Genkit v1.1.0
	//
	// Steps:
	// 1. Import "github.com/firebase/genkit/go/genkit"
	// 2. Import "github.com/firebase/genkit/go/ai"
	// 3. Import "github.com/firebase/genkit/go/plugins/googlegenai"
	// 4. Get API key from environment: apiKey := os.Getenv("GEMINI_API_KEY")
	// 5. Create Google AI plugin: googleAI := &googlegenai.GoogleAI{APIKey: apiKey}
	// 6. Initialize Genkit: g := genkit.Init(ctx, genkit.WithPlugins(googleAI))
	// 7. Get model: model := googlegenai.GoogleAIModel(g, "gemini-2.5-flash")
	// 8. Load prompts from files using loadPrompt()
	// 9. Return &Classifier with all fields

	return nil, nil
}

// Triage analyzes a customer message and returns classification
func (c *Classifier) Triage(ctx context.Context, req models.TriageRequest) (*models.TriageResponse, error) {
	// TODO: CHECKPOINT 2 - Implement triage logic
	//
	// Steps:
	// 1. Build prompt combining base + intent + urgency prompts with customer message
	// 2. Call Gemini using c.model.Generate()
	// 3. Parse AI response to extract intent, urgency, and summary
	// 4. Return TriageResponse with all fields populated
	//
	// TODO: CHECKPOINT 3 - Add PII handling
	//
	// Before calling AI:
	// 1. Detect PII: detectedPII := redactor.DetectPII(req.Message)
	// 2. Redact PII: redactedMessage := redactor.RedactPII(req.Message)
	// 3. Use redactedMessage for AI analysis
	// 4. Include detectedPII in response

	return nil, nil
}

// loadPrompt loads a prompt from file
func loadPrompt(filename string) (string, error) {
	// TODO: CHECKPOINT 2 - Implement prompt loading
	// Read file from internal/prompts/{filename}.txt
	return "", nil
}
