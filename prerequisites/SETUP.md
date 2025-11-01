# Workshop Prerequisites Setup

Complete these steps BEFORE the workshop to save time.

1. INSTALL GO

---

Download and install Go 1.23 or later:

- Visit: https://go.dev/doc/install
- Follow instructions for your OS

Verify installation:
go version

Should show: go version go1.23.x

2. INSTALL GCLOUD CLI

---

Download and install:

- Visit: https://cloud.google.com/sdk/docs/install
- Follow instructions for your OS

Authenticate:
gcloud auth login
gcloud config set project YOUR_PROJECT_ID

Verify:
gcloud config get-value project

Should show your project ID

3. GET GEMINI API KEY

---

Steps:

1. Visit: https://aistudio.google.com/app/apikey
2. Click "Create API Key"
3. Copy the key

Set as environment variable:

macOS/Linux:
echo 'export GEMINI_API_KEY=your-key-here' >> ~/.bashrc
source ~/.bashrc

Windows PowerShell:
[System.Environment]::SetEnvironmentVariable('GEMINI_API_KEY', 'your-key-here', 'User')

4. INSTALL GIT

---

If not already installed:

- Visit: https://git-scm.com/downloads
- Follow instructions for your OS

5. INSTALL CODE EDITOR

---

Recommended: VS Code

- Visit: https://code.visualstudio.com/
- Install Go extension

6. VERIFY EVERYTHING

---

Run the verification script:

cd prerequisites
./verify.sh

You should see all green checkmarks

## TROUBLESHOOTING

Go not found:

- Make sure Go bin directory is in your PATH
- Restart your terminal

gcloud not authenticated:

- Run: gcloud auth login
- Follow browser prompts

API key not set:

- Check: echo $GEMINI_API_KEY
- Should output your key
- Restart terminal if just set

## OPTIONAL (helpful)

Docker:
macOS: brew install docker
Windows/Linux: https://docs.docker.com/get-docker/

ngrok (for local webhook testing):
macOS: brew install ngrok
Others: https://ngrok.com/download

## TIME ESTIMATE

- If you have nothing: 30-45 minutes
- If you have Go installed: 15-20 minutes
- If you have everything except API key: 5 minutes

Questions? Bring them to the workshop or email the facilitator.
