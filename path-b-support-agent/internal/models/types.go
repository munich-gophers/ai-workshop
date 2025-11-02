package models

// TriageRequest represents an incoming customer support message
type TriageRequest struct {
	Message string `json:"message"`                // Customer message
	Channel string `json:"channel,omitempty"`      // email, chat, phone, etc.
	UserID  string `json:"user_id,omitempty"`      // Optional user identifier
	Ticket  string `json:"ticket_id,omitempty"`    // Optional ticket reference
}

// TriageResponse represents the AI-processed triage result
type TriageResponse struct {
	OriginalMessage  string         `json:"original_message"`           // Original customer message
	RedactedMessage  string         `json:"redacted_message"`           // Message with PII redacted
	Intent           Intent         `json:"intent"`                     // Classified intent
	Urgency          Urgency        `json:"urgency"`                    // Detected urgency level
	Summary          string         `json:"summary"`                    // AI-generated summary
	DetectedPII      []PIIDetection `json:"detected_pii"`              // List of PII found
	SuggestedRouting string         `json:"suggested_routing"`          // Routing recommendation
	ProcessingTimeMs int64          `json:"processing_time_ms"`        // Processing duration
}

// Intent represents the classified customer intent
type Intent struct {
	Category   string  `json:"category"`    // billing, technical, account, general
	Confidence float64 `json:"confidence"`  // Confidence score 0.0-1.0
	Subcategory string `json:"subcategory,omitempty"` // More specific categorization
}

// Urgency represents the priority level of the request
type Urgency struct {
	Level      string  `json:"level"`       // low, medium, high, critical
	Confidence float64 `json:"confidence"`  // Confidence score 0.0-1.0
	Reason     string  `json:"reason,omitempty"` // Explanation for urgency level
}

// PIIDetection represents a detected PII element
type PIIDetection struct {
	Type     string `json:"type"`      // email, phone, ssn, credit_card, etc.
	Location string `json:"location"`  // Where in message
	Pattern  string `json:"pattern"`   // Which pattern matched
}

// PIIPattern represents a pattern for detecting PII
type PIIPattern struct {
	Type    string // Type of PII (email, phone, ssn, etc.)
	Pattern string // Regex pattern
	Message string // Human-readable description
}
