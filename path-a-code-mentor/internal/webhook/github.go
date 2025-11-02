package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/go-github/v76/github"
	"github.com/munich-gophers/ai-workshop/code-mentor/internal/analyzer"
	"github.com/munich-gophers/ai-workshop/code-mentor/internal/models"
	"github.com/munich-gophers/ai-workshop/code-mentor/internal/security"
)

// Config holds webhook configuration
type Config struct {
	Secret       string
	GitHubClient *github.Client
}

// HandleGitHub processes GitHub webhook events
func HandleGitHub(a *analyzer.Analyzer, cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Validate webhook payload
		payload, err := github.ValidatePayload(r, []byte(cfg.Secret))
		if err != nil {
			log.Printf("❌ Invalid webhook signature: %v", err)
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		// Step 2: Parse the webhook event
		webhookType := github.WebHookType(r)
		parsedEvent, err := github.ParseWebHook(webhookType, payload)
		if err != nil {
			log.Printf("❌ Failed to parse webhook: %v", err)
			http.Error(w, "Invalid webhook payload", http.StatusBadRequest)
			return
		}

		// Handle ping events (GitHub sends this when you first set up the webhook)
		if webhookType == "ping" {
			log.Println("🏓 Received ping event")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "pong",
			})
			return
		}

		// Step 3: Type assert to *github.PullRequestEvent
		event, ok := parsedEvent.(*github.PullRequestEvent)
		if !ok {
			log.Printf("⚠️  Ignoring non-PR event: %s", webhookType)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": fmt.Sprintf("Ignoring %s event", webhookType),
			})
			return
		}

		// Step 4: Only process "opened" and "synchronize" actions
		action := event.GetAction()
		if action != "opened" && action != "synchronize" {
			log.Printf("⚠️  Ignoring PR action: %s", action)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": fmt.Sprintf("Ignoring action: %s", action),
			})
			return
		}

		// Log the PR details
		pr := event.GetPullRequest()
		repo := event.GetRepo()
		log.Printf("🔍 Processing PR #%d: %s in %s",
			pr.GetNumber(),
			pr.GetTitle(),
			repo.GetFullName(),
		)

		// Generate analysis ID for tracking
		analysisID := generateAnalysisID(pr.GetNumber())

		// Process asynchronously to respond quickly to GitHub
		go processPullRequest(r.Context(), event, a, cfg)

		// Respond immediately to GitHub
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":     "Processing PR review",
			"pr_number":   pr.GetNumber(),
			"action":      action,
			"analysis_id": analysisID,
		})
	}
}

// processPullRequest handles the PR review asynchronously
func processPullRequest(ctx context.Context, event *github.PullRequestEvent, a *analyzer.Analyzer, cfg *Config) {
	pr := event.GetPullRequest()
	repo := event.GetRepo()

	// Step 5: Get the PR diff
	diff, err := fetchPRDiff(ctx, cfg.GitHubClient, repo.GetOwner().GetLogin(), repo.GetName(), pr.GetNumber())
	if err != nil {
		log.Printf("❌ Failed to fetch PR diff: %v", err)
		return
	}

	// Step 6: Scan for secrets
	secretResult := security.ScanCode(diff, "pull_request_diff")
	if len(secretResult.Findings) > 0 {
		log.Printf("🚨 Found %d potential secrets in PR #%d", len(secretResult.Findings), pr.GetNumber())

		// Post a warning comment about secrets
		if err := postSecurityWarning(ctx, cfg.GitHubClient, repo, pr, secretResult); err != nil {
			log.Printf("❌ Failed to post security warning: %v", err)
		}

		// For workshop: we'll still continue with redacted review
		// In production, you might want to block the PR
	}

	// Step 7: Redact secrets before AI processing
	redactedDiff := security.RedactSecrets(diff)

	// Step 8: Call analyzer to review the code
	review, err := a.Review(ctx, models.ReviewRequest{
		PRNumber:    pr.GetNumber(),
		Title:       pr.GetTitle(),
		Description: pr.GetBody(),
		Author:      pr.GetUser().GetLogin(),
		Diff:        redactedDiff,
		RepoName:    repo.GetFullName(),
		Branch:      pr.GetHead().GetRef(),
	})
	if err != nil {
		log.Printf("❌ Failed to analyze PR: %v", err)
		return
	}

	// Step 9: Post the review to GitHub (convert ReviewResponse to Review)
	if err := postReview(ctx, cfg.GitHubClient, repo, pr, review.ToReview()); err != nil {
		log.Printf("❌ Failed to post review: %v", err)
		return
	}

	log.Printf("✅ Successfully reviewed PR #%d", pr.GetNumber())
}

