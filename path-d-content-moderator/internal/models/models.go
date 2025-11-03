package models

import "time"

// ContentRequest represents incoming content to analyze
type ContentRequest struct {
	Content    string `json:"content"`
	ContentID  string `json:"content_id,omitempty"`
	Author     string `json:"author,omitempty"`
	ContextURL string `json:"context_url,omitempty"`
}

// SentimentScore represents sentiment analysis results
type SentimentScore struct {
	Label      string  `json:"label"`       // positive, negative, neutral
	Confidence float64 `json:"confidence"`  // 0.0 to 1.0
	Score      float64 `json:"score"`       // -1.0 (negative) to 1.0 (positive)
}

// SentimentResponse represents sentiment analysis output
type SentimentResponse struct {
	ContentID        string         `json:"content_id,omitempty"`
	Sentiment        SentimentScore `json:"sentiment"`
	ProcessingTimeMs int64          `json:"processing_time_ms"`
	Method           string         `json:"method"` // pattern-based or ai-powered
}

// ModerationCategory represents a specific moderation issue
type ModerationCategory struct {
	Category   string  `json:"category"`   // spam, harassment, hate-speech, inappropriate, etc.
	Flagged    bool    `json:"flagged"`    // whether this category is present
	Confidence float64 `json:"confidence"` // 0.0 to 1.0
	Severity   string  `json:"severity"`   // low, medium, high
}

// ModerationResult represents content moderation analysis
type ModerationResult struct {
	ContentID        string               `json:"content_id,omitempty"`
	Safe             bool                 `json:"safe"`        // overall safety assessment
	Categories       []ModerationCategory `json:"categories"`  // specific issues detected
	OverallRisk      string               `json:"overall_risk"` // low, medium, high, critical
	ProcessingTimeMs int64                `json:"processing_time_ms"`
}

// ActionRecommendation represents an automated action to take
type ActionRecommendation struct {
	Action      string  `json:"action"`      // approve, flag, reject, escalate
	Confidence  float64 `json:"confidence"`  // 0.0 to 1.0
	Reason      string  `json:"reason"`      // explanation
	AutoExecute bool    `json:"auto_execute"` // whether to automatically execute
}

// ComprehensiveAnalysis combines sentiment, moderation, and recommendations
type ComprehensiveAnalysis struct {
	ContentID        string                 `json:"content_id,omitempty"`
	Author           string                 `json:"author,omitempty"`
	Sentiment        SentimentScore         `json:"sentiment"`
	Moderation       ModerationResult       `json:"moderation"`
	Recommendation   ActionRecommendation   `json:"recommendation"`
	ProcessingTimeMs int64                  `json:"processing_time_ms"`
	AnalyzedAt       time.Time              `json:"analyzed_at"`
}

// AnalysisRequest for comprehensive analysis endpoint
type AnalysisRequest struct {
	Content           string  `json:"content"`
	ContentID         string  `json:"content_id,omitempty"`
	Author            string  `json:"author,omitempty"`
	ContextURL        string  `json:"context_url,omitempty"`
	AutoExecuteThreshold float64 `json:"auto_execute_threshold,omitempty"` // threshold for auto-execution (0.0-1.0)
}
