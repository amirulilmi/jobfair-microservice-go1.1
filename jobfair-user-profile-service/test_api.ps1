# JobFair User Profile Service - API Testing Script (PowerShell)
# This script tests all endpoints of the User Profile Service

$BaseUrl = "http://localhost:8083"
$AuthServiceUrl = "http://localhost:8080"

# Test counters
$Passed = 0
$Failed = 0
$Total = 0

# Colors
function Write-Success { Write-Host $args -ForegroundColor Green }
function Write-Error { Write-Host $args -ForegroundColor Red }
function Write-Info { Write-Host $args -ForegroundColor Cyan }
function Write-Warning { Write-Host $args -ForegroundColor Yellow }

# Function to print test result
function Print-Result {
    param($TestName, $Success, $Response)
    
    $script:Total++
    
    if ($Success) {
        Write-Success "✓ Test $Total : $TestName PASSED"
        $script:Passed++
    } else {
        Write-Error "✗ Test $Total : $TestName FAILED"
        Write-Error "Response: $Response"
        $script:Failed++
    }
}

# Function to make API call
function Invoke-ApiCall {
    param(
        [string]$Method,
        [string]$Endpoint,
        [string]$Body = "",
        [string]$Token = ""
    )
    
    $headers = @{
        "Content-Type" = "application/json"
    }
    
    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }
    
    try {
        if ($Body) {
            $response = Invoke-RestMethod -Uri "$BaseUrl$Endpoint" -Method $Method -Headers $headers -Body $Body
        } else {
            $response = Invoke-RestMethod -Uri "$BaseUrl$Endpoint" -Method $Method -Headers $headers
        }
        return $response
    } catch {
        return $_.Exception.Message
    }
}

Write-Info "============================================"
Write-Info "JobFair User Profile Service - API Testing"
Write-Info "============================================`n"

# ============================================
# TEST 1: Health Check (No Auth)
# ============================================
Write-Warning "`n[1] Testing Health Check..."
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/health" -Method Get
    if ($response.status -eq "ok") {
        Print-Result "Health Check" $true $response
    } else {
        Print-Result "Health Check" $false $response
    }
} catch {
    Print-Result "Health Check" $false $_.Exception.Message
}

# ============================================
# SETUP: Get JWT Token from Auth Service
# ============================================
Write-Warning "`n[SETUP] Getting JWT Token from Auth Service..."
$Email = Read-Host "Enter email"
$Password = Read-Host "Enter password" -AsSecureString
$PasswordPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($Password))

$loginBody = @{
    email = $Email
    password = $PasswordPlain
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$AuthServiceUrl/api/v1/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
    $Token = $loginResponse.data.token
    
    if ($Token) {
        Write-Success "✓ JWT Token obtained successfully"
        Write-Info "Token: $($Token.Substring(0, [Math]::Min(20, $Token.Length)))...`n"
    } else {
        Write-Error "✗ Failed to get JWT token"
        exit 1
    }
} catch {
    Write-Error "✗ Failed to get JWT token. Please make sure:"
    Write-Error "  1. Auth service is running at $AuthServiceUrl"
    Write-Error "  2. Email and password are correct"
    Write-Error "  3. User is registered"
    exit 1
}

# ============================================
# TEST 2: Create Profile
# ============================================
Write-Warning "`n[2] Testing Create Profile..."
$body = @{
    full_name = "Test User"
    phone_number = "081234567890"
} | ConvertTo-Json

$response = Invoke-ApiCall -Method "POST" -Endpoint "/api/v1/profiles" -Body $body -Token $Token
if ($response.success) {
    Print-Result "Create Profile" $true $response
} else {
    Print-Result "Create Profile" $false $response
}

# ============================================
# TEST 3: Get Profile
# ============================================
Write-Warning "`n[3] Testing Get Profile..."
$response = Invoke-ApiCall -Method "GET" -Endpoint "/api/v1/profiles" -Token $Token
if ($response.success) {
    Print-Result "Get Profile" $true $response
} else {
    Print-Result "Get Profile" $false $response
}