// fetchPRDiff retrieves the unified diff for a pull request
func fetchPRDiff(ctx context.Context, client *github.Client, owner, repo string, prNumber int) (string, error) {
	// Fetch all files in the PR
	opts := &github.ListOptions{PerPage: 100}
	var allFiles []*github.CommitFile

	for {
		files, resp, err := client.PullRequests.ListFiles(ctx, owner, repo, prNumber, opts)
		if err != nil {
			return "", fmt.Errorf("failed to list PR files: %w", err)
		}

		allFiles = append(allFiles, files...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// Build a unified diff from all files
	var diff string
	for _, file := range allFiles {
		if file.Patch != nil {
			diff += fmt.Sprintf("diff --git a/%s b/%s\n", file.GetFilename(), file.GetFilename())
			diff += fmt.Sprintf("--- a/%s\n", file.GetFilename())
			diff += fmt.Sprintf("+++ b/%s\n", file.GetFilename())
			diff += file.GetPatch() + "\n"
		}
	}

	log.Printf("📁 Fetched diff with %d files (%d bytes)", len(allFiles), len(diff))
	return diff, nil
}

// postSecurityWarning posts a comment warning about potential secrets
func postSecurityWarning(ctx context.Context, client *github.Client, repo *github.Repository, pr *github.PullRequest, result *security.ScanResult) error {
	comment := buildSecurityComment(result)

	issueComment := &github.IssueComment{
		Body: github.String(comment),
	}

	_, _, err := client.Issues.CreateComment(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		pr.GetNumber(),
		issueComment,
	)

	return err
}

// buildSecurityComment formats security findings as markdown
func buildSecurityComment(result *security.ScanResult) string {
	comment := "## 🚨 Security Alert: Potential Secrets Detected\n\n"
	comment += "⚠️ **WARNING**: This PR may contain sensitive information that should not be committed.\n\n"
	comment += "### Findings:\n\n"

	for i, finding := range result.Findings {
		comment += fmt.Sprintf("%d. **%s** detected in `%s` (line %d)\n",
			i+1,
			finding.Type,
			finding.File,
			finding.Line,
		)
		comment += fmt.Sprintf("   - Pattern: `%s`\n", finding.Pattern)
		comment += fmt.Sprintf("   - Severity: %s\n\n", finding.Severity)
	}

	comment += "### ⚡ Action Required:\n\n"
	comment += "1. Remove any secrets from the code\n"
	comment += "2. Use environment variables or secret management systems\n"
	comment += "3. If secrets were already committed, rotate them immediately\n"
	comment += "4. Consider using tools like `git-secrets` or `gitleaks` locally\n\n"
	comment += "_This is an automated security scan. Please review carefully._"

	return comment
}

// postReview posts the AI review to GitHub
func postReview(ctx context.Context, client *github.Client, repo *github.Repository, pr *github.PullRequest, review *models.Review) error {
	// Format the review as a comment
	comment := buildReviewComment(review)

	issueComment := &github.IssueComment{
		Body: github.String(comment),
	}

	createdComment, _, err := client.Issues.CreateComment(
		ctx,
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		pr.GetNumber(),
		issueComment,
	)

	if err != nil {
		return err
	}

	log.Printf("💬 Posted review: %s", createdComment.GetHTMLURL())
	return nil
}

// buildReviewComment formats the AI review as markdown
func buildReviewComment(review *models.Review) string {
	comment := "## 🤖 AI Code Review\n\n"

	// Add summary
	comment += "### Summary\n\n"
	comment += review.Summary + "\n\n"

	// Add detailed suggestions
	if len(review.Suggestions) > 0 {
		comment += "### 💡 Suggestions\n\n"
		for i, suggestion := range review.Suggestions {
			emoji := getSeverityEmoji(suggestion.Severity)
			comment += fmt.Sprintf("%d. %s **%s**", i+1, emoji, suggestion.Title)

			if suggestion.File != "" {
				comment += fmt.Sprintf(" (`%s`", suggestion.File)
				if suggestion.Line > 0 {
					comment += fmt.Sprintf(":%d", suggestion.Line)
				}
				comment += ")"
			}
			comment += "\n\n"

			if suggestion.Description != "" {
				comment += "   " + suggestion.Description + "\n\n"
			}

			if suggestion.CodeExample != "" {
				comment += "   ```go\n"
				comment += "   " + suggestion.CodeExample + "\n"
				comment += "   ```\n\n"
			}
		}
	}

	// Add positive feedback
	if len(review.PositiveFeedback) > 0 {
		comment += "### ✨ What's Good\n\n"
		for _, feedback := range review.PositiveFeedback {
			comment += fmt.Sprintf("- %s\n", feedback)
		}
		comment += "\n"
	}

	// Add overall assessment
	if review.Approved {
		comment += "---\n✅ **Looks good to merge!**\n"
	} else {
		comment += "---\n💭 **Please review the suggestions above.**\n"
	}

	comment += "\n_Generated by AI Code Mentor powered by Gemini_"

	return comment
}

// getSeverityEmoji returns an emoji for a given severity level
func getSeverityEmoji(severity string) string {
	switch severity {
	case "critical":
		return "🚨"
	case "warning":
		return "⚠️"
	case "suggestion":
		return "💡"
	default:
		return "💬"
	}
}

// generateAnalysisID creates a unique ID for tracking PR analysis
func generateAnalysisID(prNumber int) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("analysis-%d-%d", prNumber, rand.Intn(100000))
}
