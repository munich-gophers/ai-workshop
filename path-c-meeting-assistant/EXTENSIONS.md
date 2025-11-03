# Path C: Meeting Assistant Extensions

Ideas to extend your AI meeting notes service after the workshop.

## REAL-TIME TRANSCRIPTION

1. **Integrate Google Speech-to-Text**

Convert live audio to meeting transcripts:

```go
import speech "cloud.google.com/go/speech/apiv1"
import speechpb "cloud.google.com/go/speech/apiv1/speechpb"

func (s *Service) TranscribeLive(audioStream io.Reader) (<-chan string, error) {
    client, _ := speech.NewClient(context.Background())
    stream, _ := client.StreamingRecognize(context.Background())

    // Configure streaming recognition
    req := &speechpb.StreamingRecognizeRequest{
        StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
            StreamingConfig: &speechpb.StreamingRecognitionConfig{
                Config: &speechpb.RecognitionConfig{
                    Encoding:        speechpb.RecognitionConfig_LINEAR16,
                    SampleRateHertz: 16000,
                    LanguageCode:    "en-US",
                },
                InterimResults: true,
            },
        },
    }

    stream.Send(req)

    // Stream audio chunks and return transcriptions
    transcripts := make(chan string)
    go func() {
        for {
            resp, err := stream.Recv()
            if err != nil {
                break
            }

            for _, result := range resp.Results {
                if result.IsFinal {
                    transcripts <- result.Alternatives[0].Transcript
                }
            }
        }
        close(transcripts)
    }()

    return transcripts, nil
}
```

**Setup:**
1. Enable Cloud Speech-to-Text API
2. Stream audio from web client (WebRTC)
3. Process transcripts in real-time

---

## SPEAKER IDENTIFICATION

2. **Multi-Speaker Diarization**

Identify who said what in meetings:

```go
type Speaker struct {
    ID        string
    Name      string
    Segments  []Segment
}

type Segment struct {
    StartTime time.Duration
    EndTime   time.Duration
    Text      string
}

func (s *Service) IdentifySpeakers(transcript string, audioFile string) ([]Speaker, error) {
    // Use Speech-to-Text with diarization
    config := &speechpb.RecognitionConfig{
        Encoding:                   speechpb.RecognitionConfig_LINEAR16,
        SampleRateHertz:            16000,
        LanguageCode:               "en-US",
        EnableSpeakerDiarization:   true,
        DiarizationSpeakerCount:    4, // Adjust based on meeting size
        Model:                      "video",
    }

    // Process and group by speaker
    // Return structured speaker segments
}
```

**Benefits:**
- Attribute quotes to specific people
- Generate speaker-specific summaries
- Track participation metrics

---

## AUTOMATED FOLLOW-UPS

3. **Email Generation from Action Items**

Auto-generate follow-up emails:

```go
type FollowUpEmail struct {
    To          []string
    Subject     string
    Body        string
    ActionItems []ActionItem
    Deadline    time.Time
}

func (s *Service) GenerateFollowUpEmail(analysis MeetingAnalysis) (*FollowUpEmail, error) {
    prompt := fmt.Sprintf(`
Generate a professional follow-up email based on this meeting:

Meeting: %s
Date: %s

Action Items:
%s

Key Decisions:
%s

Generate an email that:
- Summarizes the meeting
- Lists clear action items with owners
- Highlights decisions made
- Includes deadlines
- Has a professional but friendly tone
`, analysis.Title, analysis.Date,
   formatActionItems(analysis.ActionItems),
   formatDecisions(analysis.Decisions))

    response, err := genkit.Generate(ctx, s.genkit, ai.WithPrompt(prompt))
    // Parse and structure the email
    return email, nil
}
```

**Integration Ideas:**
- Send via SendGrid/Gmail API
- Add calendar invites for deadlines
- CC relevant stakeholders automatically

---

## CALENDAR INTEGRATION

4. **Auto-Create Calendar Events for Action Items**

```go
import "google.golang.org/api/calendar/v3"

func (s *Service) CreateActionItemEvents(items []ActionItem) error {
    calendarService, _ := calendar.NewService(ctx)

    for _, item := range items {
        event := &calendar.Event{
            Summary:     fmt.Sprintf("TODO: %s", item.Description),
            Description: fmt.Sprintf("Action item from meeting: %s\nOwner: %s",
                                   item.MeetingTitle, item.Owner),
            Start: &calendar.EventDateTime{
                Date: item.Deadline.Format("2006-01-02"),
            },
            End: &calendar.EventDateTime{
                Date: item.Deadline.Format("2006-01-02"),
            },
            Attendees: []*calendar.EventAttendee{
                {Email: item.OwnerEmail},
            },
            Reminders: &calendar.EventReminders{
                UseDefault: false,
                Overrides: []*calendar.EventReminder{
                    {Method: "email", Minutes: 24 * 60},  // 1 day before
                    {Method: "popup", Minutes: 60},        // 1 hour before
                },
            },
        }

        _, err := calendarService.Events.Insert("primary", event).Do()
        if err != nil {
            return err
        }
    }

    return nil
}
```

