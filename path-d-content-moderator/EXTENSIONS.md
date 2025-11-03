# Path D: Content Moderator Extensions

Ideas to extend your AI content moderation service after the workshop.

## CUSTOM MODERATION POLICIES

1. **Community-Specific Rule Sets**

Different communities need different moderation standards:

```go
type ModerationPolicy struct {
    ID                 string
    Name               string
    CommunityID        string
    EnabledCategories  []string
    CategoryThresholds map[string]float64
    AutoActionRules    []AutoActionRule
    AppealProcess      bool
}

type AutoActionRule struct {
    Condition  string  // "confidence > 0.9 AND category = hate-speech"
    Action     string  // "reject", "shadow-ban", "flag", "escalate"
    NotifyUser bool
}

func (m *Moderator) ApplyPolicy(policy ModerationPolicy, content ContentRequest) (*ComprehensiveAnalysis, error) {
    // Standard moderation
    analysis, _ := m.AnalyzeComprehensive(ctx, AnalysisRequest{
        ContentRequest: content,
    })

    // Apply policy-specific thresholds
    for category, threshold := range policy.CategoryThresholds {
        for i, cat := range analysis.Moderation.Categories {
            if cat.Category == category && cat.Confidence > threshold {
                analysis.Moderation.Categories[i].Flagged = true
            }
        }
    }

    // Evaluate auto-action rules
    for _, rule := range policy.AutoActionRules {
        if evaluateRule(rule.Condition, analysis) {
            analysis.Recommendation.Action = rule.Action
            analysis.Recommendation.AutoExecute = true
            analysis.Recommendation.Reason = fmt.Sprintf("Policy rule triggered: %s", rule.Condition)
        }
    }

    return analysis, nil
}
```

**Policy Examples:**
```go
var strictPolicy = ModerationPolicy{
    Name: "Strict - Family Friendly",
    CategoryThresholds: map[string]float64{
        "profanity":   0.3,  // Low tolerance
        "sexual":      0.2,
        "hate-speech": 0.1,
        "violence":    0.2,
    },
    AutoActionRules: []AutoActionRule{
        {
            Condition:  "confidence > 0.7",
            Action:     "reject",
            NotifyUser: true,
        },
    },
}

var relaxedPolicy = ModerationPolicy{
    Name: "Relaxed - Adult Community",
    CategoryThresholds: map[string]float64{
        "profanity":   0.8,  // High tolerance
        "sexual":      0.9,
        "hate-speech": 0.3,  // Still strict on hate
        "violence":    0.7,
    },
}
```

---

## MULTI-LANGUAGE CONTENT ANALYSIS

2. **Support Content in Multiple Languages**

```go
import "cloud.google.com/go/translate"

func (m *Moderator) ModerateMultiLanguage(ctx context.Context, req ContentRequest) (*ModerationResult, error) {
    // Detect language
    translateClient, _ := translate.NewClient(ctx)
    detections, _ := translateClient.DetectLanguage(ctx, []string{req.Content})
    detectedLang := detections[0][0].Language.String()

    // Translate to English for analysis if needed
    content := req.Content
    if detectedLang != "en" {
        translations, _ := translateClient.Translate(
            ctx,
            []string{req.Content},
            language.English,
            &translate.Options{
                Source: language.Make(detectedLang),
            },
        )
        content = translations[0].Text
    }

    // Moderate translated content
    result, _ := m.Moderate(ctx, ContentRequest{
        Content:    content,
        ContentID:  req.ContentID,
        Author:     req.Author,
        ContextURL: req.ContextURL,
    })

    // Add language metadata
    result.DetectedLanguage = detectedLang
    result.TranslatedContent = content

    return result, nil
}
```

**Supported Languages:**
- Automatic language detection
- Translation to English for consistent moderation
- Preserve original content for appeals

---

## IMAGE & VIDEO MODERATION

3. **Extend to Visual Content**

