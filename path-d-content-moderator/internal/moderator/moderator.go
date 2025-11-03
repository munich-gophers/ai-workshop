package moderator

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
	"github.com/munich-gophers/ai-workshop/content-moderator/internal/models"
)

// Moderator performs AI-powered content moderation using Genkit v1.1.0 API
type Moderator struct {
	genkit         *genkit.Genkit
	basePrompt     string
	moderatePrompt string
	analyzePrompt  string
}

// New creates a new Moderator instance using Genkit v1.1.0 API
func New(ctx context.Context) (*Moderator, error) {
	// Get API key from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// Create Google AI plugin
	googleAI := &googlegenai.GoogleAI{APIKey: apiKey}

	// Initialize Genkit with the plugin
	g := genkit.Init(ctx,
		genkit.WithPlugins(googleAI),
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
	)

	// Load prompts
	basePrompt, err := loadPrompt("base")
	if err != nil {
		return nil, fmt.Errorf("failed to load base prompt: %w", err)
	}

	moderatePrompt, err := loadPrompt("moderate")
	if err != nil {
		return nil, fmt.Errorf("failed to load moderate prompt: %w", err)
	}

	analyzePrompt, err := loadPrompt("analyze")
	if err != nil {
		return nil, fmt.Errorf("failed to load analyze prompt: %w", err)
	}

	return &Moderator{
		genkit:         g,
		basePrompt:     basePrompt,
		moderatePrompt: moderatePrompt,
		analyzePrompt:  analyzePrompt,
	}, nil
}

