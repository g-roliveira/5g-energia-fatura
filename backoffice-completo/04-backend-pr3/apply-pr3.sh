#!/bin/bash
set -e

# PR 3/N — PDF Generation, Worker Pool, Bulk Actions, LISTEN/NOTIFY
# Apply script for backoffice billing module

BUNDLE_FILE="backoffice-billing-pr3.bundle"
REPO_DIR="${1:-.}"

echo "======================================"
echo "PR 3/N — Backoffice Billing Module"
echo "======================================"
echo ""

# Check if in a git repo
if [ ! -d "$REPO_DIR/.git" ]; then
    echo "❌ Error: $REPO_DIR is not a git repository"
    echo "Usage: ./apply-pr3.sh [path-to-5g-energia-fatura]"
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
echo "🔀 Merging pr3-billing-pdf-workers branch..."
git merge pr3-billing-pdf-workers --no-edit

echo ""
echo "✅ PR 3 applied successfully!"
echo ""
echo "📋 Next steps:"
echo ""
echo "1. Install dependencies:"
echo "   cd services/backend-go"
echo "   go mod tidy"
echo ""
echo "2. Install chromium:"
echo "   apt-get install chromium-browser"
echo "   # ou"
echo "   apk add chromium"
echo ""
echo "3. Create directories:"
echo "   mkdir -p /var/lib/backoffice/pdfs"
echo "   mkdir -p /opt/backoffice/templates"
echo ""
echo "4. Integrate into server.go:"
echo "   See services/backend-go/internal/app/BILLING_INTEGRATION_PR3.md"
echo ""
echo "5. Test bulk sync:"
echo "   curl -X POST http://localhost:8080/v1/billing/cycles/{id}/bulk \\"
echo "     -d '{\"action\": \"sync\"}'"
echo ""
echo "6. Test SSE LISTEN/NOTIFY:"
echo "   curl -N http://localhost:8080/v1/billing/events/cycles/{id}"
echo ""
echo "📄 Read PR3_README.md for complete documentation"
