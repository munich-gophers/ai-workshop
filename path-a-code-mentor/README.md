# PATH A: AI CODE MENTOR - README

Build a GitHub-integrated PR reviewer powered by Go + Genkit + Gemini.

## OVERVIEW

A webhook service that:

- Receives GitHub webhooks when PRs are opened or updated
- Scans for secrets (API keys, passwords) before AI processing
- Analyzes code changes using Gemini 1.5 Flash
- Provides actionable feedback like a senior developer
- Returns structured JSON for automation

## PREREQUISITES

Required:

- Go 1.23+ installed
- GEMINI_API_KEY from https://aistudio.google.com/app/apikey
- Git and a code editor (VS Code recommended)

Optional (for deployment):

- Docker
- Google Cloud account

NOT needed for workshop:

- ngrok (we test locally with curl)
- GitHub webhook setup (optional post-workshop)

## QUICK START

1. Install Dependencies

go mod download

Takes about 30 seconds. Downloads Genkit SDK, Google AI plugin, and go-github.

2. Configure Environment

cp .env.example .env

Edit .env and add:
GEMINI_API_KEY=your-key-here
PORT=8080

3. Start the Server

go run cmd/server/main.go

You should see:
Server starting on port 8080
Tip: Start by implementing the health check endpoint!

4. Verify Setup

curl http://localhost:8080/health

Expected: Connection error or "not implemented" - this is your first task!

## WORKSHOP STRUCTURE

Complete 3 checkpoints in order:

CHECKPOINT 1: HTTP Server (5 minutes)
Goal: Create a working health check endpoint

What you'll do:

1. Open cmd/server/main.go
2. Find // TODO: CHECKPOINT 1
3. Add HTTP handler for /health
4. Return JSON: {"status":"healthy"}

Test:
curl http://localhost:8080/health

Stuck? See solution:
./switch.sh path-a checkpoint-1

CHECKPOINT 2: AI Integration (10 minutes)
Goal: Connect to Gemini and analyze code diffs

What you'll do:

1. Open internal/analyzer/analyzer.go
2. Find // TODO: CHECKPOINT 2 comments
3. Initialize Genkit with googleai.Init()
4. Create Gemini model instance
5. Implement Review() method
6. Parse AI responses into structured JSON

Test:
./test.sh checkpoint-2

Stuck? Compare your code:
git diff path-a/checkpoint-1..path-a/checkpoint-2

CHECKPOINT 3: Webhooks + Security (10 minutes)
Goal: Handle GitHub PR events and detect secrets

What you'll do:

1. Open internal/webhook/github.go
2. Find // TODO: CHECKPOINT 3 comments
3. Parse PullRequestEvent payloads
4. Scan code for secrets (API keys, passwords)
5. Redact secrets before sending to AI
6. Return structured analysis

Test locally (no GitHub or ngrok needed):
./test.sh checkpoint-3

Stuck? See complete solution:
./switch.sh path-a checkpoint-3

## FILE STRUCTURE

path-a-code-mentor/
├── cmd/server/main.go (Checkpoints 1,2,3)
├── internal/
│ ├── models/types.go (Complete)
│ ├── analyzer/
│ │ ├── analyzer.go (Checkpoint 2 - TODO)
│ │ └── language.go (Complete)
│ ├── handler/review.go (Complete)
│ ├── webhook/
│ │ ├── github.go (Checkpoint 3 - TODO)
│ │ └── validator.go (Complete)
│ └── security/
│ ├── detector.go (Checkpoint 3 - TODO)
│ └── patterns.go (Complete)
├── prompts/
│ ├── review-base.txt (Complete)
│ ├── review-go.txt (Complete)
│ ├── review-python.txt (Complete)
│ └── review-javascript.txt (Complete)
└── examples/
├── github-pr-payload.json (Complete)
└── simple-diff.txt (Complete)

Legend:

- Complete: No changes needed
- TODO: You implement this

## TESTING YOUR WORK

Test health endpoint (Checkpoint 1):
curl http://localhost:8080/health

Test code review API (Checkpoint 2):
curl -X POST http://localhost:8080/api/review \
 -H "Content-Type: application/json" \
 -d '{"diff":"+func add(a, b int) int {\n+ return a + b\n+}","language":"go","file_path":"math.go"}'

Test secret detection (Checkpoint 3):
curl -X POST http://localhost:8080/api/review \
 -H "Content-Type: application/json" \
 -d '{"diff":"+const API_KEY = \"sk-1234567890abcdef\"","language":"go","file_path":"config.go"}'

Test GitHub webhook handler (Checkpoint 3):
curl -X POST http://localhost:8080/webhook/github \
 -H "Content-Type: application/json" \
 -H "X-GitHub-Event: pull_request" \
 -d @examples/github-pr-payload.json

Automated testing:
./test.sh checkpoint-1
./test.sh checkpoint-2
./test.sh checkpoint-3

## ARCHITECTURE

GitHub PR Event
|
| Webhook
v
Your Service

1. Validate Signature
2. Parse PR Event
3. Extract Diff
   |
   v
   Secret Detection (Scan for API keys, etc)
   |
   v
   Redact Secrets (Replace with placeholders)
   |
   v
   Genkit -> Gemini (AI Code Review)
   |
   v
   Parse & Structure (JSON Response)
   |
   v
   Result (Issues, Secrets, Summary)

