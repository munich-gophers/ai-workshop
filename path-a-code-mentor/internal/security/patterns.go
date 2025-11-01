package security

import "regexp"

// SecretPattern represents a regex pattern for detecting secrets
type SecretPattern struct {
	Type    string
	Pattern *regexp.Regexp
	Message string
}

// GetSecretPatterns returns all secret detection patterns
func GetSecretPatterns() []SecretPattern {
	return []SecretPattern{
		{
			Type:    "api_key",
			Pattern: regexp.MustCompile(`(?i)(api[_-]?key|apikey|api[_-]?secret)[\s]*[=:]\s*["\']?([a-zA-Z0-9_\-]{20,})["\']?`),
			Message: "Potential API key detected",
		},
		{
			Type:    "api_key",
			Pattern: regexp.MustCompile(`sk-[a-zA-Z0-9]{20,}`), // OpenAI
			Message: "OpenAI API key detected",
		},
		{
			Type:    "api_key",
			Pattern: regexp.MustCompile(`AKIA[0-9A-Z]{16}`), // AWS
			Message: "AWS access key detected",
		},
		{
			Type:    "token",
			Pattern: regexp.MustCompile(`(?i)(token|bearer)[\s]*[=:]\s*["\']?([a-zA-Z0-9_\-\.]{20,})["\']?`),
			Message: "Potential authentication token detected",
		},
		{
			Type:    "token",
			Pattern: regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`), // GitHub
			Message: "GitHub personal access token detected",
		},
		{
			Type:    "password",
			Pattern: regexp.MustCompile(`(?i)(password|passwd|pwd)[\s]*[=:]\s*["\']([^"'\s]{8,})["\']`),
			Message: "Hardcoded password detected",
		},
		{
			Type:    "private_key",
			Pattern: regexp.MustCompile(`-----BEGIN\s+(RSA\s+)?PRIVATE\s+KEY-----`),
			Message: "Private key detected",
		},
		{
			Type:    "connection_string",
			Pattern: regexp.MustCompile(`(?i)(postgres|mysql|mongodb):\/\/[^\s]*:[^\s]*@`),
			Message: "Database connection string with credentials detected",
		},
	}
}
