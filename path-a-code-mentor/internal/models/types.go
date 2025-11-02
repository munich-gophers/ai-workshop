package models

// ReviewRequest represents a code review request
type ReviewRequest struct {
	Diff        string `json:"diff"`                   // Git diff content
	Language    string `json:"language"`               // Programming language (optional, auto-detected)
	FilePath    string `json:"file_path"`              // File being reviewed
	PRNumber    int    `json:"pr_number,omitempty"`    // Pull request number (for webhook)
	Title       string `json:"title,omitempty"`        // PR title (for webhook)
	Description string `json:"description,omitempty"`  // PR description (for webhook)
	Author      string `json:"author,omitempty"`       // PR author (for webhook)
	RepoName    string `json:"repo_name,omitempty"`    // Repository name (for webhook)
	Branch      string `json:"branch,omitempty"`       // Branch name (for webhook)
}

// ReviewResponse represents the AI's code review
type ReviewResponse struct {
	Suggestions      []Suggestion `json:"suggestions"`
	Summary          string       `json:"summary"`
	Severity         string       `json:"severity"` // Overall severity: low, medium, high
	Language         string       `json:"language"`
	FilePath         string       `json:"file_path"`
	ProcessingTimeMs int          `json:"processing_time_ms"`
	SecretsDetected  bool         `json:"secrets_detected,omitempty"`
	SecretLocations  []Secret     `json:"secret_locations,omitempty"`
}

// Suggestion represents a single code review comment
type Suggestion struct {
	Line        int    `json:"line,omitempty"`       // Line number
	File        string `json:"file,omitempty"`       // File path
	Severity    string `json:"severity"`             // low, medium, high, critical
	Category    string `json:"category"`             // bug, performance, style, security, best-practice
	Message     string `json:"message"`              // The feedback
	Title       string `json:"title,omitempty"`      // Short title for the suggestion
	Description string `json:"description"`          // Why this matters (alias for Explanation)
	Explanation string `json:"explanation"`          // Why this matters
	Suggestion  string `json:"suggestion,omitempty"` // Suggested fix (optional)
	CodeExample string `json:"code_example,omitempty"` // Code example for the fix
}

// Review represents a complete code review (for webhook responses)
type Review struct {
	Summary          string       `json:"summary"`
	Suggestions      []Suggestion `json:"suggestions"`
	PositiveFeedback []string     `json:"positive_feedback,omitempty"`
	Approved         bool         `json:"approved"`
}

// ToReview converts a ReviewResponse to a Review for webhook posting
func (r *ReviewResponse) ToReview() *Review {
	// Determine if the PR should be approved based on severity
	approved := r.Severity != SeverityHigh && len(r.Suggestions) < 5

	// Extract positive feedback from low-severity suggestions
	var positiveFeedback []string
	for _, suggestion := range r.Suggestions {
		if suggestion.Severity == SeverityLow && suggestion.Category == CategoryStyle {
			positiveFeedback = append(positiveFeedback, suggestion.Message)
		}
	}

	return &Review{
		Summary:          r.Summary,
		Suggestions:      r.Suggestions,
		PositiveFeedback: positiveFeedback,
		Approved:         approved,
	}
}

// Secret represents a detected secret in code
type Secret struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Type    string `json:"type"`    // api_key, token, password, etc.
	Message string `json:"message"` // Description of the issue
}

// Severity levels for suggestions (matches AI schema)
const (
	SeverityLow    = "low"    // Minor issues, nitpicks, personal preference
	SeverityMedium = "medium" // Performance issues, deprecated patterns, best practices
	SeverityHigh   = "high"   // Security issues, bugs, critical problems
)

// Category types for suggestions
const (
	CategoryBug          = "bug"           // Code defects and bugs
	CategoryPerformance  = "performance"   // Performance issues
	CategoryStyle        = "style"         // Code style and formatting
	CategorySecurity     = "security"      // Security vulnerabilities
	CategoryBestPractice = "best-practice" // Best practices and patterns
)

// Secret types that can be detected
const (
	SecretAPIKey           = "api_key"
	SecretToken            = "token"
	SecretPassword         = "password"
	SecretPrivateKey       = "private_key"
	SecretConnectionString = "connection_string"
)
