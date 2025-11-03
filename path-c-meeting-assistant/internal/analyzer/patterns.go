package analyzer

import (
	"regexp"
	"strings"
	"time"

	"github.com/munich-gophers/ai-workshop/meeting-assistant/internal/models"
)

var (
	// Action item patterns - common phrases that indicate tasks
	actionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)action\s*item[s]?:?\s*(.+)`),
		regexp.MustCompile(`(?i)to\s*do:?\s*(.+)`),
		regexp.MustCompile(`(?i)TODO:?\s*(.+)`),
		regexp.MustCompile(`(?i)task[s]?:?\s*(.+)`),
		regexp.MustCompile(`(?i)follow\s*up:?\s*(.+)`),
		regexp.MustCompile(`(?i)(\w+)\s+will\s+(.+)`),
		regexp.MustCompile(`(?i)(\w+)\s+to\s+(.+)`),
		regexp.MustCompile(`(?i)(\w+)\s+needs\s+to\s+(.+)`),
		regexp.MustCompile(`(?i)(\w+)\s+should\s+(.+)`),
	}

	// Assignee patterns - extracting who is responsible
	assigneePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)assigned\s+to:?\s*(\w+)`),
		regexp.MustCompile(`(?i)@(\w+)`),
		regexp.MustCompile(`(?i)(\w+)\s+will\s+`),
		regexp.MustCompile(`(?i)(\w+)\s+to\s+`),
		regexp.MustCompile(`(?i)owner:?\s*(\w+)`),
	}

	// Due date patterns - extracting deadlines
	dueDatePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)by\s+(\w+\s+\d{1,2})`),
		regexp.MustCompile(`(?i)due\s+(\w+\s+\d{1,2})`),
		regexp.MustCompile(`(?i)deadline:?\s*(\w+\s+\d{1,2})`),
		regexp.MustCompile(`(?i)(\d{1,2}/\d{1,2}(?:/\d{2,4})?)`),
		regexp.MustCompile(`(?i)(next\s+week|this\s+week|tomorrow|today)`),
		regexp.MustCompile(`(?i)by\s+(end\s+of\s+week|EOW)`),
	}

	// Decision patterns - identifying decisions made
	decisionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)decided\s+to\s+(.+)`),
		regexp.MustCompile(`(?i)decision:?\s*(.+)`),
		regexp.MustCompile(`(?i)agreed\s+to\s+(.+)`),
		regexp.MustCompile(`(?i)we\s+will\s+(.+)`),
		regexp.MustCompile(`(?i)approved:?\s*(.+)`),
		regexp.MustCompile(`(?i)resolution:?\s*(.+)`),
	}

	// Priority indicators
	priorityPatterns = map[string]*regexp.Regexp{
		"high":     regexp.MustCompile(`(?i)(urgent|critical|high\s+priority|ASAP|immediately)`),
		"medium":   regexp.MustCompile(`(?i)(medium\s+priority|important|soon)`),
		"low":      regexp.MustCompile(`(?i)(low\s+priority|when\s+possible|nice\s+to\s+have)`),
	}
)

// ExtractActionItems extracts action items from meeting notes using pattern matching
func ExtractActionItems(notes string) []models.ActionItem {
	var actionItems []models.ActionItem
	lines := strings.Split(notes, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if line contains action item patterns
		for _, pattern := range actionPatterns {
			if matches := pattern.FindStringSubmatch(line); len(matches) > 1 {
				description := strings.TrimSpace(matches[len(matches)-1])
				if description == "" {
					continue
				}

				actionItem := models.ActionItem{
					Description: description,
					Status:      "pending",
					Priority:    extractPriority(line),
					ExtractedAt: time.Now(),
				}

				// Try to extract assignee
				if assignee := extractAssignee(line); assignee != "" {
					actionItem.Assignee = assignee
				}

				// Try to extract due date
				if dueDate := extractDueDate(line); dueDate != "" {
					actionItem.DueDate = dueDate
				}

				actionItems = append(actionItems, actionItem)
				break
			}
		}
	}

	return actionItems
}

// ExtractDecisions extracts decisions from meeting notes
func ExtractDecisions(notes string) []models.Decision {
	var decisions []models.Decision
	lines := strings.Split(notes, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if line contains decision patterns
		for _, pattern := range decisionPatterns {
			if matches := pattern.FindStringSubmatch(line); len(matches) > 1 {
				description := strings.TrimSpace(matches[1])
				if description == "" {
					continue
				}

				decision := models.Decision{
					Description: description,
				}

				decisions = append(decisions, decision)
				break
			}
		}
	}

	return decisions
}

// ExtractParticipants extracts participant names from meeting notes
func ExtractParticipants(notes string) []models.Participant {
	participantMap := make(map[string]*models.Participant)

	// Common name patterns (capitalized words)
	namePattern := regexp.MustCompile(`\b([A-Z][a-z]+(?:\s+[A-Z][a-z]+)?)\b`)
	matches := namePattern.FindAllString(notes, -1)

	for _, name := range matches {
		name = strings.TrimSpace(name)
		// Skip common words that might be capitalized
		if isCommonWord(name) {
			continue
		}

		if participant, exists := participantMap[name]; exists {
			participant.Mentions++
		} else {
			participantMap[name] = &models.Participant{
				Name:     name,
				Mentions: 1,
			}
		}
	}

	// Convert map to slice
	var participants []models.Participant
	for _, participant := range participantMap {
		// Only include participants mentioned at least twice
		if participant.Mentions >= 2 {
			participants = append(participants, *participant)
		}
	}

	return participants
}

// extractAssignee extracts the person responsible from a line
func extractAssignee(line string) string {
	for _, pattern := range assigneePatterns {
		if matches := pattern.FindStringSubmatch(line); len(matches) > 1 {
			assignee := strings.TrimSpace(matches[1])
			// Capitalize first letter
			if len(assignee) > 0 {
				assignee = strings.ToUpper(string(assignee[0])) + strings.ToLower(assignee[1:])
			}
			return assignee
		}
	}
	return ""
}

// extractDueDate extracts the due date from a line
func extractDueDate(line string) string {
	for _, pattern := range dueDatePatterns {
		if matches := pattern.FindStringSubmatch(line); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}
	return ""
}

// extractPriority determines the priority of an action item
func extractPriority(line string) string {
	for priority, pattern := range priorityPatterns {
		if pattern.MatchString(line) {
			return priority
		}
	}
	return "medium" // default priority
}

// isCommonWord checks if a capitalized word is a common word (not a name)
func isCommonWord(word string) bool {
	commonWords := map[string]bool{
		"The": true, "This": true, "That": true, "These": true, "Those": true,
		"We": true, "They": true, "It": true, "Action": true, "Item": true,
		"Meeting": true, "Notes": true, "Agenda": true, "Team": true,
		"Next": true, "Last": true, "First": true, "Second": true,
		"Monday": true, "Tuesday": true, "Wednesday": true, "Thursday": true,
		"Friday": true, "Saturday": true, "Sunday": true,
		"January": true, "February": true, "March": true, "April": true,
		"May": true, "June": true, "July": true, "August": true,
		"September": true, "October": true, "November": true, "December": true,
	}
	return commonWords[word]
}
