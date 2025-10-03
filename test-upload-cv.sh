#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}===== JobFair CV Upload Test =====${NC}\n"

# Check if token is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: JWT token required${NC}"
    echo "Usage: ./test-upload-cv.sh <JWT_TOKEN> [file_path]"
    echo ""
    echo "Example:"
    echo "  ./test-upload-cv.sh eyJhbGc... test.pdf"
    echo ""
    exit 1
fi

TOKEN="$1"

# Use provided file or create a test PDF
if [ -z "$2" ]; then
    echo -e "${YELLOW}No file provided, creating test CV file...${NC}"
    TEST_FILE="test_cv.pdf"
    
    # Create a simple test PDF
    echo "%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj
2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj
3 0 obj
<<
/Type /Page
/Parent 2 0 R
/MediaBox [0 0 612 792]
/Contents 4 0 R
/Resources <<
/Font <<
/F1 <<
/Type /Font
/Subtype /Type1
/BaseFont /Helvetica
>>
>>
>>
>>
endobj
4 0 obj
<<
/Length 44
>>
stream
BT
/F1 12 Tf
100 700 Td
(Test CV Document) Tj
ET
endstream
endobj
xref
0 5
0000000000 65535 f
0000000009 00000 n
0000000058 00000 n
0000000115 00000 n
0000000315 00000 n
trailer
<<
/Size 5
/Root 1 0 R
>>
startxref
408
%%EOF" > "$TEST_FILE"
    
    echo -e "${GREEN}✓ Created test file: $TEST_FILE${NC}\n"
else
    TEST_FILE="$2"
    if [ ! -f "$TEST_FILE" ]; then
        echo -e "${RED}Error: File '$TEST_FILE' not found${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Using file: $TEST_FILE${NC}\n"
fi

# API endpoint
API_URL="http://localhost/api/v1/cv"

echo -e "${YELLOW}Uploading CV...${NC}"
echo "URL: $API_URL"
echo "File: $TEST_FILE"
echo ""

# Upload the file
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@$TEST_FILE")

# Split response and status code
HTTP_BODY=$(echo "$RESPONSE" | head -n -1)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

# Display results
echo -e "${YELLOW}Response:${NC}"
echo "$HTTP_BODY" | jq '.' 2>/dev/null || echo "$HTTP_BODY"
echo ""

if [ "$HTTP_CODE" -eq 201 ] || [ "$HTTP_CODE" -eq 200 ]; then
    echo -e "${GREEN}✓ Upload successful! (HTTP $HTTP_CODE)${NC}"
else
    echo -e "${RED}✗ Upload failed! (HTTP $HTTP_CODE)${NC}"
fi

# Cleanup test file if we created it
if [ "$TEST_FILE" = "test_cv.pdf" ]; then
    rm -f "$TEST_FILE"
    echo -e "\n${YELLOW}Cleaned up test file${NC}"
fi

echo ""
echo -e "${YELLOW}===== Test Complete =====${NC}"