---

## MEETING ANALYTICS

5. **Dashboard with Meeting Insights**

Track meeting effectiveness over time:

```go
type MeetingMetrics struct {
    TotalMeetings       int
    AvgDuration         time.Duration
    AvgActionItems      float64
    CompletionRate      float64
    TopParticipants     []string
    MeetingsByType      map[string]int
    ActionItemTrends    []DataPoint
}

func (s *Service) GenerateAnalytics(startDate, endDate time.Time) (*MeetingMetrics, error) {
    // Query meetings from database
    meetings := s.db.GetMeetingsByDateRange(startDate, endDate)

    metrics := &MeetingMetrics{
        TotalMeetings: len(meetings),
    }

    // Calculate metrics
    for _, meeting := range meetings {
        metrics.AvgDuration += meeting.Duration
        metrics.AvgActionItems += float64(len(meeting.ActionItems))

        // Track completion
        completed := countCompleted(meeting.ActionItems)
        metrics.CompletionRate += float64(completed) / float64(len(meeting.ActionItems))
    }

    // Average calculations
    if len(meetings) > 0 {
        metrics.AvgDuration /= time.Duration(len(meetings))
        metrics.AvgActionItems /= float64(len(meetings))
        metrics.CompletionRate = (metrics.CompletionRate / float64(len(meetings))) * 100
    }

    return metrics, nil
}
```

**Dashboard Features:**
- Meeting frequency trends
- Action item completion rates
- Most productive meeting types
- Participant engagement scores
- Time-to-completion for tasks

---

## MEETING TEMPLATES

6. **Pre-defined Templates for Common Meeting Types**

```go
type MeetingTemplate struct {
    Type            string
    DefaultAgenda   []string
    RequiredTopics  []string
    SuggestedDuration time.Duration
    Participants    []string
}

var templates = map[string]MeetingTemplate{
    "standup": {
        Type: "Daily Standup",
        DefaultAgenda: []string{
            "What did you do yesterday?",
            "What will you do today?",
            "Any blockers?",
        },
        SuggestedDuration: 15 * time.Minute,
    },
    "sprint-planning": {
        Type: "Sprint Planning",
        DefaultAgenda: []string{
            "Review sprint goal",
            "Review backlog",
            "Estimate stories",
            "Commit to sprint",
        },
        SuggestedDuration: 2 * time.Hour,
    },
    "retrospective": {
        Type: "Sprint Retrospective",
        DefaultAgenda: []string{
            "What went well?",
            "What could be improved?",
            "Action items for next sprint",
        },
        SuggestedDuration: 1 * time.Hour,
    },
}

func (s *Service) ApplyTemplate(templateType string, transcript string) (*MeetingAnalysis, error) {
    template := templates[templateType]

    prompt := fmt.Sprintf(`
Analyze this %s meeting transcript using this template:

Agenda:
%s

Transcript:
%s

Extract information according to the template structure.
`, template.Type, strings.Join(template.DefaultAgenda, "\n"), transcript)

    // Generate structured analysis based on template
    return analysis, nil
}
```

---

## INTEGRATION WITH PROJECT TOOLS

7. **Sync Action Items to Jira**

```go
import "github.com/felixgeelhaar/jirasdk"

func (s *Service) CreateJiraIssues(items []ActionItem) error {
    // Initialize Jira client
    tp := jirasdk.BasicAuthTransport{
        Username: s.jiraEmail,
        Password: s.jiraAPIToken,
    }

    jiraClient, err := jirasdk.NewClient(tp.Client(), s.jiraURL)
    if err != nil {
        return fmt.Errorf("failed to create jira client: %w", err)
    }

    for _, item := range items {
        issue := &jirasdk.Issue{
            Fields: &jirasdk.IssueFields{
                Project: jirasdk.Project{
                    Key: s.projectKey,
                },
                Summary:     item.Description,
                Description: fmt.Sprintf("Action item from meeting: %s\n\nDecision context:\n%s\n\nOwner: %s",
                                       item.MeetingTitle, item.Context, item.Owner),
                Type: jirasdk.IssueType{
                    Name: "Task",
                },
                Assignee: &jirasdk.User{
                    Name: item.Owner,
                },
                Priority: &jirasdk.Priority{
                    Name: item.Priority,
                },
            },
        }

        // Set due date if available
        if !item.Deadline.IsZero() {
            dueDate := jirasdk.Date(item.Deadline)
            issue.Fields.Duedate = &dueDate
        }

        createdIssue, resp, err := jiraClient.Issue.Create(context.Background(), issue)
        if err != nil {
            return fmt.Errorf("failed to create issue: %w (status: %d)", err, resp.StatusCode)
        }

        // Store Jira issue key with action item
        item.ExternalID = createdIssue.Key
        log.Printf("Created Jira issue: %s for action item: %s", createdIssue.Key, item.Description)
    }

    return nil
}

// Optional: Update existing Jira issues
func (s *Service) UpdateJiraIssue(item ActionItem) error {
    tp := jirasdk.BasicAuthTransport{
        Username: s.jiraEmail,
        Password: s.jiraAPIToken,
    }

    jiraClient, _ := jirasdk.NewClient(tp.Client(), s.jiraURL)

    issue := &jirasdk.Issue{
        Fields: &jirasdk.IssueFields{
            Description: fmt.Sprintf("Updated from meeting notes\n\n%s", item.Context),
        },
    }

    _, resp, err := jiraClient.Issue.Update(context.Background(), item.ExternalID, issue)
    if err != nil {
        return fmt.Errorf("failed to update issue %s: %w (status: %d)",
                         item.ExternalID, err, resp.StatusCode)
    }

    return nil
}
```

