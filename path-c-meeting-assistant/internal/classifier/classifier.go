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
	"github.com/munich-gophers/ai-workshop/meeting-assistant/internal/models"
)

// Classifier performs AI-powered meeting notes analysis using Genkit v1.1.0 API
type Classifier struct {
	genkit        *genkit.Genkit
	basePrompt    string
	extractPrompt string
	analyzePrompt string
	emailPrompt   string
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

	extractPrompt, err := loadPrompt("extract")
	if err != nil {
		return nil, fmt.Errorf("failed to load extract prompt: %w", err)
	}

	analyzePrompt, err := loadPrompt("analyze")
	if err != nil {
		return nil, fmt.Errorf("failed to load analyze prompt: %w", err)
	}

	emailPrompt, err := loadPrompt("email")
	if err != nil {
		return nil, fmt.Errorf("failed to load email prompt: %w", err)
	}

	return &Classifier{
		genkit:        g,
		basePrompt:    basePrompt,
		extractPrompt: extractPrompt,
		analyzePrompt: analyzePrompt,
		emailPrompt:   emailPrompt,
	}, nil
}

// Analyze performs comprehensive meeting analysis
func (c *Classifier) Analyze(ctx context.Context, req models.MeetingNotesRequest) (*models.AnalyzeResponse, error) {
	startTime := time.Now()

	// Build the complete prompt
	fullPrompt := fmt.Sprintf("%s\n\n%s\n\nMeeting Notes:\n%s\n\nProvide your analysis in JSON format with the following structure:\n{\n  \"action_items\": [\n    {\n      \"description\": \"task description\",\n      \"assignee\": \"person name\",\n      \"due_date\": \"date or deadline\",\n      \"priority\": \"high|medium|low\",\n      \"status\": \"pending\"\n    }\n  ],\n  \"decisions\": [\n    {\n      \"description\": \"what was decided\",\n      \"context\": \"why it was decided\",\n      \"impact\": \"expected impact\",\n      \"owners\": [\"person1\", \"person2\"]\n    }\n  ],\n  \"participants\": [\n    {\n      \"name\": \"person name\",\n      \"role\": \"their role or contribution\",\n      \"mentions\": 1\n    }\n  ],\n  \"topics\": [\n    {\n      \"title\": \"topic title\",\n      \"summary\": \"brief summary\",\n      \"key_points\": [\"point1\", \"point2\"]\n    }\n  ],\n  \"meeting_date\": \"date if mentioned\",\n  \"next_meeting\": \"next meeting info if mentioned\",\n  \"summary\": \"overall meeting summary\"\n}",
		c.basePrompt,
		c.analyzePrompt,
		req.Notes,
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
		ActionItems []struct {
			Description string `json:"description"`
			Assignee    string `json:"assignee"`
			DueDate     string `json:"due_date"`
			Priority    string `json:"priority"`
			Status      string `json:"status"`
		} `json:"action_items"`
		Decisions []struct {
			Description string   `json:"description"`
			Context     string   `json:"context"`
			Impact      string   `json:"impact"`
			Owners      []string `json:"owners"`
		} `json:"decisions"`
		Participants []struct {
			Name     string `json:"name"`
			Role     string `json:"role"`
			Mentions int    `json:"mentions"`
		} `json:"participants"`
		Topics []struct {
			Title     string   `json:"title"`
			Summary   string   `json:"summary"`
			KeyPoints []string `json:"key_points"`
		} `json:"topics"`
		MeetingDate string `json:"meeting_date"`
		NextMeeting string `json:"next_meeting"`
		Summary     string `json:"summary"`
	}

	if err := json.Unmarshal([]byte(responseText), &aiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w\nResponse: %s", err, responseText)
	}

	// Convert to response types
	var actionItems []models.ActionItem
	for _, item := range aiResponse.ActionItems {
		actionItems = append(actionItems, models.ActionItem{
			Description: item.Description,
			Assignee:    item.Assignee,
			DueDate:     item.DueDate,
			Priority:    item.Priority,
			Status:      item.Status,
			ExtractedAt: time.Now(),
		})
	}

	var decisions []models.Decision
	for _, dec := range aiResponse.Decisions {
		decisions = append(decisions, models.Decision{
			Description: dec.Description,
			Context:     dec.Context,
			Impact:      dec.Impact,
			Owners:      dec.Owners,
		})
	}

	var participants []models.Participant
	for _, part := range aiResponse.Participants {
		participants = append(participants, models.Participant{
			Name:     part.Name,
			Role:     part.Role,
			Mentions: part.Mentions,
		})
	}

	var topics []models.Topic
	for _, topic := range aiResponse.Topics {
		topics = append(topics, models.Topic{
			Title:     topic.Title,
			Summary:   topic.Summary,
			KeyPoints: topic.KeyPoints,
		})
	}

	// Build response
	processingTime := time.Since(startTime).Milliseconds()

	return &models.AnalyzeResponse{
		ActionItems:      actionItems,
		Decisions:        decisions,
		Participants:     participants,
		Topics:           topics,
		MeetingDate:      aiResponse.MeetingDate,
		NextMeeting:      aiResponse.NextMeeting,
		Summary:          aiResponse.Summary,
		ProcessingTimeMs: processingTime,
	}, nil
}

// GenerateEmail generates a follow-up email from meeting notes
func (c *Classifier) GenerateEmail(ctx context.Context, req models.EmailRequest) (*models.EmailResponse, error) {
	startTime := time.Now()

	// Build context from provided data or notes
	var actionItemsText string
	if len(req.ActionItems) > 0 {
		actionItemsText = "\n\nAction Items:\n"
		for _, item := range req.ActionItems {
			actionItemsText += fmt.Sprintf("- %s", item.Description)
			if item.Assignee != "" {
				actionItemsText += fmt.Sprintf(" (@%s)", item.Assignee)
			}
			if item.DueDate != "" {
				actionItemsText += fmt.Sprintf(" - Due: %s", item.DueDate)
			}
			actionItemsText += "\n"
		}
	}

	var decisionsText string
	if len(req.Decisions) > 0 {
		decisionsText = "\n\nKey Decisions:\n"
		for _, dec := range req.Decisions {
			decisionsText += fmt.Sprintf("- %s\n", dec.Description)
		}
	}

	// Determine tone instructions
	toneInstructions := "Use a professional and friendly tone."
	if req.Tone == "formal" {
		toneInstructions = "Use a formal, professional tone with full names and structured language."
	} else if req.Tone == "casual" {
		toneInstructions = "Use a casual but professional tone with first names and conversational language."
	} else if req.Tone == "friendly" {
		toneInstructions = "Use a warm, friendly, and collaborative tone."
	}

	recipientContext := ""
	if req.RecipientName != "" {
		recipientContext = fmt.Sprintf("\n\nRecipient: %s", req.RecipientName)
	}

	// Build the complete prompt
	fullPrompt := fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\nMeeting Notes:\n%s%s%s\n\nGenerate a professional follow-up email in JSON format:\n{\n  \"subject\": \"clear and specific subject line\",\n  \"body\": \"complete email body with proper formatting\"\n}",
		c.basePrompt,
		c.emailPrompt,
		toneInstructions,
		recipientContext,
		req.Notes,
		actionItemsText,
		decisionsText,
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
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	if err := json.Unmarshal([]byte(responseText), &aiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w\nResponse: %s", err, responseText)
	}

	// Build response
	processingTime := time.Since(startTime).Milliseconds()

	return &models.EmailResponse{
		Subject:          aiResponse.Subject,
		Body:             aiResponse.Body,
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