```go
import vision "cloud.google.com/go/vision/apiv1"
import visionpb "cloud.google.com/go/vision/apiv1/visionpb"

type ImageModerationResult struct {
    ImageURL         string
    SafeSearchResult SafeSearch
    ExplicitContent  []ExplicitAnnotation
    Violence         float64
    Racy             float64
    Adult            float64
    Medical          float64
    Spoof            float64
}

type SafeSearch struct {
    Adult    string  // "VERY_LIKELY", "LIKELY", "POSSIBLE", "UNLIKELY", "VERY_UNLIKELY"
    Violence string
    Racy     string
    Medical  string
    Spoof    string
}

func (m *Moderator) ModerateImage(ctx context.Context, imageURL string) (*ImageModerationResult, error) {
    client, _ := vision.NewImageAnnotatorClient(ctx)

    image := vision.NewImageFromURI(imageURL)

    // Safe search detection
    safeSearch, _ := client.DetectSafeSearch(ctx, image, nil)

    result := &ImageModerationResult{
        ImageURL: imageURL,
        SafeSearchResult: SafeSearch{
            Adult:    safeSearch.Adult.String(),
            Violence: safeSearch.Violence.String(),
            Racy:     safeSearch.Racy.String(),
            Medical:  safeSearch.Medical.String(),
            Spoof:    safeSearch.Spoof.String(),
        },
    }

    // Convert likelihood to scores
    result.Adult = likelihoodToScore(safeSearch.Adult)
    result.Violence = likelihoodToScore(safeSearch.Violence)
    result.Racy = likelihoodToScore(safeSearch.Racy)
    result.Medical = likelihoodToScore(safeSearch.Medical)
    result.Spoof = likelihoodToScore(safeSearch.Spoof)

    return result, nil
}

func likelihoodToScore(likelihood visionpb.Likelihood) float64 {
    scores := map[visionpb.Likelihood]float64{
        visionpb.Likelihood_VERY_UNLIKELY: 0.1,
        visionpb.Likelihood_UNLIKELY:      0.3,
        visionpb.Likelihood_POSSIBLE:      0.5,
        visionpb.Likelihood_LIKELY:        0.7,
        visionpb.Likelihood_VERY_LIKELY:   0.9,
    }
    return scores[likelihood]
}
```

**Video Moderation:**
```go
func (m *Moderator) ModerateVideo(ctx context.Context, videoURL string) (*VideoModerationResult, error) {
    // Extract frames at intervals
    frames := extractFrames(videoURL, 5*time.Second) // Every 5 seconds

    var results []ImageModerationResult
    for _, frame := range frames {
        result, _ := m.ModerateImage(ctx, frame.URL)
        results = append(results, *result)
    }

    // Aggregate results
    return aggregateFrameResults(results), nil
}
```

---

## APPEAL WORKFLOW SYSTEM

4. **Let Users Contest Moderation Decisions**

```go
type Appeal struct {
    ID               string
    ContentID        string
    OriginalDecision ComprehensiveAnalysis
    UserExplanation  string
    Status           string // "pending", "under-review", "upheld", "overturned"
    ReviewedBy       string
    ReviewedAt       time.Time
    FinalDecision    string
    ReviewNotes      string
}

func (m *Moderator) SubmitAppeal(contentID, explanation string) (*Appeal, error) {
    appeal := &Appeal{
        ID:              generateID(),
        ContentID:       contentID,
        UserExplanation: explanation,
        Status:          "pending",
    }

    // Store in database
    m.db.SaveAppeal(appeal)

    // Trigger review workflow
    m.notifyReviewers(appeal)

    return appeal, nil
}

func (m *Moderator) ReviewAppeal(appealID, reviewerID string) (*Appeal, error) {
    appeal, _ := m.db.GetAppeal(appealID)

    // Re-analyze with human context
    prompt := fmt.Sprintf(`
Re-evaluate this content moderation decision with the user's appeal context:

Original Content: %s
Original Decision: %s
Confidence: %.2f

User's Explanation: %s

Considering the user's context, should this decision be:
1. Upheld (original decision correct)
2. Overturned (original decision wrong)

Provide reasoning.
`, appeal.OriginalContent, appeal.OriginalDecision.Recommendation.Action,
   appeal.OriginalDecision.Recommendation.Confidence, appeal.UserExplanation)

    // Get AI-assisted review
    response, _ := genkit.Generate(ctx, m.genkit, ai.WithPrompt(prompt))

    appeal.Status = "under-review"
    appeal.ReviewedBy = reviewerID
    appeal.ReviewedAt = time.Now()
    // Human makes final decision with AI assistance

    return appeal, nil
}
```

---

## MODERATION DASHBOARD

5. **Build Analytics & Monitoring Interface**

```go
type ModerationDashboard struct {
    TotalReviewed      int
    FlaggedRate        float64
    TopCategories      []CategoryStats
    TrendsByDay        []TrendPoint
    ModeratorMetrics   []ModeratorStats
    AppealRate         float64
    OverturnRate       float64
    AvgConfidence      float64
}

type CategoryStats struct {
    Category   string
    Count      int
    Percentage float64
}

type TrendPoint struct {
    Date       time.Time
    Flagged    int
    Total      int
    Categories map[string]int
}

func (m *Moderator) GenerateDashboard(startDate, endDate time.Time) (*ModerationDashboard, error) {
    results, _ := m.db.GetModerationResults(startDate, endDate)

    dashboard := &ModerationDashboard{
        TotalReviewed: len(results),
    }

    // Calculate metrics
    flaggedCount := 0
    categoryCount := make(map[string]int)
    totalConfidence := 0.0

    for _, result := range results {
        if result.IsFlagged {
            flaggedCount++
        }

        for _, cat := range result.Categories {
            if cat.Flagged {
                categoryCount[cat.Category]++
            }
        }

        totalConfidence += result.Recommendation.Confidence
    }

    dashboard.FlaggedRate = float64(flaggedCount) / float64(len(results)) * 100
    dashboard.AvgConfidence = totalConfidence / float64(len(results))

    // Build category stats
    for category, count := range categoryCount {
        dashboard.TopCategories = append(dashboard.TopCategories, CategoryStats{
            Category:   category,
            Count:      count,
            Percentage: float64(count) / float64(len(results)) * 100,
        })
    }

    // Sort by count
    sort.Slice(dashboard.TopCategories, func(i, j int) bool {
        return dashboard.TopCategories[i].Count > dashboard.TopCategories[j].Count
    })

    return dashboard, nil
}
```

