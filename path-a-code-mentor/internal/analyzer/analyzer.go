package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/munich-gophers/ai-workshop/code-mentor/internal/models"
)

// Analyzer performs AI-powered code analysis using Genkit v1.1.0 API
type Analyzer struct {
	genkit       *genkit.Genkit
	model        ai.Model
	basePrompt   string
	goPrompt     string
	pythonPrompt string
	jsPrompt     string
}

// New creates a new Analyzer instance using Genkit v1.1.0 API
func New(ctx context.Context) (*Analyzer, error) {
	// Get API key from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// Initialize Google AI plugin
	googleAI := &googlegenai.GoogleAI{
		APIKey: apiKey,
	}

	// Initialize Genkit with the Google AI plugin
	g := genkit.Init(ctx, genkit.WithPlugins(googleAI))

	// Get the model
	model := googlegenai.GoogleAIModel(g, "gemini-2.5-flash")

	// Load prompts from files
	basePrompt, err := loadPrompt("internal/prompts/review-base.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load base prompt: %w", err)
	}

	goPrompt, err := loadPrompt("internal/prompts/review-go.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load Go prompt: %w", err)
	}

	pythonPrompt, err := loadPrompt("internal/prompts/review-python.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load Python prompt: %w", err)
	}

	jsPrompt, err := loadPrompt("internal/prompts/review-javascript.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load JavaScript prompt: %w", err)
	}

	return &Analyzer{
		genkit:       g,
		model:        model,
		basePrompt:   basePrompt,
		goPrompt:     goPrompt,
		pythonPrompt: pythonPrompt,
		jsPrompt:     jsPrompt,
	}, nil
}

// Review analyzes code and returns suggestions using Genkit v1.1.0 API
func (a *Analyzer) Review(ctx context.Context, req models.ReviewRequest) (*models.ReviewResponse, error) {
	// Build prompt based on language
	prompt := a.buildPrompt(req)

	// Build the model request with structured output
	modelReq := &ai.ModelRequest{
		Messages: []*ai.Message{
			ai.NewUserMessage(ai.NewTextPart(prompt)),
		},
		Output: buildOutputSchema(),
	}

	// Generate response using the model directly
	resp, err := a.model.Generate(ctx, modelReq, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate review: %w", err)
	}

	// Extract text from response
	text, err := extractResponseText(resp)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	response, err := a.parseResponse(text, req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse review response: %w", err)
	}

	return response, nil
}

// ReviewSimple is an alternative without structured output (more flexible, less strict)
func (a *Analyzer) ReviewSimple(ctx context.Context, req models.ReviewRequest) (*models.ReviewResponse, error) {
	// Build prompt based on language
	prompt := a.buildPrompt(req)

	// Build the model request without output schema
	modelReq := &ai.ModelRequest{
		Messages: []*ai.Message{
			ai.NewUserMessage(ai.NewTextPart(prompt)),
		},
	}

	// Generate response using the model directly
	resp, err := a.model.Generate(ctx, modelReq, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate review: %w", err)
	}

	// Extract text from response
	text, err := extractResponseText(resp)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response with more lenient error handling
	response, err := a.parseResponse(text, req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse review response: %w", err)
	}

	return response, nil
}

// buildOutputSchema creates the JSON schema for structured output
func buildOutputSchema() *ai.ModelOutputConfig {
	return &ai.ModelOutputConfig{
		Format: ai.OutputFormatJSON,
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"suggestions": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"line": map[string]any{
								"type": "integer",
							},
							"message": map[string]any{
								"type": "string",
							},
							"severity": map[string]any{
								"type": "string",
								"enum": []string{"low", "medium", "high"},
							},
							"category": map[string]any{
								"type": "string",
								"enum": []string{"bug", "performance", "style", "security", "best-practice"},
							},
							"explanation": map[string]any{
								"type": "string",
							},
						},
						"required": []string{"line", "message", "severity", "category", "explanation"},
					},
				},
				"summary": map[string]any{
					"type": "string",
				},
				"severity": map[string]any{
					"type": "string",
					"enum": []string{"low", "medium", "high"},
				},
			},
			"required": []string{"suggestions", "summary", "severity"},
		},
	}
}

// extractResponseText safely extracts text from AI response
func extractResponseText(resp *ai.ModelResponse) (string, error) {
	if resp.Message == nil {
		return "", fmt.Errorf("no message in response")
	}

	if len(resp.Message.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return resp.Message.Content[0].Text, nil
}

// buildPrompt constructs the full prompt for the AI
func (a *Analyzer) buildPrompt(req models.ReviewRequest) string {
	var sb strings.Builder

	// Start with base prompt
	sb.WriteString(a.basePrompt)
	sb.WriteString("\n\n")

	// Add language-specific guidance
	switch req.Language {
	case "go":
		sb.WriteString(a.goPrompt)
		sb.WriteString("\n\n")
	case "python":
		sb.WriteString(a.pythonPrompt)
		sb.WriteString("\n\n")
	case "javascript", "typescript":
		sb.WriteString(a.jsPrompt)
		sb.WriteString("\n\n")
	}

	// Add the code diff and instructions
	sb.WriteString(fmt.Sprintf(`File: %s
Language: %s

Code Diff:
%s

Please analyze this code and provide a review in JSON format with the following structure:
{
  "suggestions": [
    {
      "line": <line_number>,
      "message": "<suggestion_text>",
      "severity": "low|medium|high",
      "category": "bug|performance|style|security|best-practice",
      "explanation": "<why_this_matters>"
    }
  ],
  "summary": "<brief_overall_assessment>",
  "severity": "low|medium|high"
}

Return ONLY valid JSON without markdown code blocks.`,
		req.FilePath, req.Language, req.Diff))

	return sb.String()
}

// parseResponse handles parsing the AI response into structured data
func (a *Analyzer) parseResponse(text string, req models.ReviewRequest) (*models.ReviewResponse, error) {
	// Clean up response (remove markdown code blocks if present)
	cleanText := cleanJSONResponse(text)

	var response models.ReviewResponse
	if err := json.Unmarshal([]byte(cleanText), &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w (text: %s)", err, cleanText)
	}

	// Add metadata
	response.Language = req.Language
	response.FilePath = req.FilePath

	return &response, nil
}

// cleanJSONResponse removes markdown code blocks from response text
func cleanJSONResponse(text string) string {
	cleanText := strings.TrimSpace(text)

	// Remove ```json``` blocks
	if strings.HasPrefix(cleanText, "```json") {
		cleanText = strings.TrimPrefix(cleanText, "```json")
		cleanText = strings.TrimSuffix(cleanText, "```")
		cleanText = strings.TrimSpace(cleanText)
	} else if strings.HasPrefix(cleanText, "```") {
		// Remove generic ``` blocks
		cleanText = strings.TrimPrefix(cleanText, "```")
		cleanText = strings.TrimSuffix(cleanText, "```")
		cleanText = strings.TrimSpace(cleanText)
	}

	return cleanText
}

// loadPrompt loads a prompt template from file
func loadPrompt(filename string) (string, error) {
	// Try current directory first
	content, err := os.ReadFile(filename)
	if err != nil {
		// Try relative to project root
		content, err = os.ReadFile(filepath.Join("..", "..", filename))
		if err != nil {
			return "", fmt.Errorf("failed to read prompt file %s: %w", filename, err)
		}
	}
	return string(content), nil
}
