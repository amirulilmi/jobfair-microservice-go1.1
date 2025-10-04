#!/bin/bash

# Make all scripts executable
chmod +x scripts/apply-bug-fixes.sh
chmod +x scripts/test-bug-fixes.sh
chmod +x scripts/populate-company-mappings.sh
chmod +x scripts/quick-fix-company-mapping.sh

echo "✅ All scripts are now executable!"
echo ""
echo "Available scripts:"
echo "  ./scripts/apply-bug-fixes.sh           - Apply bug fixes"
echo "  ./scripts/test-bug-fixes.sh [TOKEN]    - Test bug fixes"
echo "  ./scripts/populate-company-mappings.sh - Populate company mappings"
echo "  ./scripts/quick-fix-company-mapping.sh - Quick fix company mapping"
