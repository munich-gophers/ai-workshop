package models

// ReviewRequest represents a code review request
type ReviewRequest struct {
	Diff     string `json:"diff"`      // Git diff content
	Language string `json:"language"`  // Programming language (optional, auto-detected)
	FilePath string `json:"file_path"` // File being reviewed
}

// ReviewResponse represents the AI's code review
type ReviewResponse struct {
	Suggestions      []Suggestion `json:"suggestions"`
	Summary          string       `json:"summary"`
	Language         string       `json:"language"`
	ProcessingTimeMs int          `json:"processing_time_ms"`
	SecretsDetected  bool         `json:"secrets_detected,omitempty"`
	SecretLocations  []Secret     `json:"secret_locations,omitempty"`
}

// Suggestion represents a single code review comment
type Suggestion struct {
	Line       int    `json:"line,omitempty"`       // Line number
	File       string `json:"file,omitempty"`       // File path
	Severity   string `json:"severity"`             // critical, warning, info, nitpick
	Message    string `json:"message"`              // The feedback
	Suggestion string `json:"suggestion,omitempty"` // Suggested fix (optional)
}

// Secret represents a detected secret in code
type Secret struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Type    string `json:"type"`    // api_key, token, password, etc.
	Message string `json:"message"` // Description of the issue
}

// Severity levels for suggestions
const (
	SeverityCritical = "critical" // Security issues, bugs that will crash
	SeverityWarning  = "warning"  // Performance issues, deprecated patterns
	SeverityInfo     = "info"     // Best practices, documentation
	SeverityNitpick  = "nitpick"  // Minor formatting, personal preference
)

// Secret types that can be detected
const (
	SecretAPIKey           = "api_key"
	SecretToken            = "token"
	SecretPassword         = "password"
	SecretPrivateKey       = "private_key"
	SecretConnectionString = "connection_string"
)