# ============================================
# TEST 4: Update Profile
# ============================================
Write-Warning "`n[4] Testing Update Profile..."
$body = @{
    bio = "Software Engineer with passion for microservices"
    city = "Jakarta"
} | ConvertTo-Json

$response = Invoke-ApiCall -Method "PUT" -Endpoint "/api/v1/profiles" -Body $body -Token $Token
if ($response.success) {
    Print-Result "Update Profile" $true $response
} else {
    Print-Result "Update Profile" $false $response
}

# ============================================
# TEST 5: Create Work Experience
# ============================================
Write-Warning "`n[5] Testing Create Work Experience..."
$body = @{
    company_name = "PT Test"
    job_position = "Software Engineer"
    start_date = "2020-01-01T00:00:00Z"
    is_current_job = $true
    job_description = "Developing microservices"
} | ConvertTo-Json

$response = Invoke-ApiCall -Method "POST" -Endpoint "/api/v1/work-experiences" -Body $body -Token $Token
if ($response.success) {
    Print-Result "Create Work Experience" $true $response
} else {
    Print-Result "Create Work Experience" $false $response
}

# ============================================
# TEST 6: Get All Work Experiences
# ============================================
Write-Warning "`n[6] Testing Get All Work Experiences..."
$response = Invoke-ApiCall -Method "GET" -Endpoint "/api/v1/work-experiences" -Token $Token
if ($response.success) {
    Print-Result "Get All Work Experiences" $true $response
} else {
    Print-Result "Get All Work Experiences" $false $response
}

# ============================================
# TEST 7: Create Education
# ============================================
Write-Warning "`n[7] Testing Create Education..."
$body = @{
    university = "Universitas Indonesia"
    major = "Computer Science"
    degree = "Bachelor"
    start_date = "2016-08-01T00:00:00Z"
    end_date = "2020-07-31T00:00:00Z"
    gpa = 3.75
} | ConvertTo-Json

$response = Invoke-ApiCall -Method "POST" -Endpoint "/api/v1/educations" -Body $body -Token $Token
if ($response.success) {
    Print-Result "Create Education" $true $response
} else {
    Print-Result "Create Education" $false $response
}

# ============================================
# TEST 8: Get All Educations
# ============================================
Write-Warning "`n[8] Testing Get All Educations..."
$response = Invoke-ApiCall -Method "GET" -Endpoint "/api/v1/educations" -Token $Token
if ($response.success) {
    Print-Result "Get All Educations" $true $response
} else {
    Print-Result "Get All Educations" $false $response
}

# ============================================
# TEST 9: Create Certification
# ============================================
Write-Warning "`n[9] Testing Create Certification..."
$body = @{
    certification_name = "AWS Certified Developer"
    organizer = "Amazon Web Services"
    issue_date = "2023-01-15T00:00:00Z"
} | ConvertTo-Json

$response = Invoke-ApiCall -Method "POST" -Endpoint "/api/v1/certifications" -Body $body -Token $Token
if ($response.success) {
    Print-Result "Create Certification" $true $response
} else {
    Print-Result "Create Certification" $false $response
}

# ============================================
# TEST 10: Get All Certifications
# ============================================
Write-Warning "`n[10] Testing Get All Certifications..."
$response = Invoke-ApiCall -Method "GET" -Endpoint "/api/v1/certifications" -Token $Token
if ($response.success) {
    Print-Result "Get All Certifications" $true $response
} else {
    Print-Result "Get All Certifications" $false $response
}

# ============================================
# TEST 11: Create Skills (Bulk)
# ============================================
Write-Warning "`n[11] Testing Create Skills (Bulk)..."
$body = @{
    technical_skills = @(
        @{
            skill_name = "Go"
            skill_type = "technical"
            proficiency_level = "expert"
        },
        @{
            skill_name = "PostgreSQL"
            skill_type = "technical"
            proficiency_level = "advanced"
        }
    )
    soft_skills = @(
        @{
            skill_name = "Leadership"
            skill_type = "soft"
            proficiency_level = "advanced"
        }
    )
} | ConvertTo-Json -Depth 10