// Moderate performs AI-powered content moderation
func (m *Moderator) Moderate(ctx context.Context, req models.ContentRequest) (*models.ModerationResult, error) {
	startTime := time.Now()

	// Build the complete prompt
	fullPrompt := fmt.Sprintf("%s\n\n%s\n\nContent to analyze:\n%s\n\nProvide your moderation analysis in JSON format with the following structure:\n{\n  \"safe\": true/false,\n  \"categories\": [\n    {\n      \"category\": \"category name\",\n      \"flagged\": true/false,\n      \"confidence\": 0.0-1.0,\n      \"severity\": \"low|medium|high\"\n    }\n  ],\n  \"overall_risk\": \"low|medium|high|critical\"\n}",
		m.basePrompt,
		m.moderatePrompt,
		req.Content,
	)

	// Call Gemini using Generate
	response, err := genkit.Generate(ctx, m.genkit,
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

	// Clean markdown code blocks
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

	// Parse JSON response
	var aiResponse struct {
		Safe        bool `json:"safe"`
		Categories  []struct {
			Category   string  `json:"category"`
			Flagged    bool    `json:"flagged"`
			Confidence float64 `json:"confidence"`
			Severity   string  `json:"severity"`
		} `json:"categories"`
		OverallRisk string `json:"overall_risk"`
	}

	if err := json.Unmarshal([]byte(responseText), &aiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w\nResponse: %s", err, responseText)
	}

	// Convert to response types
	var categories []models.ModerationCategory
	for _, cat := range aiResponse.Categories {
		categories = append(categories, models.ModerationCategory{
			Category:   cat.Category,
			Flagged:    cat.Flagged,
			Confidence: cat.Confidence,
			Severity:   cat.Severity,
		})
	}

	// Build response
	processingTime := time.Since(startTime).Milliseconds()

	return &models.ModerationResult{
		ContentID:        req.ContentID,
		Safe:             aiResponse.Safe,
		Categories:       categories,
		OverallRisk:      aiResponse.OverallRisk,
		ProcessingTimeMs: processingTime,
	}, nil
}

// AnalyzeComprehensive performs comprehensive analysis combining sentiment and moderation
func (m *Moderator) AnalyzeComprehensive(ctx context.Context, req models.AnalysisRequest) (*models.ComprehensiveAnalysis, error) {
	startTime := time.Now()

	// Build the complete prompt
	fullPrompt := fmt.Sprintf("%s\n\n%s\n\nContent to analyze:\n%s\n\nAuthor: %s\nContext URL: %s\n\nProvide comprehensive analysis in JSON format with the following structure:\n{\n  \"sentiment\": {\n    \"label\": \"positive|negative|neutral\",\n    \"confidence\": 0.0-1.0,\n    \"score\": -1.0 to 1.0\n  },\n  \"moderation\": {\n    \"safe\": true/false,\n    \"categories\": [\n      {\n        \"category\": \"category name\",\n        \"flagged\": true/false,\n        \"confidence\": 0.0-1.0,\n        \"severity\": \"low|medium|high\"\n      }\n    ],\n    \"overall_risk\": \"low|medium|high|critical\"\n  },\n  \"recommendation\": {\n    \"action\": \"approve|flag|reject|escalate\",\n    \"confidence\": 0.0-1.0,\n    \"reason\": \"explanation\",\n    \"auto_execute\": true/false\n  }\n}",
		m.basePrompt,
		m.analyzePrompt,
		req.Content,
		req.Author,
		req.ContextURL,
	)

	// Call Gemini using Generate
	response, err := genkit.Generate(ctx, m.genkit,
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

	// Clean markdown code blocks
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

	// Parse JSON response
	var aiResponse struct {
		Sentiment struct {
			Label      string  `json:"label"`
			Confidence float64 `json:"confidence"`
			Score      float64 `json:"score"`
		} `json:"sentiment"`
		Moderation struct {
			Safe        bool `json:"safe"`
			Categories  []struct {
				Category   string  `json:"category"`
				Flagged    bool    `json:"flagged"`
				Confidence float64 `json:"confidence"`
				Severity   string  `json:"severity"`
			} `json:"categories"`
			OverallRisk string `json:"overall_risk"`
		} `json:"moderation"`
		Recommendation struct {
			Action      string  `json:"action"`
			Confidence  float64 `json:"confidence"`
			Reason      string  `json:"reason"`
			AutoExecute bool    `json:"auto_execute"`
		} `json:"recommendation"`
	}

	if err := json.Unmarshal([]byte(responseText), &aiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w\nResponse: %s", err, responseText)
	}

	// Convert moderation categories
	var categories []models.ModerationCategory
	for _, cat := range aiResponse.Moderation.Categories {
		categories = append(categories, models.ModerationCategory{
			Category:   cat.Category,
			Flagged:    cat.Flagged,
			Confidence: cat.Confidence,
			Severity:   cat.Severity,
		})
	}

	// Apply auto-execute threshold if specified
	autoExecute := aiResponse.Recommendation.AutoExecute
	if req.AutoExecuteThreshold > 0 {
		autoExecute = aiResponse.Recommendation.Confidence >= req.AutoExecuteThreshold
	}

	// Build response
	processingTime := time.Since(startTime).Milliseconds()

	return &models.ComprehensiveAnalysis{
		ContentID: req.ContentID,
		Author:    req.Author,
		Sentiment: models.SentimentScore{
			Label:      aiResponse.Sentiment.Label,
			Confidence: aiResponse.Sentiment.Confidence,
			Score:      aiResponse.Sentiment.Score,
		},
		Moderation: models.ModerationResult{
			ContentID:        req.ContentID,
			Safe:             aiResponse.Moderation.Safe,
			Categories:       categories,
			OverallRisk:      aiResponse.Moderation.OverallRisk,
			ProcessingTimeMs: 0, // Not tracked separately
		},
		Recommendation: models.ActionRecommendation{
			Action:      aiResponse.Recommendation.Action,
			Confidence:  aiResponse.Recommendation.Confidence,
			Reason:      aiResponse.Recommendation.Reason,
			AutoExecute: autoExecute,
		},
		ProcessingTimeMs: processingTime,
		AnalyzedAt:       time.Now(),
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
