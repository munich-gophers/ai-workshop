# Ship an AI Assistant in 45 Minutes

Build and deploy a production-ready AI service using Go + Genkit + Gemini.

Choose your adventure - four complete learning paths:

## WHY THIS WORKSHOP?

Learn production-ready AI development through hands-on building. Each path teaches:
- **Real-world patterns** you'll use in production
- **Go + Genkit** for type-safe AI applications
- **Google Gemini** for powerful AI capabilities
- **Incremental learning** through 3 guided checkpoints

**All paths are standalone** - pick one that matches your interests or complete multiple to master different AI patterns!

## CHOOSE YOUR PATH

Path A: AI Code Mentor
Build a GitHub/GitLab PR reviewer that:

- Receives webhook events from your repo
- Analyzes code changes with Gemini
- Suggests improvements like a senior developer
- (Optional) Learns your team's style over time

Best for: DevOps engineers, platform teams, tech leads

Start Path A: See path-a-code-mentor/README.md

---

Path B: AI Support Agent
Build a customer support triage service that:

- Redacts PII before processing
- Classifies intent and urgency
- Summarizes user messages
- Returns structured JSON for automation

Best for: Support engineers, SRE teams, product builders

Start Path B: See path-b-support-agent/README.md

---

Path C: AI Meeting Notes Assistant
Build a meeting transcription and summarization service that:

- Extracts structured information from meeting transcripts
- Identifies action items and decisions
- Generates executive summaries
- Returns organized JSON with key insights

Best for: Product managers, team leads, documentation enthusiasts

Start Path C: `./switch.sh path-c start`

---

Path D: AI Content Moderator & Sentiment Analyzer
Build a content moderation and analysis service that:

- Performs pattern-based sentiment analysis
- Classifies content across multiple safety categories
- Provides comprehensive risk assessment
- Recommends automated actions with confidence scores

Best for: Community managers, social platform builders, content teams

Start Path D: `./switch.sh path-d start`

---

## PREREQUISITES

Complete these BEFORE the workshop:

**Required (15 min setup):**

1. **Go 1.23+** installed ([go.dev](https://go.dev))
2. **Gemini API key** from [AI Studio](https://aistudio.google.com/app/apikey)
3. **Git** and a code editor (VS Code recommended)

**Optional (for deployment & integrations):**

4. **gcloud CLI** - For deploying to Google Cloud Run ([install guide](https://cloud.google.com/sdk/docs/install))
   - Run `gcloud auth login` after installation
   - Requires Google Cloud account with billing enabled
5. **ngrok** - For local testing with webhooks/integrations ([ngrok.com](https://ngrok.com))
   - Useful for Path A (Code Mentor) webhook testing
   - Exposes local server to the internet

**Verify Setup:**
```bash
cd prerequisites
./verify.sh
```

Expected output:
```
✅ Go 1.23.2 detected
✅ GEMINI_API_KEY set
✅ Ready to ship!
```

---

## WORKSHOP STRUCTURE

Phase 1: Choose Your Path (2 min)
Pick Path A, B, C, or D and switch to that branch:
```bash
./switch.sh <path-a|path-b|path-c|path-d> start
```

Phase 2: Build (25 min)
Each path has 3 checkpoints:

- Checkpoint 1: Foundation & basics (5 min)
- Checkpoint 2: AI integration (10 min)
- Checkpoint 3: Full feature set (10 min)

Navigate between checkpoints:
```bash
./switch.sh <path-name> checkpoint-1
./switch.sh <path-name> checkpoint-2
./switch.sh <path-name> checkpoint-3
```

Phase 3: Deploy (Optional, 10 min)
Deploy your service to production

---

## QUICK START

1. **Clone the repository:**
   ```bash
   git clone https://github.com/munich-gophers/ai-workshop
   cd ai-workshop
   ```

2. **Verify your setup:**
   ```bash
   cd prerequisites
   ./verify.sh
   cd ..
   ```

3. **Choose your path and start building:**
   ```bash
   # Example: Start Path A (Code Mentor)
   ./switch.sh path-a start

   # Or start any other path
   # ./switch.sh path-b start  # Support Agent
   # ./switch.sh path-c start  # Meeting Assistant
   # ./switch.sh path-d start  # Content Moderator
   ```

4. **Follow the checkpoints:**
   - Each path has TODOs marked in the code
   - Complete each checkpoint before moving to the next
   - Use `./switch.sh <path-name> checkpoint-X` to jump to solutions

---

## WHAT YOU'LL LEARN

Each path teaches different AI patterns and techniques:

| Path | Core AI Pattern | Key Skills |
|------|----------------|------------|
| **Path A: Code Mentor** | Code Analysis & Generation | Webhook processing, structured code review, context-aware suggestions |
| **Path B: Support Agent** | Text Classification & PII Handling | Data redaction, intent classification, structured extraction |
| **Path C: Meeting Assistant** | Information Extraction | Transcript parsing, action item detection, executive summarization |
| **Path D: Content Moderator** | Multi-label Classification | Sentiment analysis, content safety, confidence scoring, decision automation |

**Common to all paths:**
- Setting up Go + Genkit projects
- Integrating Google Gemini AI
- Prompt engineering techniques
- Structured JSON responses
- Error handling & validation
- Testing AI applications

---

## RESOURCES

**AI & Development:**
- [Genkit Documentation](https://firebase.google.com/docs/genkit) - Firebase Genkit framework
- [Gemini API](https://ai.google.dev/docs) - Google's AI model documentation
- [Go Documentation](https://go.dev/doc/) - Official Go language docs

**Deployment & Testing:**
- [Google Cloud Run](https://cloud.google.com/run/docs) - Serverless container deployment
- [ngrok](https://ngrok.com/docs) - Local tunnel for webhook testing
- [gcloud CLI](https://cloud.google.com/sdk/gcloud) - Google Cloud command-line tool

---

## LICENSE

MIT License - See LICENSE file

Ready? Pick your path and start building!