$response = Invoke-ApiCall -Method "POST" -Endpoint "/api/v1/skills/bulk" -Body $body -Token $Token
if ($response.success) {
    Print-Result "Create Skills (Bulk)" $true $response
} else {
    Print-Result "Create Skills (Bulk)" $false $response
}

# ============================================
# TEST 12: Get All Skills
# ============================================
Write-Warning "`n[12] Testing Get All Skills..."
$response = Invoke-ApiCall -Method "GET" -Endpoint "/api/v1/skills" -Token $Token
if ($response.success) {
    Print-Result "Get All Skills" $true $response
} else {
    Print-Result "Get All Skills" $false $response
}

# ============================================
# TEST 13: Create Career Preference
# ============================================
Write-Warning "`n[13] Testing Create Career Preference..."
$body = @{
    is_actively_looking = $true
    expected_salary_min = 15000000
    expected_salary_max = 25000000
    salary_currency = "IDR"
    is_negotiable = $true
    preferred_work_types = @("remote", "hybrid")
    preferred_locations = @("Jakarta")
    willing_to_relocate = $false
} | ConvertTo-Json

$response = Invoke-ApiCall -Method "POST" -Endpoint "/api/v1/career-preference" -Body $body -Token $Token
if ($response.success) {
    Print-Result "Create Career Preference" $true $response
} else {
    Print-Result "Create Career Preference" $false $response
}

# ============================================
# TEST 14: Get Career Preference
# ============================================
Write-Warning "`n[14] Testing Get Career Preference..."
$response = Invoke-ApiCall -Method "GET" -Endpoint "/api/v1/career-preference" -Token $Token
if ($response.success) {
    Print-Result "Get Career Preference" $true $response
} else {
    Print-Result "Get Career Preference" $false $response
}

# ============================================
# TEST 15: Create Position Preferences
# ============================================
Write-Warning "`n[15] Testing Create Position Preferences..."
$body = @{
    positions = @(
        @{
            position_name = "Software Engineer"
            priority = 1
        },
        @{
            position_name = "Backend Developer"
            priority = 2
        }
    )
} | ConvertTo-Json -Depth 10

$response = Invoke-ApiCall -Method "POST" -Endpoint "/api/v1/position-preferences" -Body $body -Token $Token
if ($response.success) {
    Print-Result "Create Position Preferences" $true $response
} else {
    Print-Result "Create Position Preferences" $false $response
}

# ============================================
# TEST 16: Get Position Preferences
# ============================================
Write-Warning "`n[16] Testing Get Position Preferences..."
$response = Invoke-ApiCall -Method "GET" -Endpoint "/api/v1/position-preferences" -Token $Token
if ($response.success) {
    Print-Result "Get Position Preferences" $true $response
} else {
    Print-Result "Get Position Preferences" $false $response
}

# ============================================
# TEST 17: Get Complete Profile
# ============================================
Write-Warning "`n[17] Testing Get Complete Profile..."
$response = Invoke-ApiCall -Method "GET" -Endpoint "/api/v1/profiles/full" -Token $Token
if ($response.success) {
    Print-Result "Get Complete Profile" $true $response
} else {
    Print-Result "Get Complete Profile" $false $response
}

# ============================================
# TEST 18: Get Profile Completion Status
# ============================================
Write-Warning "`n[18] Testing Get Profile Completion Status..."
$response = Invoke-ApiCall -Method "GET" -Endpoint "/api/v1/profiles/completion" -Token $Token
if ($response.data.completion_status -ge 0) {
    Print-Result "Get Profile Completion Status" $true $response
    Write-Info "Profile Completion: $($response.data.completion_status)%"
} else {
    Print-Result "Get Profile Completion Status" $false $response
}

# ============================================
# SUMMARY
# ============================================
Write-Info "`n============================================"
Write-Info "Test Summary"
Write-Info "============================================"
Write-Host "Total Tests: $Total"
Write-Success "Passed: $Passed"
Write-Error "Failed: $Failed"

if ($Failed -eq 0) {
    Write-Success "`n🎉 All tests passed! Service is working correctly."
    exit 0
} else {
    Write-Error "`n⚠️  Some tests failed. Please check the errors above."
    exit 1
}
