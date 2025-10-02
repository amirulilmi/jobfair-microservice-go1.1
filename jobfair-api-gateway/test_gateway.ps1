# API Gateway Testing Script (PowerShell)

$GatewayURL = "http://localhost:8000"
$Passed = 0
$Failed = 0
$Total = 0

function Write-TestResult {
    param($TestName, $Success, $Response)
    
    $script:Total++
    
    if ($Success) {
        Write-Host "✓ Test $Total : $TestName " -ForegroundColor Green -NoNewline
        Write-Host "PASSED" -ForegroundColor Green
        $script:Passed++
    } else {
        Write-Host "✗ Test $Total : $TestName " -ForegroundColor Red -NoNewline
        Write-Host "FAILED" -ForegroundColor Red
        Write-Host "Response: $Response" -ForegroundColor Red
        $script:Failed++
    }
}

Write-Host "============================================" -ForegroundColor Cyan
Write-Host "API Gateway Testing" -ForegroundColor Cyan
Write-Host "Gateway URL: $GatewayURL" -ForegroundColor Cyan
Write-Host "============================================`n" -ForegroundColor Cyan

# Test 1: Gateway Health Check
Write-Host "[1] Testing Gateway Health Check..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$GatewayURL/health" -Method Get
    if ($response.status -eq "healthy") {
        Write-TestResult "Gateway Health Check" $true $response
    } else {
        Write-TestResult "Gateway Health Check" $false $response
    }
} catch {
    Write-TestResult "Gateway Health Check" $false $_.Exception.Message
}

# Test 2: Auth Service via Gateway - Register
Write-Host "`n[2] Testing Auth Service (Register) via Gateway..." -ForegroundColor Yellow
$randomEmail = "test$(Get-Random)@example.com"
$registerBody = @{
    email = $randomEmail
    password = "password123"
    full_name = "Test User"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$GatewayURL/api/v1/auth/register" `
        -Method Post `
        -Body $registerBody `
        -ContentType "application/json"
    
    if ($response.success) {
        Write-TestResult "Auth Register via Gateway" $true $response
    } else {
        Write-TestResult "Auth Register via Gateway" $false $response
    }
} catch {
    Write-TestResult "Auth Register via Gateway" $false $_.Exception.Message
}

# Test 3: Auth Service via Gateway - Login
Write-Host "`n[3] Testing Auth Service (Login) via Gateway..." -ForegroundColor Yellow
$loginBody = @{
    email = $randomEmail
    password = "password123"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$GatewayURL/api/v1/auth/login" `
        -Method Post `
        -Body $loginBody `
        -ContentType "application/json"
    
    if ($response.success -and $response.data.token) {
        Write-TestResult "Auth Login via Gateway" $true $response
        $Token = $response.data.token
        Write-Host "Token obtained: $($Token.Substring(0, 20))..." -ForegroundColor Cyan
    } else {
        Write-TestResult "Auth Login via Gateway" $false $response
    }
} catch {
    Write-TestResult "Auth Login via Gateway" $false $_.Exception.Message
}

