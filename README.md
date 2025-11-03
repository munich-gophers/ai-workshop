# Ship an AI Assistant in 45 Minutes

Build and deploy a production-ready AI service using Go + Genkit + Gemini.

Choose your adventure - four complete learning paths:

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

Start Path C: See path-c-meeting-assistant/README.md

---

Path D: AI Content Moderator & Sentiment Analyzer
Build a content moderation and analysis service that:

- Performs pattern-based sentiment analysis
- Classifies content across multiple safety categories
- Provides comprehensive risk assessment
- Recommends automated actions with confidence scores

Best for: Community managers, social platform builders, content teams

Start Path D: See path-d-content-moderator/README.md

---

## PREREQUISITES

Complete these BEFORE the workshop:

Required (30 min setup):

1. Go 1.23+ installed (go.dev)
2. Google Cloud account with billing enabled
3. gcloud CLI authenticated: gcloud auth login
4. Gemini API key from: https://aistudio.google.com/app/apikey
5. Git and a code editor (VS Code recommended)

Verify Setup:
cd prerequisites
./verify.sh

Expected output:
✅ Go 1.23.2 detected
✅ gcloud authenticated
✅ GEMINI_API_KEY set
✅ Ready to ship!

---

## WORKSHOP STRUCTURE

Phase 1: Choose Your Path (2 min)
Pick Path A, B, C, or D and navigate to that directory

Phase 2: Build (25 min)
Each path has 3 checkpoints:

- Checkpoint 1: Basic server (5 min)
- Checkpoint 2: AI integration (10 min)
- Checkpoint 3: Full feature (10 min)

Test at each checkpoint: ./test.sh checkpoint-X
Stuck? See solution: ./switch.sh <path-a|path-b|path-c|path-d> checkpoint-X

Phase 3: Deploy (Optional, 10 min)
./deploy.sh

---

## QUICK START

1. Clone:
   git clone https://github.com/munich-gophers/ai-workshop
   cd ai-workshop

2. Verify:
   cd prerequisites
   ./verify.sh

3. Choose path:
   cd ../path-a-code-mentor

4. Follow README in that directory

---

## RESOURCES

- Genkit: https://firebase.google.com/docs/genkit
- Gemini API: https://ai.google.dev/docs
- Cloud Run: https://cloud.google.com/run/docs

---

## LICENSE

MIT License - See LICENSE file

Ready? Pick your path and start building!
