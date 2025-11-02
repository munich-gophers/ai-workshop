package security

import (
	"strings"
)

// Finding represents a detected secret with additional metadata
type Finding struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Type     string `json:"type"`
	Pattern  string `json:"pattern"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// ScanResult contains all findings from a security scan
type ScanResult struct {
	Findings []Finding `json:"findings"`
}

// ScanCode scans code for secrets and returns detected issues
func ScanCode(code string, filePath string) *ScanResult {
	result := &ScanResult{
		Findings: []Finding{},
	}

	// Get all secret detection patterns
	patterns := GetSecretPatterns()

	// Split code into lines for line-by-line scanning
	lines := strings.Split(code, "\n")

	// Scan each line with each pattern
	for lineNum, line := range lines {
		for _, pattern := range patterns {
			if pattern.Pattern.MatchString(line) {
				// Determine severity based on secret type
				severity := "high"
				if pattern.Type == "password" || pattern.Type == "private_key" {
					severity = "critical"
				}

				finding := Finding{
					File:     filePath,
					Line:     lineNum + 1, // Line numbers are 1-indexed
					Type:     pattern.Type,
					Pattern:  pattern.Pattern.String(),
					Severity: severity,
					Message:  pattern.Message,
				}

				result.Findings = append(result.Findings, finding)
			}
		}
	}

	return result
}

// RedactSecrets replaces secrets in code with placeholders
// Example: "sk-abc123" → "[API_KEY_REDACTED]"
func RedactSecrets(code string) string {
	redacted := code
	patterns := GetSecretPatterns()

	for _, pattern := range patterns {
		placeholder := "[" + strings.ToUpper(pattern.Type) + "_REDACTED]"
		redacted = pattern.Pattern.ReplaceAllString(redacted, placeholder)
	}

	return redacted
}
