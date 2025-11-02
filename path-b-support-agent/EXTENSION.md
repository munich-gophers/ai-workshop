# Path B: Support Agent Extensions

Ideas to extend your AI support triage service after the workshop.

## SENTIMENT AND EMOTION

1. Sentiment Analysis

Detect frustrated or angry customers:

go
type SentimentResult struct {
Score float64 // -1.0 (very negative) to 1.0 (very positive)
Confidence float64
Emotions []string // "frustrated", "angry", "satisfied", "confused"
Urgency string // Adjusted based on emotion
}

func AnalyzeSentiment(ctx context.Context, message string) (\*SentimentResult, error) {
prompt := `Analyze the sentiment and emotions in this customer message.

Return JSON with:

- score: -1.0 (very negative) to 1.0 (very positive)
- confidence: 0.0 to 1.0
- emotions: array of detected emotions
- recommended_urgency: based on emotional state

Message: ` + message

    // Call Gemini
    response, err := model.Generate(ctx, genkit.NewUserTextMessage(prompt), nil)
    if err != nil {
        return nil, err
    }

    // Parse JSON
    var result SentimentResult
    json.Unmarshal([]byte(response.Text), &result)

    return &result, nil

}

Use cases:

- Escalate angry customers immediately
- Route frustrated users to senior agents
- Track customer satisfaction trends
- Trigger proactive outreach

Auto-escalate if:

- Score < -0.7 (very negative)
- Emotions include "angry" or "furious"
- Multiple frustration indicators

---

2. Empathy Scoring

Measure how empathetic agent responses are:

go
func ScoreEmpathy(agentResponse string) (float64, error) {
prompt := `Rate the empathy level of this support response from 0.0 (robotic) to 1.0 (highly empathetic).

Consider:

- Acknowledgment of customer feelings
- Understanding of their situation
- Personal touch vs template language
- Problem ownership

Response: ` + agentResponse

    // Returns empathy score
    // Use to train agents
    // Identify best practices

}

Track:

- Average empathy score per agent
- Correlation with satisfaction
- Best performing responses
- Training opportunities

## MULTILINGUAL SUPPORT

1. Auto-Translation

Detect and translate messages:

go
type MultilingualClassifier struct {
classifier *Classifier
translator *Translator
}

func (mc *MultilingualClassifier) ClassifyAnyLanguage(req SupportRequest) (*ClassificationResponse, error) {
// Detect language
detectedLang, confidence := detectLanguage(req.Message)

    // Store original
    originalMessage := req.Message
    originalLang := detectedLang

    // Translate to English if needed
    if detectedLang != "en" && confidence > 0.8 {
        translated, err := mc.translator.Translate(req.Message, detectedLang, "en")
        if err != nil {
            return nil, err
        }
        req.Message = translated
    }

    // Classify in English
    result, err := mc.classifier.Classify(req)
    if err != nil {
        return nil, err
    }

    // Add metadata
    result.OriginalLanguage = originalLang
    result.TranslationUsed = (originalLang != "en")

    return result, nil

}

Supported languages:

- Spanish, French, German
- Portuguese, Italian
- Japanese, Korean, Chinese
- And 100+ more

Benefits:

- Global support
- Single AI model
- Consistent quality
- Cost effective

---

2. Cultural Context

Adapt responses to cultural norms:

go
type CulturalContext struct {
Language string
Country string
Formality string // "formal", "neutral", "casual"
Greeting string
Closing string
}

func getC ulturalContext(language string, country string) *CulturalContext {
contexts := map[string]*CulturalContext{
"ja-JP": {
Formality: "formal",
Greeting: "いつもお世話になっております",
Closing: "よろしくお願いいたします",
},
"de-DE": {
Formality: "formal",
Greeting: "Sehr geehrte/r",
Closing: "Mit freundlichen Grüßen",
},
}

    return contexts[language+"-"+country]

}

Adapt:

- Formality level
- Greeting style
- Closing phrases
- Direct vs indirect communication
- Response length

## INTELLIGENT ROUTING

1. Similar Ticket Search

Find related resolved tickets:

go
import "github.com/firebase/genkit/go/ai"

type SimilarTicket struct {
ID string
Title string
Similarity float64
Resolution string
ResolvedBy string
TimeToResolve time.Duration
}

func FindSimilarTickets(ctx context.Context, message string, limit int) ([]SimilarTicket, error) {
// Generate embedding
embedding, err := ai.GenerateEmbedding(ctx, ai.EmbeddingRequest{
Text: message,
})
if err != nil {
return nil, err
}

    // Search vector database
    results, err := vectorDB.SimilaritySearch(embedding, limit)
    if err != nil {
        return nil, err
    }

    // Load ticket details
    tickets := []SimilarTicket{}
    for _, result := range results {
        ticket := loadTicketFromDB(result.ID)
        ticket.Similarity = result.Score
        tickets = append(tickets, ticket)
    }

    return tickets, nil

}

Benefits:

- Show agents proven solutions
- Reduce resolution time
- Improve consistency
- Knowledge base building

