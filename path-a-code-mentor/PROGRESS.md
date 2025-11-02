# Progress Tracker

Current Branch: CHECKPOINT-2

Status: AI integration works with Genkit v1.0+!

Next Step: Checkpoint 3 - Webhooks and Security

---

WHAT YOU COMPLETED

✅ Genkit v1.0+ initialized with unified API
✅ Google AI plugin configured
✅ Gemini 2.5 Flash model (latest version!)
✅ Prompt loading from files
✅ /api/review endpoint working
✅ Structured JSON responses with GenerateData
✅ Language-specific prompts (Go, Python, JS)

Excellent work! Your AI code reviewer uses the latest Genkit API!

---

KEY API CHANGES (Genkit v1.0+)

Old way (pre-1.0):
googleai.Init(ctx, &googleai.Config{APIKey: key})
model := googleai.Model("gemini-1.5-flash")
resp, err := model.Generate(ctx, request, nil)

New way (v1.0+):
g := genkit.Init(ctx,
genkit.WithPlugins(&googlegenai.GoogleAI{}),
genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
)
response, \_, err := genkit.GenerateData[T](ctx, g,
ai.WithPrompt("..."),
)

Benefits:

- Cleaner initialization
- Type-safe structured output
- Better error handling
- Production-ready (v1.0+ has API stability guarantees)

---

WHAT'S NEXT

Implement webhook handling and secret detection:

1. internal/webhook/github.go - Handle GitHub PR events
2. internal/security/detector.go - Scan for secrets

This completes the production-ready features.

---

TEST IT

Run the server:
go run cmd/server/main.go

Test code review:
curl -X POST http://localhost:8080/api/review \
 -H "Content-Type: application/json" \
 -d '{
"diff": "+func add(a, b int) int {\n+ return a + b\n+}",
"language": "go",
"file_path": "math.go"
}'

You should get AI-powered suggestions!

Run automated test:
./test.sh checkpoint-2

---

UNDERSTANDING GENKIT V1.0+ API

Package changes:

- googleai → googlegenai (note the "genai" suffix)
- Unified genkit.Init() replaces plugin-specific init

Model names now include provider:

- "googleai/gemini-2.5-flash" (newest, fastest)
- "googleai/gemini-2.0-flash" (multimodal)
- "googleai/gemini-1.5-flash" (legacy)

Configuration structure:
googlegenai.GeminiConfig{
Temperature: 0.3, // 0.0-2.0
MaxOutputTokens: 2000,
TopP: 0.95, // Nucleus sampling
TopK: 40, // Vocabulary filtering
}

Structured output:
GenerateData[T](ctx, g, options...)

- Type-safe JSON parsing
- Automatic schema generation from struct tags
- Fallback to text parsing if needed

---

TROUBLESHOOTING

Error: "cannot find package googleai":
Solution: Package renamed to googlegenai in v1.0+
Update import: "github.com/firebase/genkit/go/plugins/googlegenai"

Error: "undefined: googleai.Init":
Solution: v1.0+ uses unified initialization
Use: genkit.Init(ctx, genkit.WithPlugins(&googlegenai.GoogleAI{}))

Error: "GEMINI_API_KEY not set":
Solution: Set environment variable or pass to plugin
export GEMINI_API_KEY=your-key-here

Error: "Model not found":
Solution: Use provider prefix in model name
Correct: "googleai/gemini-2.5-flash"
Wrong: "gemini-2.5-flash"

---

QUICK COMMANDS

Run server:
go run cmd/server/main.go

Test with curl:
curl -X POST http://localhost:8080/api/review \
 -H "Content-Type: application/json" \
 -d @examples/simple-diff.txt

Automated test:
./test.sh checkpoint-2

Update dependencies:
go get github.com/firebase/genkit/go@latest

---

NEXT CHECKPOINT PREVIEW

In Checkpoint 3, you'll:

- Parse GitHub webhook payloads with go-github
- Validate HMAC signatures
- Scan code for secrets (API keys, passwords)
- Redact secrets before AI processing
- Return comprehensive security analysis

Time estimate: 10 minutes

Ready for the final checkpoint? Let's add webhooks and security!
