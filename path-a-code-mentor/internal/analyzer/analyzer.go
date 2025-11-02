package analyzer

import (
	"context"
	"fmt"
	"os"

	"github.com/munich-gophers/ai-workshop/code-mentor/internal/models"
)

type Analyzer struct {
	// TODO: CHECKPOINT 2 - Add fields
	//
	// You'll need:
	// - genkit *genkit.Genkit
	// - model ai.Model
	// - basePrompt string
	// - goPrompt string
	// - pythonPrompt string
	// - jsPrompt string
}

// New creates a new code analyzer
func New(ctx context.Context) (*Analyzer, error) {
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
	// 9. Return &Analyzer with all fields
	//
	// Example:
	//   apiKey := os.Getenv("GEMINI_API_KEY")
	//   if apiKey == "" {
	//       return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	//   }
	//
	//   googleAI := &googlegenai.GoogleAI{APIKey: apiKey}
	//   g := genkit.Init(ctx, genkit.WithPlugins(googleAI))
	//   model := googlegenai.GoogleAIModel(g, "gemini-2.5-flash")
	//
	//   basePrompt, err := loadPrompt("internal/prompts/review-base.txt")
	//   if err != nil {
	//       return nil, fmt.Errorf("failed to load base prompt: %w", err)
	//   }
	//
	//   goPrompt, err := loadPrompt("internal/prompts/review-go.txt")
	//   // ... load other prompts
	//
	//   return &Analyzer{
	//       genkit:       g,
	//       model:        model,
	//       basePrompt:   basePrompt,
	//       goPrompt:     goPrompt,
	//       pythonPrompt: pythonPrompt,
	//       jsPrompt:     jsPrompt,
	//   }, nil

	return nil, fmt.Errorf("not implemented - see TODO comments in analyzer.go")
}

// Review analyzes a code change and returns suggestions
func (a *Analyzer) Review(ctx context.Context, req models.ReviewRequest) (*models.ReviewResponse, error) {
	// TODO: CHECKPOINT 2 - Implement code review logic
	//
	// Steps:
	// 1. Detect language if not provided:
	//    if req.Language == "" {
	//        req.Language = detectLanguage(req.FilePath)
	//    }
	//
	// 2. Build prompt:
	//    prompt := a.buildPrompt(req)
	//
	// 3. Call Gemini:
	//    response, err := a.model.Generate(ctx,
	//        genkit.NewUserTextMessage(prompt),
	//        nil,
	//    )
	//
	// 4. Parse JSON response:
	//    var result models.ReviewResponse
	//    err = json.Unmarshal([]byte(response.Text), &result)
	//
	// 5. Set language and return:
	//    result.Language = req.Language
	//    return &result, nil

	return nil, fmt.Errorf("not implemented - see TODO comments in analyzer.go")
}

// buildPrompt creates the full prompt for Gemini
func (a *Analyzer) buildPrompt(req models.ReviewRequest) string {
	// TODO: CHECKPOINT 2 - Build the prompt
	//
	// Combine:
	// 1. a.basePrompt (general instructions)
	// 2. a.languagePrompts[req.Language] (if available)
	// 3. The code diff with file info
	//
	// Example format:
	//   prompt := a.basePrompt + "\n\n"
	//   if langPrompt, ok := a.languagePrompts[req.Language]; ok {
	//       prompt += langPrompt + "\n\n"
	//   }
	//   prompt += fmt.Sprintf("FILE: %s\nLANGUAGE: %s\n\nCODE DIFF:\n%s\n\n",
	//       req.FilePath, req.Language, req.Diff)
	//   prompt += "Respond with JSON only: {\"suggestions\":[...],\"summary\":\"...\"}"
	//   return prompt

	return ""
}

// loadPrompts reads prompt files from disk
func loadPrompts() (string, map[string]string, error) {
	// Load base prompt
	basePromptBytes, err := os.ReadFile("prompts/review-base.txt")
	if err != nil {
		return "", nil, fmt.Errorf("load base prompt: %w", err)
	}
	basePrompt := string(basePromptBytes)

	// TODO: CHECKPOINT 2 - Load language-specific prompts
	//
	// Hint:
	//   languagePrompts := make(map[string]string)
	//   languages := []string{"go", "python", "javascript", "typescript"}
	//   for _, lang := range languages {
	//       filename := fmt.Sprintf("prompts/review-%s.txt", lang)
	//       content, err := os.ReadFile(filename)
	//       if err != nil {
	//           continue // Optional prompts
	//       }
	//       languagePrompts[lang] = string(content)
	//   }

	languagePrompts := make(map[string]string)
	// Your implementation here

	return basePrompt, languagePrompts, nil
}
