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
	// TODO: CHECKPOINT 3 - Implement secret detection
	// Hint:
	// 1. Get patterns with GetSecretPatterns()
	// 2. Split code into lines
	// 3. For each line, check each pattern
	// 4. If match found, create Finding with line number
	// 5. Return ScanResult with all findings
	//
	// Example:
	//   result := &ScanResult{Findings: []Finding{}}
	//   patterns := GetSecretPatterns()
	//   lines := strings.Split(code, "\n")
	//   for lineNum, line := range lines {
	//       for _, pattern := range patterns {
	//           if pattern.Pattern.MatchString(line) {
	//               finding := Finding{
	//                   File: filePath,
	//                   Line: lineNum + 1,
	//                   Type: pattern.Type,
	//                   ...
	//               }
	//               result.Findings = append(result.Findings, finding)
	//           }
	//       }
	//   }
	//   return result

	// Your implementation here
	return &ScanResult{Findings: []Finding{}}
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
