package security

import (
	"strings"

	"github.com/munich-gophers/ai-workshop/code-mentor/internal/models"
)

// ScanCode scans code for secrets and returns detected issues
func ScanCode(code string, filePath string) []models.Secret {
	// TODO: CHECKPOINT 3 - Implement secret detection
	// Hint:
	// 1. Get patterns with GetSecretPatterns()
	// 2. Split code into lines
	// 3. For each line, check each pattern
	// 4. If match found, append to secrets slice with line number
	// 5. Return all detected secrets

	secrets := []models.Secret{}
	// Your implementation here

	return secrets
}

// RedactSecrets replaces secrets in code with placeholders
func RedactSecrets(code string) string {
	// TODO: CHECKPOINT 3 - Implement secret redaction
	// Hint:
	// 1. Get patterns with GetSecretPatterns()
	// 2. For each pattern, replace matches with [TYPE_REDACTED]
	// 3. Example: "sk-abc123" → "[API_KEY_REDACTED]"

	redacted := code
	patterns := GetSecretPatterns()

	for _, pattern := range patterns {
		placeholder := "[" + strings.ToUpper(pattern.Type) + "_REDACTED]"
		redacted = pattern.Pattern.ReplaceAllString(redacted, placeholder)
	}

	return redacted
}
