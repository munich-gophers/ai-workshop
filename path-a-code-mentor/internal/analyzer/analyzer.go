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
	// - model *genkit.Model
	// - basePrompt string
	// - languagePrompts map[string]string
}

// New creates a new code analyzer
func New(ctx context.Context) (*Analyzer, error) {
	// TODO: CHECKPOINT 2 - Initialize Genkit
	//
	// Steps:
	// 1. Import "github.com/firebase/genkit/go/genkit"
	// 2. Import "github.com/firebase/genkit/go/plugins/googleai"
	// 3. Initialize: err := googleai.Init(ctx, nil)
	// 4. Create model: model := googleai.Model("gemini-1.5-flash")
	// 5. Load prompts: basePrompt, langPrompts, err := loadPrompts()
	// 6. Return &Analyzer{model: model, basePrompt: basePrompt, languagePrompts: langPrompts}
	//
	// Example:
	//   if err := googleai.Init(ctx, nil); err != nil {
	//       return nil, fmt.Errorf("genkit init: %w", err)
	//   }
	//
	//   model := googleai.Model("gemini-1.5-flash")
	//   basePrompt, languagePrompts, err := loadPrompts()
	//   if err != nil {
	//       return nil, err
	//   }
	//
	//   return &Analyzer{
	//       model: model,
	//       basePrompt: basePrompt,
	//       languagePrompts: languagePrompts,
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