if ($Token) {
    # Test 4: Create Profile via Gateway
    Write-Host "`n[4] Testing Profile Service (Create Profile) via Gateway..." -ForegroundColor Yellow
    $profileBody = @{
        full_name = "John Doe"
        phone_number = "081234567890"
    } | ConvertTo-Json

    try {
        $response = Invoke-RestMethod -Uri "$GatewayURL/api/v1/profiles" `
            -Method Post `
            -Headers @{Authorization = "Bearer $Token"} `
            -Body $profileBody `
            -ContentType "application/json"
        
        if ($response.success) {
            Write-TestResult "Create Profile via Gateway" $true $response
        } else {
            Write-TestResult "Create Profile via Gateway" $false $response
        }
    } catch {
        Write-TestResult "Create Profile via Gateway" $false $_.Exception.Message
    }

    # Test 5: Get Profile via Gateway
    Write-Host "`n[5] Testing Profile Service (Get Profile) via Gateway..." -ForegroundColor Yellow
    try {
        $response = Invoke-RestMethod -Uri "$GatewayURL/api/v1/profiles" `
            -Method Get `
            -Headers @{Authorization = "Bearer $Token"}
        
        if ($response.success) {
            Write-TestResult "Get Profile via Gateway" $true $response
        } else {
            Write-TestResult "Get Profile via Gateway" $false $response
        }
    } catch {
        Write-TestResult "Get Profile via Gateway" $false $_.Exception.Message
    }

    # Test 6: Create Work Experience via Gateway
    Write-Host "`n[6] Testing Work Experience (Create) via Gateway..." -ForegroundColor Yellow
    $workExpBody = @{
        company_name = "PT Test Company"
        job_position = "Software Engineer"
        start_date = "2020-01-01T00:00:00Z"
        is_current_job = $true
        job_description = "Developing applications"
    } | ConvertTo-Json

    try {
        $response = Invoke-RestMethod -Uri "$GatewayURL/api/v1/work-experiences" `
            -Method Post `
            -Headers @{Authorization = "Bearer $Token"} `
            -Body $workExpBody `
            -ContentType "application/json"
        
        if ($response.success) {
            Write-TestResult "Create Work Experience via Gateway" $true $response
        } else {
            Write-TestResult "Create Work Experience via Gateway" $false $response
        }
    } catch {
        Write-TestResult "Create Work Experience via Gateway" $false $_.Exception.Message
    }

    # Test 7: Get Work Experiences via Gateway
    Write-Host "`n[7] Testing Work Experience (Get All) via Gateway..." -ForegroundColor Yellow
    try {
        $response = Invoke-RestMethod -Uri "$GatewayURL/api/v1/work-experiences" `
            -Method Get `
            -Headers @{Authorization = "Bearer $Token"}
        
        if ($response.success) {
            Write-TestResult "Get Work Experiences via Gateway" $true $response
        } else {
            Write-TestResult "Get Work Experiences via Gateway" $false $response
        }
    } catch {
        Write-TestResult "Get Work Experiences via Gateway" $false $_.Exception.Message
    }

    # Test 8: Create Education via Gateway
    Write-Host "`n[8] Testing Education (Create) via Gateway..." -ForegroundColor Yellow
    $educationBody = @{
        university = "Institut Teknologi Bandung"
        major = "Computer Science"
        degree = "Bachelor"
        start_date = "2016-08-01T00:00:00Z"
        end_date = "2020-07-31T00:00:00Z"
        gpa = 3.75
    } | ConvertTo-Json

    try {
        $response = Invoke-RestMethod -Uri "$GatewayURL/api/v1/educations" `
            -Method Post `
            -Headers @{Authorization = "Bearer $Token"} `
            -Body $educationBody `
            -ContentType "application/json"
        
        if ($response.success) {
            Write-TestResult "Create Education via Gateway" $true $response
        } else {
            Write-TestResult "Create Education via Gateway" $false $response
        }
    } catch {
        Write-TestResult "Create Education via Gateway" $false $_.Exception.Message
    }

    # Test 9: Create Skills (Bulk) via Gateway
    Write-Host "`n[9] Testing Skills (Bulk Create) via Gateway..." -ForegroundColor Yellow
    $skillsBody = @{
        technical_skills = @(
            @{
                skill_name = "Go Programming"
                skill_type = "technical"
                proficiency_level = "advanced"
            }
        )
        soft_skills = @(
            @{
                skill_name = "Leadership"
                skill_type = "soft"
                proficiency_level = "intermediate"
            }
        )
    } | ConvertTo-Json -Depth 10

    try {
        $response = Invoke-RestMethod -Uri "$GatewayURL/api/v1/skills/bulk" `
            -Method Post `
            -Headers @{Authorization = "Bearer $Token"} `
            -Body $skillsBody `
            -ContentType "application/json"
        
        if ($response.success) {
            Write-TestResult "Create Skills (Bulk) via Gateway" $true $response
        } else {
            Write-TestResult "Create Skills (Bulk) via Gateway" $false $response
        }
    } catch {
        Write-TestResult "Create Skills (Bulk) via Gateway" $false $_.Exception.Message
    }

    # Test 10: Get Complete Profile via Gateway
    Write-Host "`n[10] Testing Get Complete Profile via Gateway..." -ForegroundColor Yellow
    try {
        $response = Invoke-RestMethod -Uri "$GatewayURL/api/v1/profiles/full" `
            -Method Get `
            -Headers @{Authorization = "Bearer $Token"}
        
        if ($response.success) {
            Write-TestResult "Get Complete Profile via Gateway" $true $response
        } else {
            Write-TestResult "Get Complete Profile via Gateway" $false $response
        }
    } catch {
        Write-TestResult "Get Complete Profile via Gateway" $false $_.Exception.Message
    }
} else {
    Write-Host "`n⚠️ Skipping authenticated tests (no token)" -ForegroundColor Yellow
}

# Summary
Write-Host "`n============================================" -ForegroundColor Cyan
Write-Host "Test Summary" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "Total Tests: $Total"
Write-Host "Passed: $Passed" -ForegroundColor Green
Write-Host "Failed: $Failed" -ForegroundColor Red

if ($Failed -eq 0) {
    Write-Host "`n🎉 All tests passed! API Gateway is working correctly." -ForegroundColor Green
    exit 0
} else {
    Write-Host "`n⚠️  Some tests failed. Please check the errors above." -ForegroundColor Red
    exit 1
}
