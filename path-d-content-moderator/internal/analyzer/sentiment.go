package analyzer

import (
	"math"
	"regexp"
	"strings"

	"github.com/munich-gophers/ai-workshop/content-moderator/internal/models"
)

// Sentiment word lists for pattern-based analysis
var (
	positiveWords = []string{
		"good", "great", "excellent", "amazing", "wonderful", "fantastic",
		"love", "perfect", "best", "awesome", "brilliant", "outstanding",
		"happy", "pleased", "satisfied", "delighted", "impressed", "thank",
		"appreciate", "helpful", "quality", "recommend", "superb", "terrific",
	}

	negativeWords = []string{
		"bad", "terrible", "awful", "horrible", "worst", "poor",
		"hate", "disappointing", "disappointed", "frustrated", "angry", "upset",
		"useless", "waste", "broken", "failed", "failure", "problem",
		"issue", "bug", "slow", "crash", "error", "never", "wrong",
	}

	intensifiers = []string{
		"very", "extremely", "absolutely", "totally", "completely", "really",
		"incredibly", "exceptionally", "remarkably",
	}

	negations = []string{
		"not", "no", "never", "nothing", "nowhere", "neither", "nobody",
		"none", "don't", "doesn't", "didn't", "won't", "wouldn't", "can't",
		"cannot", "couldn't",
	}
)

// AnalyzeSentiment performs pattern-based sentiment analysis
func AnalyzeSentiment(content string) models.SentimentScore {
	// Normalize content
	content = strings.ToLower(content)
	words := strings.Fields(content)

	var positiveScore, negativeScore float64
	var totalWords int

	// Check for emojis
	positiveScore += countPositiveEmojis(content) * 1.5
	negativeScore += countNegativeEmojis(content) * 1.5

	// Analyze each word
	for i, word := range words {
		totalWords++

		// Clean punctuation
		word = cleanWord(word)

		// Check for intensifiers before the word
		intensity := 1.0
		if i > 0 {
			prevWord := cleanWord(words[i-1])
			if contains(intensifiers, prevWord) {
				intensity = 1.5
			}
		}

		// Check for negations before the word
		isNegated := false
		if i > 0 {
			prevWord := cleanWord(words[i-1])
			if contains(negations, prevWord) {
				isNegated = true
			}
		}

		// Score positive words
		if contains(positiveWords, word) {
			if isNegated {
				negativeScore += 1.0 * intensity
			} else {
				positiveScore += 1.0 * intensity
			}
		}

		// Score negative words
		if contains(negativeWords, word) {
			if isNegated {
				positiveScore += 0.5 * intensity
			} else {
				negativeScore += 1.0 * intensity
			}
		}

		// Check for exclamation marks (increase intensity)
		if strings.Contains(words[i], "!") {
			if positiveScore > negativeScore {
				positiveScore += 0.3
			} else if negativeScore > positiveScore {
				negativeScore += 0.3
			}
		}

		// Check for question marks (slight negative bias for uncertainty)
		if strings.Contains(words[i], "?") {
			negativeScore += 0.1
		}
	}

	// Calculate final score
	totalScore := positiveScore + negativeScore
	if totalScore == 0 {
		// Neutral by default
		return models.SentimentScore{
			Label:      "neutral",
			Confidence: 0.5,
			Score:      0.0,
		}
	}

	// Normalize score between -1 and 1
	normalizedScore := (positiveScore - negativeScore) / math.Max(totalScore, 1.0)

	// Determine label and confidence
	var label string
	var confidence float64

	if normalizedScore > 0.2 {
		label = "positive"
		confidence = math.Min((normalizedScore-0.2)/0.8, 1.0)
	} else if normalizedScore < -0.2 {
		label = "negative"
		confidence = math.Min((math.Abs(normalizedScore)-0.2)/0.8, 1.0)
	} else {
		label = "neutral"
		confidence = 1.0 - math.Abs(normalizedScore)/0.2
	}

	// Ensure minimum confidence
	confidence = math.Max(confidence, 0.3)

	return models.SentimentScore{
		Label:      label,
		Confidence: confidence,
		Score:      normalizedScore,
	}
}

// Helper functions
func cleanWord(word string) string {
	// Remove common punctuation
	reg := regexp.MustCompile(`[^\w]`)
	return reg.ReplaceAllString(strings.ToLower(word), "")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func countPositiveEmojis(content string) float64 {
	positiveEmojis := []string{
		"😀", "😃", "😄", "😁", "😊", "🙂", "😍", "🥰", "😘",
		"👍", "👏", "🎉", "❤️", "💯", "✅", "⭐",
	}

	count := 0.0
	for _, emoji := range positiveEmojis {
		count += float64(strings.Count(content, emoji))
	}
	return count
}

func countNegativeEmojis(content string) float64 {
	negativeEmojis := []string{
		"😠", "😡", "😤", "😞", "😢", "😭", "😩", "😫",
		"👎", "💔", "❌", "⛔",
	}

	count := 0.0
	for _, emoji := range negativeEmojis {
		count += float64(strings.Count(content, emoji))
	}
	return count
}
