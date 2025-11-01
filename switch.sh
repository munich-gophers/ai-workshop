#!/bin/bash
# Usage: ./switch.sh path-a checkpoint-2

set -e

PATH_NAME=$1
CHECKPOINT=$2

if [ -z "$PATH_NAME" ] || [ -z "$CHECKPOINT" ]; then
    echo "Usage: ./switch.sh <path-a|path-b> <start|checkpoint-1|checkpoint-2|checkpoint-3|complete>"
    echo ""
    echo "Examples:"
    echo "  ./switch.sh path-a start          # Start Path A from scratch"
    echo "  ./switch.sh path-a checkpoint-2   # Jump to Checkpoint 2"
    echo "  ./switch.sh path-b complete       # See the full solution for Path B"
    exit 1
fi

BRANCH="${PATH_NAME}/${CHECKPOINT}"

# Check if branch exists
if ! git show-ref --verify --quiet refs/heads/${BRANCH}; then
    echo "❌ Branch '${BRANCH}' does not exist"
    echo ""
    echo "Available branches:"
    git branch | grep ${PATH_NAME} || echo "No branches found for ${PATH_NAME}"
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

# Show PROGRESS.md if it exists
PROGRESS_FILE="${PATH_NAME}-code-mentor/PROGRESS.md"
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

echo ""
