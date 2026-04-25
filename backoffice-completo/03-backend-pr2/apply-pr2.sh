#!/bin/bash
set -e

# PR 2/N — Billing Cycles, Adjustments, and SSE
# Apply script for backoffice billing module

BUNDLE_FILE="backoffice-billing-pr2.bundle"
REPO_DIR="${1:-.}"

echo "======================================"
echo "PR 2/N — Backoffice Billing Module"
echo "======================================"
echo ""

# Check if in a git repo
if [ ! -d "$REPO_DIR/.git" ]; then
    echo "❌ Error: $REPO_DIR is not a git repository"
    echo "Usage: ./apply.sh [path-to-5g-energia-fatura]"
    exit 1
fi

cd "$REPO_DIR"

# Check if bundle file exists
if [ ! -f "$BUNDLE_FILE" ]; then
    echo "❌ Error: $BUNDLE_FILE not found in current directory"
    echo "Make sure to extract the tar.gz first"
    exit 1
fi

echo "📦 Extracting bundle..."
git bundle unbundle "$BUNDLE_FILE"

echo ""
echo "🔀 Merging pr2-billing-cycles branch..."
git merge pr2-billing-cycles --no-edit

echo ""
echo "✅ PR 2 applied successfully!"
echo ""
echo "📋 Next steps:"
echo ""
echo "1. Run migration:"
echo "   cd services/backend-go"
echo "   migrate -path migrations -database \"\$BACKOFFICE_PG_URL\" up"
echo ""
echo "2. Verify migration:"
echo "   psql \"\$BACKOFFICE_PG_URL\" -c \"SELECT table_name FROM information_schema.tables WHERE table_schema = 'core' AND table_name = 'notification';\""
echo ""
echo "3. Integrate into server.go:"
echo "   See services/backend-go/internal/app/BILLING_INTEGRATION_PR2.md"
echo ""
echo "4. Test SSE endpoint:"
echo "   curl -N http://localhost:8080/v1/billing/events/cycles/{cycle_id}"
echo ""
echo "📄 Read PR2_README.md for complete documentation"
