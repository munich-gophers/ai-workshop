package redactor

import (
	"regexp"
	"strings"

	"github.com/munich-gophers/ai-workshop/support-agent/internal/models"
)

// GetPIIPatterns returns all PII detection patterns
func GetPIIPatterns() []models.PIIPattern {
	return []models.PIIPattern{
		{
			Type:    "email",
			Pattern: `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`,
			Message: "Email address detected",
		},
		{
			Type:    "phone",
			Pattern: `(\+?1[-.]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}`,
			Message: "Phone number detected",
		},
		{
			Type:    "ssn",
			Pattern: `\d{3}-\d{2}-\d{4}`,
			Message: "Social Security Number detected",
		},
		{
			Type:    "credit_card",
			Pattern: `\b(?:\d{4}[-\s]?){3}\d{4}\b`,
			Message: "Credit card number detected",
		},
		{
			Type:    "zipcode",
			Pattern: `\b\d{5}(?:-\d{4})?\b`,
			Message: "ZIP code detected",
		},
		{
			Type:    "ip_address",
			Pattern: `\b(?:\d{1,3}\.){3}\d{1,3}\b`,
			Message: "IP address detected",
		},
	}
}

// DetectPII scans text for PII and returns detected instances
func DetectPII(text string) []models.PIIDetection {
	var detections []models.PIIDetection
	patterns := GetPIIPatterns()

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.Pattern)
		matches := re.FindAllString(text, -1)

		for _, match := range matches {
			detections = append(detections, models.PIIDetection{
				Type:     pattern.Type,
				Location: match,
				Pattern:  pattern.Message,
			})
		}
	}

	return detections
}

// RedactPII replaces PII in text with placeholders
func RedactPII(text string) string {
	redacted := text
	patterns := GetPIIPatterns()

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.Pattern)
		placeholder := "[" + strings.ToUpper(pattern.Type) + "_REDACTED]"
		redacted = re.ReplaceAllString(redacted, placeholder)
	}

	return redacted
}