**Dashboard Endpoints:**
```go
func (h *Handler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
    endDate := time.Now()
    startDate := endDate.AddDate(0, 0, -30) // Last 30 days

    dashboard, _ := h.moderator.GenerateDashboard(startDate, endDate)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(dashboard)
}
```

---

## AUTOMATED CONTENT CATEGORIZATION

6. **Smart Content Classification**

```go
type ContentCategory struct {
    Type        string  // "question", "complaint", "spam", "feedback", "support"
    Confidence  float64
    Tags        []string
    Topics      []string
    Intent      string
}

func (m *Moderator) CategorizeContent(ctx context.Context, content string) (*ContentCategory, error) {
    prompt := `
Analyze and categorize this user-generated content:

Content: ` + content + `

Classify it:
1. Type: question, complaint, spam, feedback, support, discussion, announcement
2. Tags: Relevant tags (e.g., "technical", "billing", "feature-request")
3. Topics: Main topics discussed
4. Intent: User's primary intent

Provide structured JSON response.
`

    response, _ := genkit.Generate(ctx, m.genkit, ai.WithPrompt(prompt))

    var category ContentCategory
    // Parse AI response into category
    return &category, nil
}
```

**Use Cases:**
- Auto-routing to appropriate teams
- Priority assignment
- Topic tracking
- Community insights

---

## TOXICITY SCORING & TRENDS

7. **Track Content Health Over Time**

```go
type ToxicityTrend struct {
    Period           string
    AvgToxicity      float64
    PeakToxicity     float64
    ToxicPosts       int
    TotalPosts       int
    TopOffenders     []UserStats
    ImprovedUsers    []UserStats
}

type UserStats struct {
    UserID          string
    Username        string
    PostCount       int
    AvgToxicity     float64
    FlaggedCount    int
    TrendDirection  string // "improving", "worsening", "stable"
}

func (m *Moderator) AnalyzeToxicityTrends(startDate, endDate time.Time) (*ToxicityTrend, error) {
    results, _ := m.db.GetResultsByDateRange(startDate, endDate)

    trend := &ToxicityTrend{
        Period:     fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
        TotalPosts: len(results),
    }

    userStats := make(map[string]*UserStats)
    totalToxicity := 0.0

    for _, result := range results {
        toxicity := calculateToxicity(result)
        totalToxicity += toxicity

        if toxicity > trend.PeakToxicity {
            trend.PeakToxicity = toxicity
        }

        if toxicity > 0.7 {
            trend.ToxicPosts++
        }

        // Track per-user stats
        if stats, exists := userStats[result.Author]; exists {
            stats.PostCount++
            stats.AvgToxicity = (stats.AvgToxicity*float64(stats.PostCount-1) + toxicity) / float64(stats.PostCount)
            if result.IsFlagged {
                stats.FlaggedCount++
            }
        } else {
            userStats[result.Author] = &UserStats{
                UserID:       result.Author,
                PostCount:    1,
                AvgToxicity:  toxicity,
                FlaggedCount: 0,
            }
        }
    }

    trend.AvgToxicity = totalToxicity / float64(len(results))

    // Identify top offenders and improved users
    // ... (sorting and filtering logic)

    return trend, nil
}

func calculateToxicity(result ComprehensiveAnalysis) float64 {
    // Weighted toxicity score
    weights := map[string]float64{
        "hate-speech":     1.0,
        "harassment":      0.9,
        "violence":        0.8,
        "sexual":          0.7,
        "profanity":       0.3,
        "spam":            0.2,
    }

    toxicity := 0.0
    for _, cat := range result.Moderation.Categories {
        if cat.Flagged {
            toxicity += weights[cat.Category] * cat.Confidence
        }
    }

    return toxicity
}
```

---

## CONTEXT-AWARE MODERATION

8. **Consider Thread/Conversation Context**

```go
type ConversationContext struct {
    ThreadID       string
    ParentPosts    []string
    Participants   []string
    Topic          string
    OverallTone    string
}

func (m *Moderator) ModerateWithContext(ctx context.Context, req ContentRequest, conversation ConversationContext) (*ComprehensiveAnalysis, error) {
    // Build enriched prompt with context
    prompt := fmt.Sprintf(`