Display to agent:
"3 similar tickets found - average resolution time: 15 minutes"

---

2. Skill-Based Routing

Route to agent with right expertise:

go
type AgentSkills struct {
AgentID string
Languages []string
Products []string
Expertise []string // "billing", "technical", "account"
AvgRating float64
Available bool
}

func RouteToAgent(ticket *Ticket, agents []AgentSkills) *AgentSkills {
scored := []struct {
agent \*AgentSkills
score float64
}{}

    for _, agent := range agents {
        if !agent.Available {
            continue
        }

        score := 0.0

        // Language match
        if contains(agent.Languages, ticket.Language) {
            score += 10.0
        }

        // Expertise match
        if contains(agent.Expertise, ticket.Intent) {
            score += 20.0
        }

        // Rating bonus
        score += agent.AvgRating * 5.0

        scored = append(scored, struct {
            agent *AgentSkills
            score float64
        }{&agent, score})
    }

    // Sort by score
    sort.Slice(scored, func(i, j int) bool {
        return scored[i].score > scored[j].score
    })

    if len(scored) > 0 {
        return scored[0].agent
    }

    return nil

}

Route based on:

- Language proficiency
- Product knowledge
- Topic expertise
- Historical performance
- Current workload

## AUTOMATION

1. Auto-Response

Automatically handle common questions:

go
type AutoResponseRule struct {
ID string
TriggerKeywords []string
Intent string
MinConfidence float64
ResponseTemplate string
RequiresHuman bool
}

func CheckAutoResponse(ticket *Ticket) (*AutoResponseRule, error) {
rules := loadAutoResponseRules()

    for _, rule := range rules {
        // Check if intent matches
        if ticket.Intent != rule.Intent {
            continue
        }

        // Check confidence threshold
        if ticket.Confidence < rule.MinConfidence {
            continue
        }

        // Check keywords
        keywordMatch := false
        for _, keyword := range rule.TriggerKeywords {
            if strings.Contains(strings.ToLower(ticket.Message), keyword) {
                keywordMatch = true
                break
            }
        }

        if keywordMatch {
            return &rule, nil
        }
    }

    return nil, nil

}

Examples:

Password Reset:
Trigger: "forgot password", "can't login", "reset password"
Intent: account
Confidence: > 0.9
Response: "I'll send you a password reset link to {email}"
Auto-execute: Send reset email

Order Status:
Trigger: "where is my order", "track order", "shipping status"
Intent: other
Confidence: > 0.85
Response: "Your order #{order_id} is {status}. Expected delivery: {date}"
Auto-execute: Query order API

Refund Request:
Trigger: "refund", "money back", "cancel order"
Intent: billing
Confidence: > 0.8
Response: "I've started your refund request. You'll receive ${amount} within 5-7 business days."
RequiresHuman: true (for approval)

---

2. Smart Macros

Context-aware response templates:

go
type Macro struct {
Name string
Template string
Variables []string
Conditions map[string]string
}

func ApplyMacro(macro *Macro, ticket *Ticket, context map[string]string) string {
response := macro.Template

    // Replace variables
    for _, variable := range macro.Variables {
        value := context[variable]
        if value == "" {
            value = fetchFromContext(ticket, variable)
        }
        response = strings.ReplaceAll(response, "{"+variable+"}", value)
    }

    return response

}

Macro examples:

Escalation:
"I understand how frustrating this is. I'm escalating your case to our {team} team who specializes in {issue_type}. They'll reach out within {sla_time}."

Technical Issue:
"Thank you for reporting this. I've created ticket #{ticket_id} for our engineering team. We'll investigate and update you within {timeframe}."

Variables auto-filled from:

- Ticket metadata
- User profile
- System state
- SLA definitions

## ESCALATION LOGIC

1. Smart Escalation

Automatically escalate based on criteria:

go
type EscalationRule struct {
Name string
Priority int // Higher = more important
Condition func(\*Ticket) bool
TargetTeam string
SLA time.Duration
NotifySlack bool
}

var escalationRules = []EscalationRule{
{
Name: "angry-vip",
Priority: 100,
Condition: func(t _Ticket) bool {
return t.Sentiment.Score < -0.7 && t.UserTier == "VIP"
},
TargetTeam: "senior-support",
SLA: 15 _ time.Minute,
NotifySlack: true,
},
{
Name: "data-loss",
Priority: 90,
Condition: func(t _Ticket) bool {
keywords := []string{"lost data", "deleted", "missing files", "disappeared"}
for \_, kw := range keywords {
if strings.Contains(strings.ToLower(t.Message), kw) {
return true
}
}
return false
},
TargetTeam: "engineering",
SLA: 30 _ time.Minute,
NotifySlack: true,
},
{
Name: "legal-threat",
Priority: 95,
Condition: func(t _Ticket) bool {
legalWords := []string{"lawyer", "attorney", "lawsuit", "legal action", "sue"}
for \_, word := range legalWords {
if strings.Contains(strings.ToLower(t.Message), word) {
return true
}
}
return false
},
TargetTeam: "legal",
SLA: 1 _ time.Hour,
NotifySlack: true,
},
}

