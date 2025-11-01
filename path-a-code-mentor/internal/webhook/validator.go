package webhook

import (
	"github.com/google/go-github/v76/github"
	"io"
	"net/http"
	"os"
)

// ValidateGitHubWebhook validates the webhook signature using go-github
func ValidateGitHubWebhook(r *http.Request) ([]byte, error) {
	secret := []byte(os.Getenv("GITHUB_WEBHOOK_SECRET"))

	// If no secret is configured, skip validation (development only!)
	if len(secret) == 0 {
		// Read body without validation
		return io.ReadAll(r.Body)
	}

	// Use go-github's built-in validation
	return github.ValidatePayload(r, secret)
}

// GetWebhookType returns the type of GitHub webhook event
func GetWebhookType(r *http.Request) string {
	return github.WebHookType(r)
}