Moderate this content considering the conversation context:

Current Post: %s
Author: %s

Conversation Thread:
%s

Topic: %s
Overall Tone: %s

Consider:
- Is this response appropriate given the thread context?
- Is it targeted harassment or legitimate debate?
- Does it escalate or de-escalate tension?

Provide moderation analysis with context-aware reasoning.
`, req.Content, req.Author,
   formatThreadHistory(conversation.ParentPosts),
   conversation.Topic,
   conversation.OverallTone)

    response, _ := genkit.Generate(ctx, m.genkit, ai.WithPrompt(prompt))

    // Parse and return context-aware analysis
    return analysis, nil
}
```

---

## AUTOMATED ACTIONS & WEBHOOKS

9. **Trigger Actions Based on Moderation Results**

```go
type ModerationWebhook struct {
    URL         string
    Events      []string // "content.flagged", "user.warned", "content.deleted"
    Headers     map[string]string
    RetryPolicy RetryPolicy
}

type RetryPolicy struct {
    MaxRetries int
    BackoffMs  int
}

func (m *Moderator) TriggerWebhooks(analysis *ComprehensiveAnalysis) error {
    webhooks, _ := m.db.GetActiveWebhooks()

    for _, webhook := range webhooks {
        event := determineEvent(analysis)

        if contains(webhook.Events, event) {
            go m.sendWebhook(webhook, analysis, event)
        }
    }

    return nil
}

func (m *Moderator) sendWebhook(webhook ModerationWebhook, analysis *ComprehensiveAnalysis, event string) error {
    payload := map[string]interface{}{
        "event":      event,
        "timestamp":  time.Now().Unix(),
        "content_id": analysis.ContentID,
        "author":     analysis.Author,
        "decision":   analysis.Recommendation,
        "categories": analysis.Moderation.Categories,
    }

    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(body))

    for key, value := range webhook.Headers {
        req.Header.Set(key, value)
    }

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)

    // Implement retry logic
    if err != nil || resp.StatusCode >= 500 {
        m.retryWebhook(webhook, payload, webhook.RetryPolicy)
    }

    return nil
}
```

---

## FALSE POSITIVE LEARNING

10. **Improve Accuracy with Feedback Loops**

```go
type FeedbackLoop struct {
    ContentID        string
    OriginalDecision ComprehensiveAnalysis
    ActualOutcome    string // "correct", "false-positive", "false-negative"
    ReviewerNotes    string
    Timestamp        time.Time
}

func (m *Moderator) RecordFeedback(contentID, outcome, notes string) error {
    feedback := &FeedbackLoop{
        ContentID:     contentID,
        ActualOutcome: outcome,
        ReviewerNotes: notes,
        Timestamp:     time.Now(),
    }

    // Store feedback
    m.db.SaveFeedback(feedback)

    // Use feedback to improve future decisions
    if outcome == "false-positive" {
        m.adjustThresholds(feedback)
    }

    return nil
}

func (m *Moderator) adjustThresholds(feedback *FeedbackLoop) {
    // Analyze patterns in false positives
    // Adjust category thresholds or add exceptions
    // Update prompt engineering based on common mistakes
}
```

---

## RESOURCES

**APIs & Services:**
- [Cloud Vision API](https://cloud.google.com/vision/docs) - Image/video moderation
- [Cloud Translation API](https://cloud.google.com/translate/docs) - Multi-language support
- [Perspective API](https://perspectiveapi.com/) - Toxicity detection
- [Sightengine](https://sightengine.com/) - Content moderation API

**Tools:**
- [Prometheus](https://prometheus.io/) - Metrics & monitoring
- [Grafana](https://grafana.com/) - Dashboard visualization
- [Redis](https://redis.io/) - Rate limiting & caching

**Libraries:**
- [go-github](https://github.com/google/go-github) - GitHub webhooks
- [felixgeelhaar/jirasdk](https://github.com/felixgeelhaar/jirasdk) - Issue tracking
- [go-chi](https://github.com/go-chi/chi) - HTTP routing

---

## NEXT STEPS

Pick one extension that interests you:

1. **Quick win** (1-2 hours): Custom policies, toxicity scoring
2. **Medium** (half day): Image moderation, appeal workflow
3. **Advanced** (full day): Multi-language support, context-aware moderation, dashboard

**Remember:**
- Start with clear moderation guidelines
- Always provide transparency to users
- Build appeals process from day one
- Monitor for bias in AI decisions
- Regularly review and adjust thresholds

**Privacy & Ethics:**
- Store minimal user data
- Be transparent about AI usage
- Allow users to export their data
- Implement fair appeal processes
- Regular bias audits
