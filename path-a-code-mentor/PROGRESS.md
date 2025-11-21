# Progress Tracker

Current Branch: START

Status: Ready to begin!

Next Step: Implement the health endpoint and basic server

---

STARTING POINT

This is the scaffold with TODOs. Your task is to:

1. Implement the health check endpoint
2. Set up the HTTP server
3. Add basic error handling
4. Test that it works

Look for `// TODO:` comments in the code to guide you.

---

QUICK START

1. Copy .env.example to .env:
   ```bash
   cp .env.example .env
   ```

2. Add your Gemini API key to .env:
   ```
   GEMINI_API_KEY=your_key_here
   ```

3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

4. Test the health endpoint:
   ```bash
   curl http://localhost:8080/health
   ```

---

WHAT YOU'LL BUILD

In this workshop path, you'll create an AI-powered code review service that:

- Receives webhook events from GitHub/GitLab
- Analyzes code changes using Gemini AI
- Provides intelligent code review suggestions
- Returns structured JSON responses

---

CHECKPOINTS

- **START** (current): Scaffold with TODOs
- **Checkpoint 1**: Health endpoint working
- **Checkpoint 2**: AI integration with Genkit
- **Checkpoint 3**: Full webhook integration

---

MOVE TO NEXT CHECKPOINT

When you're ready to see the solution for the health endpoint:

```bash
./switch.sh path-a checkpoint-1
```

Or if you want to see what changed:

```bash
git diff path-a/start..path-a/checkpoint-1
```

---

TIPS

- Read the README.md for detailed instructions
- Check EXTENSIONS.md for bonus features
- Use the TODO comments as your guide
- Don't hesitate to ask for help!

Good luck! 🚀
