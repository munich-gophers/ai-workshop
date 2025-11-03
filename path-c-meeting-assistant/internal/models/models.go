package models

import "time"

// MeetingNotesRequest represents the input for meeting analysis
type MeetingNotesRequest struct {
	Notes string `json:"notes"` // Raw meeting notes text
}

// ActionItem represents a task extracted from meeting notes
type ActionItem struct {
	Description string    `json:"description"`           // What needs to be done
	Assignee    string    `json:"assignee,omitempty"`    // Who is responsible
	DueDate     string    `json:"due_date,omitempty"`    // When it's due (parsed from notes)
	Priority    string    `json:"priority,omitempty"`    // high, medium, low
	Status      string    `json:"status"`                // pending, in_progress, completed
	ExtractedAt time.Time `json:"extracted_at"`          // When the action was extracted
}

// Decision represents a key decision made in the meeting
type Decision struct {
	Description string   `json:"description"`         // What was decided
	Context     string   `json:"context,omitempty"`   // Why the decision was made
	Impact      string   `json:"impact,omitempty"`    // Expected impact
	Owners      []string `json:"owners,omitempty"`    // Who is responsible for the decision
	MadeAt      string   `json:"made_at,omitempty"`   // When the decision was made (parsed from notes)
}

// Participant represents a meeting attendee
type Participant struct {
	Name       string   `json:"name"`                  // Person's name
	Role       string   `json:"role,omitempty"`        // Their role in the meeting
	Mentions   int      `json:"mentions"`              // How many times they were mentioned
	ActionItems []string `json:"action_items,omitempty"` // Action items assigned to them
}

// Topic represents a discussion topic
type Topic struct {
	Title       string   `json:"title"`                 // Topic title
	Summary     string   `json:"summary,omitempty"`     // Brief summary
	Duration    string   `json:"duration,omitempty"`    // How long it was discussed
	KeyPoints   []string `json:"key_points,omitempty"`  // Main points discussed
}

// ExtractResponse contains extracted action items (Checkpoint 1)
type ExtractResponse struct {
	ActionItems      []ActionItem `json:"action_items"`       // List of action items found
	TotalActions     int          `json:"total_actions"`      // Count of actions
	ProcessingTimeMs int64        `json:"processing_time_ms"` // How long extraction took
	Summary          string       `json:"summary,omitempty"`  // Brief summary of the meeting
}

// AnalyzeResponse contains full meeting analysis (Checkpoint 2)
type AnalyzeResponse struct {
	ActionItems      []ActionItem  `json:"action_items"`       // Action items extracted
	Decisions        []Decision    `json:"decisions"`          // Key decisions made
	Participants     []Participant `json:"participants"`       // Meeting participants
	Topics           []Topic       `json:"topics,omitempty"`   // Topics discussed
	MeetingDate      string        `json:"meeting_date,omitempty"` // Meeting date (parsed)
	NextMeeting      string        `json:"next_meeting,omitempty"` // Next meeting info
	Summary          string        `json:"summary"`            // Overall meeting summary
	ProcessingTimeMs int64         `json:"processing_time_ms"` // Processing time
}

// EmailRequest requests follow-up email generation (Checkpoint 3)
type EmailRequest struct {
	Notes            string       `json:"notes"`                       // Original meeting notes
	ActionItems      []ActionItem `json:"action_items,omitempty"`      // Pre-extracted action items
	Decisions        []Decision   `json:"decisions,omitempty"`         // Pre-analyzed decisions
	Participants     []Participant `json:"participants,omitempty"`     // Pre-analyzed participants
	Tone             string       `json:"tone,omitempty"`              // formal, casual, friendly
	IncludeSummary   bool         `json:"include_summary"`             // Include meeting summary
	IncludeNextSteps bool         `json:"include_next_steps"`          // Include next steps section
	RecipientName    string       `json:"recipient_name,omitempty"`    // Recipient's name for personalization
}

// EmailResponse contains generated follow-up email (Checkpoint 3)
type EmailResponse struct {
	Subject          string `json:"subject"`            // Email subject line
	Body             string `json:"body"`               // Email body content
	ProcessingTimeMs int64  `json:"processing_time_ms"` // Processing time
}