## DOCKER DEPLOYMENT

Build image:
make build
OR
docker build -t code-mentor:latest .

Run container:
make run
OR
docker run -d --name code-mentor -p 8080:8080 -e GEMINI_API_KEY="${GEMINI_API_KEY}" code-mentor:latest

View logs:
make logs
OR
docker logs -f code-mentor

## CLOUD RUN DEPLOYMENT

Deploy:
./deploy.sh

This will:

1. Build Docker image
2. Push to Google Container Registry
3. Deploy to Cloud Run
4. Output your webhook URL

After deployment you'll get a public URL like:
https://code-mentor-abc123.run.app

## KEY CONCEPTS

Genkit SDK:
Google's AI orchestration framework providing model abstraction, prompt management,
type safety, and observability. We use it to connect to Gemini, send code diffs
with prompts, and parse structured JSON responses.

go-github Library:
Official GitHub library for Go with type-safe parsing, built-in signature validation,
complete API coverage, and automatic updates. We use it for parsing webhook payloads,
validating signatures, and type-safe event handling.

Secret Detection:
Before sending code to AI, we scan for API keys (OpenAI, AWS, Google, Stripe),
authentication tokens (GitHub, Slack), hardcoded passwords, private keys (RSA, SSH),
and database connection strings. This prevents sending secrets to external APIs,
prevents accidental credential leaks, and meets data handling compliance requirements.

## TROUBLESHOOTING

Server won't start - "address already in use":
Change PORT in .env to 8081

Server won't start - "GEMINI_API_KEY not found":
Ensure .env file exists with valid key: cat .env

Tests failing - "connection refused":
Start server first: go run cmd/server/main.go

Tests failing - "not implemented" errors:
You're on START branch. Implement TODOs or switch: ./switch.sh path-a checkpoint-1

AI not responding - Empty or error responses:
Check API key: curl https://generativelanguage.googleapis.com/v1/models -H "x-goog-api-key: ${GEMINI_API_KEY}"

## WHAT YOU'LL LEARN

Go Development:

- HTTP server creation
- JSON API design
- Error handling patterns
- Package organization
- Context propagation

AI Integration:

- Genkit SDK usage
- Prompt engineering
- Structured outputs
- Token management
- Response parsing

Production Patterns:

- Webhook handling
- Signature validation
- Security scanning
- Secret redaction
- Logging and observability

DevOps:

- Docker containerization
- Multi-stage builds
- Cloud deployment
- Environment configuration

## NEXT STEPS

After completing the workshop:

1. Auto-comment on PRs - Use GitHub API to post review comments
2. Multiple file analysis - Parse complete diffs, understand cross-file changes
3. Custom rules - Add team-specific patterns, enforce coding standards
4. Learning system - Track accepted/rejected suggestions, adapt to team preferences
5. Metrics dashboard - Track review patterns, identify common issues
6. Slack integration - Notify on critical issues, send PR review summaries
7. IDE integration - Pre-commit hooks, real-time feedback

See EXTENSIONS.md for detailed implementation guides.

## RESOURCES

Documentation:

- Genkit: https://firebase.google.com/docs/genkit
- Gemini API: https://ai.google.dev/docs
- go-github: https://github.com/google/go-github
- Cloud Run: https://cloud.google.com/run/docs

Tools:

- Google AI Studio: https://aistudio.google.com/

Community:

- Genkit Discord: https://discord.gg/genkit
- Go Forum: https://forum.golangbridge.org/

## SUPPORT

During Workshop:

- Raise your hand for facilitator help
- Check PROGRESS.md for hints
- Use ./switch.sh to see solutions

After Workshop:

- Open GitHub issues for bugs
- Share implementations in Discussions
- Star the repo to save for later

## APPENDIX: REAL GITHUB WEBHOOKS (POST-WORKSHOP)

If you want to connect real GitHub webhooks after the workshop:

OPTION A: Deploy to Cloud Run (Recommended)

1. Deploy: ./deploy.sh
2. Use Cloud Run URL: https://code-mentor-abc123.run.app/webhook/github
3. Configure GitHub webhook:
   - Repo Settings -> Webhooks -> Add webhook
   - Payload URL: Your Cloud Run URL
   - Content type: application/json
   - Secret: Generate with openssl rand -hex 32
   - Events: Pull requests
4. Add secret to deployment:
   gcloud run services update code-mentor --set-env-vars GITHUB_WEBHOOK_SECRET=your-secret

OPTION B: Local Testing with ngrok

1. Install ngrok:
   macOS: brew install ngrok
   Linux: wget https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-linux-amd64.tgz

2. Start tunnel: ngrok http 8080
3. Copy HTTPS URL (e.g., https://abc123.ngrok.io)
4. Use this URL in GitHub webhook settings
5. Keep ngrok running while testing

Why ngrok? GitHub needs a public HTTPS URL to send webhooks. Your localhost:8080
is not accessible from the internet. ngrok creates a temporary public tunnel to
your local machine.

Note: ngrok URLs change each time you restart, so you'll need to update GitHub
webhook settings accordingly.

## LICENSE

MIT License - See LICENSE file

## ACKNOWLEDGMENTS

Built with Go, Genkit, Gemini, and go-github

Ready to start? Open cmd/server/main.go and look for the first TODO!