Escalate for:

- VIP customers with negative sentiment
- Security or privacy concerns
- Data loss reports
- Legal mentions
- Service outages
- Payment failures (high value)
- Repeated issues (3+ contacts)

## ANALYTICS AND INSIGHTS

1. Trend Detection

Identify emerging issues:

go
func DetectTrends(timeWindow time.Duration) []Trend {
tickets := getRecentTickets(timeWindow)

    // Group by intent
    intentCounts := make(map[string]int)
    for _, ticket := range tickets {
        intentCounts[ticket.Intent]++
    }

    // Compare to baseline
    baseline := getBaselineCounts()
    trends := []Trend{}

    for intent, count := range intentCounts {
        baselineCount := baseline[intent]
        increase := float64(count-baselineCount) / float64(baselineCount)

        if increase > 0.5 {  // 50% increase
            trends = append(trends, Trend{
                Intent:     intent,
                Increase:   increase,
                Count:      count,
                Timeframe:  timeWindow,
            })
        }
    }

    return trends

}

Detect:

- Sudden spike in specific issues
- Product problems affecting many users
- Service degradation patterns
- Feature confusion
- Documentation gaps

Alert when:

- 50% increase in technical issues
- 3x increase in billing questions
- New error messages appearing
- Negative sentiment trending

---

2. Agent Performance

Track and improve agent effectiveness:

go
type AgentMetrics struct {
AgentID string
TicketsHandled int
AvgResolutionTime time.Duration
AvgSatisfaction float64
FirstContactResolve float64 // %
EscalationRate float64 // %
ResponseTime time.Duration
}

func CalculateMetrics(agentID string, period time.Duration) \*AgentMetrics {
tickets := getAgentTickets(agentID, period)

    metrics := &AgentMetrics{AgentID: agentID}

    totalTime := time.Duration(0)
    totalSat := 0.0
    fcr := 0
    escalations := 0

    for _, ticket := range tickets {
        metrics.TicketsHandled++
        totalTime += ticket.ResolutionTime
        totalSat += ticket.SatisfactionScore

        if ticket.ResolvedFirstContact {
            fcr++
        }
        if ticket.Escalated {
            escalations++
        }
    }

    if metrics.TicketsHandled > 0 {
        metrics.AvgResolutionTime = totalTime / time.Duration(metrics.TicketsHandled)
        metrics.AvgSatisfaction = totalSat / float64(metrics.TicketsHandled)
        metrics.FirstContactResolve = float64(fcr) / float64(metrics.TicketsHandled)
        metrics.EscalationRate = float64(escalations) / float64(metrics.TicketsHandled)
    }

    return metrics

}

Track:

- Resolution time
- Customer satisfaction
- First contact resolution rate
- Escalation rate
- Response time
- Handle time

## COMMON PATTERNS (BOTH PATHS)

Slack Integration

go
func notifySlackUrgent(webhookURL string, ticket \*Ticket) error {
msg := &slack.WebhookMessage{
Text: "🚨 Urgent Support Ticket",
Attachments: []slack.Attachment{
{
Title: ticket.Subject,
TitleLink: ticket.URL,
Color: "danger",
Fields: []slack.AttachmentField{
{Title: "User", Value: ticket.UserID, Short: true},
{Title: "Urgency", Value: ticket.Urgency, Short: true},
{Title: "Intent", Value: ticket.Intent, Short: true},
{Title: "Sentiment", Value: fmt.Sprintf("%.2f", ticket.Sentiment.Score), Short: true},
},
Text: truncate(ticket.Message, 200),
},
},
}

    return slack.PostWebhook(webhookURL, msg)

}

Metrics Dashboard

go
var (
ticketsTotal = prometheus.NewCounterVec(
prometheus.CounterOpts{Name: "tickets_total"},
[]string{"intent", "urgency"},
)

    resolutionTime = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{Name: "resolution_time_seconds"},
        []string{"intent"},
    )

    satisfactionScore = prometheus.NewHistogram(
        prometheus.HistogramOpts{Name: "satisfaction_score"},
    )

)

Rate Limiting

go
func (rl \*RateLimiter) AllowUser(userID string) bool {
limiter := rl.getLimiter(userID)

    // Free tier: 10 requests per day
    // Pro tier: 100 requests per day
    return limiter.Allow()

}

## PRODUCTION DEPLOYMENT

High Availability:

- Multi-region deployment
- Load balancing
- Auto-scaling
- Health checks

Security:

- PII encryption at rest
- Audit logs
- Access controls
- Compliance (GDPR, CCPA)

Monitoring:

- Uptime tracking
- Error rate alerts
- Performance metrics
- Cost tracking

## SHARE YOUR WORK

Built something cool?

1. Open a PR to this file
2. Post in GitHub Discussions
3. Write a blog post
4. Present at meetups

Happy building! 🚀
