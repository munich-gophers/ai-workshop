# PROGRESS TRACKER

Current Branch: START

Next Step: Checkpoint 1 - Basic Server

## WHAT YOU NEED TO DO

Implement a health check endpoint in cmd/server/main.go

Step 1: Open the file

code cmd/server/main.go

OR

vim cmd/server/main.go

Step 2: Find the TODO comment

Look for:
// TODO: CHECKPOINT 1 - Add health check endpoint

Step 3: Add the health check handler

Replace the TODO with this code:

mux.HandleFunc("/health", func(w http.ResponseWriter, r \*http.Request) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
w.Write([]byte(`{"status":"healthy","service":"code-mentor","version":"1.0.0"}`))
})

## TEST IT

Start the server:

go run cmd/server/main.go

You should see:
Server starting on port 8080

Test the endpoint in another terminal:

curl http://localhost:8080/health

Expected output:

{"status":"healthy","service":"code-mentor","version":"1.0.0"}

## NEXT STEPS

Once your health check works, move to Checkpoint 1:

./switch.sh path-a checkpoint-1

This shows the complete solution and sets you up for Checkpoint 2.

Compare your solution:

git diff path-a/start..path-a/checkpoint-1

This shows exactly what changes between START and CHECKPOINT-1.

See all checkpoints:

git branch | grep path-a

## WORKSHOP TIMELINE

Checkpoint Focus Time Status
START Setup and health endpoint 5 min YOU ARE HERE
1 HTTP server basics - Next
2 Genkit + Gemini integration 10 min Later
3 Webhooks + security 10 min Later
Complete Deployment + extensions 5 min Final

Total workshop time: 30 minutes of coding + 15 minutes buffer

## KEY FILES YOU'LL EDIT

1. cmd/server/main.go (Checkpoint 1, 2, 3)
   - Add HTTP handlers
   - Initialize analyzer
   - Wire up webhooks

2. internal/analyzer/analyzer.go (Checkpoint 2)
   - Initialize Genkit
   - Implement AI review logic
   - Build prompts

3. internal/webhook/github.go (Checkpoint 3)
   - Parse webhook events
   - Validate signatures
   - Orchestrate review pipeline

4. internal/security/detector.go (Checkpoint 3)
   - Scan for secrets
   - Redact sensitive data

## GETTING HELP

Stuck on implementation?

1. Check the TODO comments - they have hints
2. Look at examples in examples/ directory
3. Switch to next checkpoint to see solution
4. Ask the facilitator

Common issues:

Import errors:
go mod tidy

Port already in use:
Change PORT in .env to 8081

Can't find prompts:
Make sure you're running from project root
pwd should end in path-a-code-mentor

## UNDERSTANDING THE STRUCTURE

path-a-code-mentor/
├── cmd/server/main.go <- You start here (Checkpoint 1)
├── internal/
│ ├── analyzer/ <- Checkpoint 2
│ │ └── analyzer.go
│ ├── webhook/ <- Checkpoint 3
│ │ └── github.go
│ ├── security/ <- Checkpoint 3
│ │ └── detector.go
│ └── models/
│ └── types.go <- Already done
└── prompts/ <- Already done
└── \*.txt

Files marked "Already done" are complete - you don't need to edit them.
Files with checkpoint numbers are where you'll add code.

## TIPS FOR SUCCESS

1. Read the TODO comments carefully - they contain step-by-step hints
2. Test after each checkpoint - use ./test.sh checkpoint-X
3. Don't worry about perfection - the goal is learning, not perfect code
4. Use the switch script - if stuck for more than 3 minutes, see the solution
5. Ask questions early - don't suffer in silence!

## WHAT YOU'RE LEARNING

Checkpoint 1: HTTP Basics

- Go HTTP handlers
- JSON responses
- Server setup

Checkpoint 2: AI Integration

- Genkit SDK
- Calling Gemini API
- Prompt engineering
- JSON parsing

Checkpoint 3: Production Patterns

- Webhook handling
- Security scanning
- Event processing
- Error handling

## QUICK COMMANDS REFERENCE

Run server:
go run cmd/server/main.go

Test health:
curl http://localhost:8080/health

Run checkpoint tests:
./test.sh checkpoint-1
./test.sh checkpoint-2
./test.sh checkpoint-3

Switch to next checkpoint:
./switch.sh path-a checkpoint-1
./switch.sh path-a checkpoint-2
./switch.sh path-a checkpoint-3

See solution differences:
git diff path-a/start..path-a/checkpoint-1
git diff path-a/checkpoint-1..path-a/checkpoint-2
git diff path-a/checkpoint-2..path-a/checkpoint-3

Deploy (after checkpoint 3):
./deploy.sh

## READY? LET'S GO!

Open cmd/server/main.go and look for the first TODO comment.

You've got this!
