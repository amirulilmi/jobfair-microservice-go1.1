# JobFair CV Upload Test Script (PowerShell)
param(
    [Parameter(Mandatory=$true)]
    [string]$Token,
    
    [Parameter(Mandatory=$false)]
    [string]$FilePath
)

Write-Host "`n===== JobFair CV Upload Test =====" -ForegroundColor Yellow
Write-Host ""

# API endpoint
$ApiUrl = "http://localhost/api/v1/cv"

# If no file provided, create a test PDF
if ([string]::IsNullOrEmpty($FilePath)) {
    Write-Host "No file provided, creating test CV file..." -ForegroundColor Yellow
    $FilePath = "test_cv.pdf"
    
    $pdfContent = @"
%PDF-1.4
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
%%EOF
"@
    
    [System.IO.File]::WriteAllText($FilePath, $pdfContent)
    Write-Host "✓ Created test file: $FilePath" -ForegroundColor Green
    Write-Host ""
    $CleanupFile = $true
} else {
    if (-not (Test-Path $FilePath)) {
        Write-Host "Error: File '$FilePath' not found" -ForegroundColor Red
        exit 1
    }
    Write-Host "✓ Using file: $FilePath" -ForegroundColor Green
    Write-Host ""
    $CleanupFile = $false
}

Write-Host "Uploading CV..." -ForegroundColor Yellow
Write-Host "URL: $ApiUrl"
Write-Host "File: $FilePath"
Write-Host ""

# Prepare the file for upload
$fileBytes = [System.IO.File]::ReadAllBytes($FilePath)
$fileName = [System.IO.Path]::GetFileName($FilePath)

# Create multipart form data
$boundary = [System.Guid]::NewGuid().ToString()
$LF = "`r`n"

$bodyLines = (
    "--$boundary",
    "Content-Disposition: form-data; name=`"file`"; filename=`"$fileName`"",
    "Content-Type: application/octet-stream$LF",
    [System.Text.Encoding]::UTF8.GetString($fileBytes),
    "--$boundary--$LF"
) -join $LF

# Make the request
try {
    $response = Invoke-WebRequest -Uri $ApiUrl `
        -Method POST `
        -Headers @{
            "Authorization" = "Bearer $Token"
        } `
        -ContentType "multipart/form-data; boundary=$boundary" `
        -Body $bodyLines `
        -UseBasicParsing
    
    Write-Host "Response:" -ForegroundColor Yellow
    $response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 10
    Write-Host ""
    
    if ($response.StatusCode -eq 200 -or $response.StatusCode -eq 201) {
        Write-Host "✓ Upload successful! (HTTP $($response.StatusCode))" -ForegroundColor Green
    }
} catch {
    Write-Host "Response:" -ForegroundColor Yellow
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        try {
            $responseBody | ConvertFrom-Json | ConvertTo-Json -Depth 10
        } catch {
            Write-Host $responseBody
        }
        Write-Host ""
        Write-Host "✗ Upload failed! (HTTP $($_.Exception.Response.StatusCode.value__))" -ForegroundColor Red
    } else {
        Write-Host $_.Exception.Message -ForegroundColor Red
    }
}

# Cleanup test file if we created it
if ($CleanupFile) {
    Remove-Item $FilePath -Force
    Write-Host "`nCleaned up test file" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "===== Test Complete =====" -ForegroundColor Yellow
Write-Host ""

# Usage instructions
Write-Host "Usage Examples:" -ForegroundColor Cyan
Write-Host "  .\test-upload-cv.ps1 -Token 'eyJhbGc...'"
Write-Host "  .\test-upload-cv.ps1 -Token 'eyJhbGc...' -FilePath 'C:\path\to\cv.pdf'"