**Setup:**
```bash
go get github.com/felixgeelhaar/jirasdk
```

**Configuration:**
```go
type JiraConfig struct {
    URL        string // e.g., "https://yourcompany.atlassian.net"
    Email      string
    APIToken   string // Generate from Atlassian account settings
    ProjectKey string // e.g., "PROJ"
}
```

---

## MULTI-LANGUAGE SUPPORT

8. **Translate Meetings to Multiple Languages**

```go
import "cloud.google.com/go/translate"

func (s *Service) TranslateMeeting(analysis *MeetingAnalysis, targetLang string) (*MeetingAnalysis, error) {
    client, _ := translate.NewClient(context.Background())

    // Translate summary
    summaryTranslations, _ := client.Translate(
        context.Background(),
        []string{analysis.Summary},
        language.Make(targetLang),
        nil,
    )

    translatedAnalysis := *analysis
    translatedAnalysis.Summary = summaryTranslations[0].Text

    // Translate action items
    for i, item := range translatedAnalysis.ActionItems {
        descriptions, _ := client.Translate(
            context.Background(),
            []string{item.Description},
            language.Make(targetLang),
            nil,
        )
        translatedAnalysis.ActionItems[i].Description = descriptions[0].Text
    }

    return &translatedAnalysis, nil
}
```

**Use Cases:**
- Global distributed teams
- Multilingual meetings
- Accessibility for non-native speakers

---

## MEETING SEARCH & RETRIEVAL

9. **Semantic Search Across All Meetings**

```go
import "github.com/firebase/genkit/go/ai"

func (s *Service) SearchMeetings(query string) ([]MeetingAnalysis, error) {
    // Generate embedding for query
    queryEmbedding, _ := ai.Embed(ctx, s.genkit, ai.WithText(query))

    // Search vector database (e.g., Pinecone, Weaviate, or Firestore)
    results, _ := s.vectorDB.Search(queryEmbedding, 10)

    // Retrieve full meeting analyses
    var meetings []MeetingAnalysis
    for _, result := range results {
        meeting, _ := s.db.GetMeeting(result.ID)
        meetings = append(meetings, meeting)
    }

    return meetings, nil
}
```

**Example Queries:**
- "What decisions did we make about the new API?"
- "All action items assigned to Sarah this month"
- "Meetings where we discussed budget"

---

## SENTIMENT ANALYSIS

10. **Track Meeting Mood and Engagement**

```go
type SentimentScore struct {
    Overall     float64 // -1 to 1
    BySegment   []float64
    Concerns    []string
    Positive    []string
}

func (s *Service) AnalyzeSentiment(transcript string) (*SentimentScore, error) {
    prompt := `
Analyze the sentiment of this meeting transcript.

Transcript:
` + transcript + `

Provide:
1. Overall sentiment score (-1 to 1)
2. Segment-by-segment sentiment
3. Key concerns or negative points
4. Positive highlights
5. Team morale indicators
`

    response, _ := genkit.Generate(ctx, s.genkit, ai.WithPrompt(prompt))
    // Parse and return sentiment analysis
    return score, nil
}
```

**Insights:**
- Detect team morale issues
- Identify contentious topics
- Track meeting energy over time

---

## RESOURCES

**APIs & Services:**
- [Google Cloud Speech-to-Text](https://cloud.google.com/speech-to-text/docs)
- [Google Calendar API](https://developers.google.com/calendar/api)
- [Cloud Translation API](https://cloud.google.com/translate/docs)
- [SendGrid Email API](https://docs.sendgrid.com/)
- [Jira REST API](https://developer.atlassian.com/cloud/jira/platform/rest/v3/)
- [Linear API](https://developers.linear.app/)

**Vector Databases:**
- [Pinecone](https://www.pinecone.io/)
- [Weaviate](https://weaviate.io/)
- [Firestore Vector Search](https://firebase.google.com/docs/firestore/vector-search)

**Libraries:**
- [felixgeelhaar/jirasdk](https://github.com/felixgeelhaar/jirasdk) - Jira client
- [go-github](https://github.com/google/go-github) - GitHub API client
- [webhooks](https://github.com/go-playground/webhooks) - Webhook handling

---

## NEXT STEPS

Pick one extension that interests you:
1. **Quick win** (1-2 hours): Email generation, meeting templates
2. **Medium** (half day): Speaker identification, calendar integration
3. **Advanced** (full day): Real-time transcription, semantic search, Jira integration

Remember: Start small, test thoroughly, and iterate based on user feedback!
