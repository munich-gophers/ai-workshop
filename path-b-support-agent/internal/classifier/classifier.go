package classifier

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/munich-gophers/ai-workshop/support-agent/internal/models"
	"github.com/munich-gophers/ai-workshop/support-agent/internal/redactor"
)

// Classifier performs AI-powered customer message triage using Genkit v1.1.0 API
type Classifier struct {
	genkit        *genkit.Genkit
	basePrompt    string
	intentPrompt  string
	urgencyPrompt string
}

// New creates a new Classifier instance using Genkit v1.1.0 API
func New(ctx context.Context) (*Classifier, error) {
	// Get API key from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// Create Google AI plugin
	googleAI := &googlegenai.GoogleAI{APIKey: apiKey}

	// Initialize Genkit with the plugin - returns just the genkit instance
	g := genkit.Init(ctx,
		genkit.WithPlugins(googleAI),
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
	)

	// Load prompts
	basePrompt, err := loadPrompt("base")
	if err != nil {
		return nil, fmt.Errorf("failed to load base prompt: %w", err)
	}

	intentPrompt, err := loadPrompt("intent")
	if err != nil {
		return nil, fmt.Errorf("failed to load intent prompt: %w", err)
	}

	urgencyPrompt, err := loadPrompt("urgency")
	if err != nil {
		return nil, fmt.Errorf("failed to load urgency prompt: %w", err)
	}

	return &Classifier{
		genkit:        g,
		basePrompt:    basePrompt,
		intentPrompt:  intentPrompt,
		urgencyPrompt: urgencyPrompt,
	}, nil
}

// Triage analyzes a customer message and returns classification
func (c *Classifier) Triage(ctx context.Context, req models.TriageRequest) (*models.TriageResponse, error) {
	startTime := time.Now()

	// ✅ CHECKPOINT 3 - Detect PII before AI processing
	detectedPII := redactor.DetectPII(req.Message)

	// ✅ CHECKPOINT 3 - Redact PII from message before sending to AI
	redactedMessage := redactor.RedactPII(req.Message)

	// Build the complete prompt with redacted message
	fullPrompt := fmt.Sprintf("%s\n\n%s\n\n%s\n\nCustomer Message:\n%s\n\nProvide your analysis in JSON format with the following structure:\n{\n  \"intent\": {\n    \"category\": \"billing|technical|account|general\",\n    \"confidence\": 0.95,\n    \"subcategory\": \"optional subcategory\"\n  },\n  \"urgency\": {\n    \"level\": \"critical|high|medium|low\",\n    \"confidence\": 0.90,\n    \"reason\": \"explanation for urgency level\"\n  },\n  \"summary\": \"Brief summary of the customer's request\",\n  \"suggested_routing\": \"Recommendation for which team should handle this\"\n}",
		c.basePrompt,
		c.intentPrompt,
		c.urgencyPrompt,
		redactedMessage, // Use redacted message for AI
	)

	// Call Gemini using Generate
	response, err := genkit.Generate(ctx, c.genkit,
		ai.WithPrompt(fullPrompt),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	// Extract text from response
	responseText := response.Text()
	if responseText == "" {
		return nil, fmt.Errorf("empty response from AI model")
	}

	// Parse JSON response
	// Clean the response text (remove markdown code blocks if present)
	responseText = strings.TrimSpace(responseText)
	if strings.HasPrefix(responseText, "```json") {
		responseText = strings.TrimPrefix(responseText, "```json")
		responseText = strings.TrimSuffix(responseText, "```")
		responseText = strings.TrimSpace(responseText)
	} else if strings.HasPrefix(responseText, "```") {
		responseText = strings.TrimPrefix(responseText, "```")
		responseText = strings.TrimSuffix(responseText, "```")
		responseText = strings.TrimSpace(responseText)
	}

	var aiResponse struct {
		Intent struct {
			Category    string  `json:"category"`
			Confidence  float64 `json:"confidence"`
			Subcategory string  `json:"subcategory,omitempty"`
		} `json:"intent"`
		Urgency struct {
			Level      string  `json:"level"`
			Confidence float64 `json:"confidence"`
			Reason     string  `json:"reason,omitempty"`
		} `json:"urgency"`
		Summary          string `json:"summary"`
		SuggestedRouting string `json:"suggested_routing"`
	}

	if err := json.Unmarshal([]byte(responseText), &aiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w\nResponse: %s", err, responseText)
	}

	// Build response
	processingTime := time.Since(startTime).Milliseconds()

	return &models.TriageResponse{
		OriginalMessage: req.Message,
		RedactedMessage: redactedMessage, // ✅ CHECKPOINT 3 - Include redacted message
		Intent: models.Intent{
			Category:    aiResponse.Intent.Category,
			Confidence:  aiResponse.Intent.Confidence,
			Subcategory: aiResponse.Intent.Subcategory,
		},
		Urgency: models.Urgency{
			Level:      aiResponse.Urgency.Level,
			Confidence: aiResponse.Urgency.Confidence,
			Reason:     aiResponse.Urgency.Reason,
		},
		Summary:          aiResponse.Summary,
		DetectedPII:      detectedPII, // ✅ CHECKPOINT 3 - Include detected PII
		SuggestedRouting: aiResponse.SuggestedRouting,
		ProcessingTimeMs: processingTime,
	}, nil
}

// loadPrompt loads a prompt from file
func loadPrompt(filename string) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("internal/prompts/%s.txt", filename))
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file: %w", err)
	}
	return string(data), nil
}
