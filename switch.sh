#!/bin/bash
# Usage: ./switch.sh path-a checkpoint-2

set -e

PATH_NAME=$1
CHECKPOINT=$2

# Function to display progress information for the current path/checkpoint
show_progress() {
    local PATH_NAME=$1
    local CHECKPOINT=$2

    # Map path name to directory name
    case $PATH_NAME in
        path-a)
            DIR_NAME="path-a-code-mentor"
            ;;
        path-b)
            DIR_NAME="path-b-support-agent"
            ;;
        path-c)
            DIR_NAME="path-c-meeting-assistant"
            ;;
        path-d)
            DIR_NAME="path-d-content-moderator"
            ;;
        *)
            DIR_NAME="${PATH_NAME}"
            ;;
    esac

    # Show PROGRESS.md if it exists
    PROGRESS_FILE="${DIR_NAME}/PROGRESS.md"
    if [ -f "${PROGRESS_FILE}" ]; then
        echo "=========================================="
        head -n 20 "${PROGRESS_FILE}"
        echo "=========================================="
        echo ""
        echo "📖 Full progress details: cat ${PROGRESS_FILE}"
    else
        # Generic message if no PROGRESS.md
        case $CHECKPOINT in
            start)
                echo "📍 Starting point - scaffold with TODOs"
                echo "👉 Look for // TODO: comments in the code"
                ;;
            checkpoint-1)
                echo "📍 Checkpoint 1 - Health endpoint works"
                echo "👉 Next: ./switch.sh ${PATH_NAME} checkpoint-2"
                ;;
            checkpoint-2)
                echo "📍 Checkpoint 2 - AI integration works"
                echo "👉 Next: ./switch.sh ${PATH_NAME} checkpoint-3"
                ;;
            checkpoint-3)
                echo "📍 Checkpoint 3 - Full integration works"
                echo "👉 Next: ./switch.sh ${PATH_NAME} complete"
                ;;
            complete)
                echo "📍 Complete solution with bonus features"
                echo "🎉 Great job! Check EXTENSIONS.md for ideas"
                ;;
        esac
    fi
}

if [ -z "$PATH_NAME" ] || [ -z "$CHECKPOINT" ]; then
    echo "Usage: ./switch.sh <path-a|path-b|path-c|path-d> <start|checkpoint-1|checkpoint-2|checkpoint-3>"
    echo ""
    echo "Examples:"
    echo "  ./switch.sh path-a start          # Start Path A (Code Mentor) from scratch"
    echo "  ./switch.sh path-a checkpoint-2   # Jump to Checkpoint 2"
    echo "  ./switch.sh path-b checkpoint-3   # Complete Path B (Support Agent)"
    echo "  ./switch.sh path-c start          # Start Path C (Meeting Assistant)"
    echo "  ./switch.sh path-d checkpoint-1   # Path D (Content Moderator) Checkpoint 1"
    exit 1
fi

BRANCH="${PATH_NAME}/${CHECKPOINT}"

# Fetch latest from remote to ensure we have all branches
echo "🔄 Fetching latest branches from remote..."
git fetch origin --quiet

# Check if branch exists locally
if git show-ref --verify --quiet refs/heads/${BRANCH}; then
    echo "✅ Local branch '${BRANCH}' found"
# Check if branch exists on remote
elif git show-ref --verify --quiet refs/remotes/origin/${BRANCH}; then
    echo "✅ Remote branch 'origin/${BRANCH}' found, creating local tracking branch..."
    git checkout --track origin/${BRANCH}
    echo ""
    echo "✅ You're now on ${BRANCH}"
    echo ""
    show_progress "$PATH_NAME" "$CHECKPOINT"
    echo ""
    exit 0
else
    echo "❌ Branch '${BRANCH}' does not exist locally or remotely"
    echo ""
    echo "Available branches:"
    git branch -a | grep ${PATH_NAME} || echo "No branches found for ${PATH_NAME}"
    exit 1
fi

# Stash any changes
if [[ -n $(git status -s) ]]; then
    echo "💾 Stashing your changes..."
    git stash push -m "Auto-stash before switching to ${BRANCH}"
    echo ""
fi

# Switch branch
echo "🔄 Switching to ${BRANCH}..."
git checkout ${BRANCH}

echo ""
echo "✅ You're now on ${BRANCH}"
echo ""
show_progress "$PATH_NAME" "$CHECKPOINT"
echo ""
