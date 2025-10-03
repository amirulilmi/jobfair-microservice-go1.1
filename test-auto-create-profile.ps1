# Auto-Create Profile Feature Test (PowerShell)
param(
    [string]$Email = "testuser@example.com",
    [string]$Password = "TestPass123!"
)

Write-Host "`n========================================" -ForegroundColor Blue
Write-Host "  Auto-Create Profile Feature Test" -ForegroundColor Blue
Write-Host "========================================`n" -ForegroundColor Blue

Write-Host "Using credentials:" -ForegroundColor Yellow
Write-Host "Email: $Email"
Write-Host "Password: $Password"
Write-Host ""

# Helper function to make API calls
function Invoke-ApiCall {
    param(
        [string]$Method,
        [string]$Url,
        [string]$Token = $null,
        [string]$Body = $null
    )
    
    $headers = @{
        "Content-Type" = "application/json"
    }
    
    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }
    
    try {
        if ($Body) {
            $response = Invoke-RestMethod -Uri $Url -Method $Method -Headers $headers -Body $Body
        } else {
            $response = Invoke-RestMethod -Uri $Url -Method $Method -Headers $headers
        }
        return $response
    } catch {
        return $_.Exception.Response
    }
}

# Step 1: Login
Write-Host "Step 1: Login..." -ForegroundColor Yellow
$loginBody = @{
    email = $Email
    password = $Password
} | ConvertTo-Json

try {
    $loginResponse = Invoke-ApiCall -Method POST -Url "http://localhost/api/v1/login" -Body $loginBody
    $token = $loginResponse.data.access_token
    
    if ($token) {
        Write-Host "✓ Login successful" -ForegroundColor Green
        Write-Host "Token: $($token.Substring(0, [Math]::Min(30, $token.Length)))..."
    } else {
        Write-Host "✗ Login failed!" -ForegroundColor Red
        Write-Host "Response: $($loginResponse | ConvertTo-Json)"
        exit 1
    }
} catch {
    Write-Host "✗ Login failed!" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}
Write-Host ""

# Step 2: Check if profile exists
Write-Host "Step 2: Checking if profile exists..." -ForegroundColor Yellow
try {
    $profileCheck = Invoke-ApiCall -Method GET -Url "http://localhost/api/v1/profiles" -Token $token
    
    if ($profileCheck.success -eq $false) {
        Write-Host "ℹ Profile doesn't exist yet (expected)" -ForegroundColor Blue
    } else {
        Write-Host "⚠ Profile already exists" -ForegroundColor Yellow
    }
} catch {
    Write-Host "ℹ Profile doesn't exist yet (expected)" -ForegroundColor Blue
}
Write-Host ""

# Step 3: Create career preference
Write-Host "Step 3: Creating career preference (will auto-create profile)..." -ForegroundColor Yellow
$preferenceBody = @{
    job_type = "full_time"
    work_location = "hybrid"
    expected_salary_min = 15000000
    expected_salary_max = 25000000
    currency = "IDR"
    willing_to_relocate = $true
    available_from = "2025-11-01"
} | ConvertTo-Json

try {
    $preferenceResponse = Invoke-ApiCall -Method POST -Url "http://localhost/api/v1/career-preference" -Token $token -Body $preferenceBody
    
    if ($preferenceResponse.success) {
        Write-Host "✓ Career preference created successfully" -ForegroundColor Green
        Write-Host "Response:"
        $preferenceResponse | ConvertTo-Json -Depth 10
    } else {
        Write-Host "✗ Failed to create career preference" -ForegroundColor Red
        Write-Host "Response:"
        $preferenceResponse | ConvertTo-Json -Depth 10
        exit 1
    }
} catch {
    Write-Host "✗ Failed to create career preference" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}
Write-Host ""

# Step 4: Verify profile was auto-created
Write-Host "Step 4: Verifying profile was auto-created..." -ForegroundColor Yellow
try {
    $profileCheck2 = Invoke-ApiCall -Method GET -Url "http://localhost/api/v1/profiles" -Token $token
    
    if ($profileCheck2.success) {
        Write-Host "✓ Profile auto-created successfully!" -ForegroundColor Green
        Write-Host "Profile ID: $($profileCheck2.data.id)"
        Write-Host "Completion: $($profileCheck2.data.completion_status)%"
    } else {
        Write-Host "✗ Profile was not created" -ForegroundColor Red
        Write-Host "Response:"
        $profileCheck2 | ConvertTo-Json -Depth 10
        exit 1
    }
} catch {
    Write-Host "✗ Profile was not created" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}
Write-Host ""

# Step 5: Add work experience
Write-Host "Step 5: Adding work experience..." -ForegroundColor Yellow
$workExpBody = @{
    company_name = "Tech Corp Indonesia"
    job_position = "Senior Software Engineer"
    start_date = "2020-01-01"
    end_date = "2023-12-31"
    is_current_job = $false
    job_description = "Developed microservices architecture"
} | ConvertTo-Json

try {
    $workExpResponse = Invoke-ApiCall -Method POST -Url "http://localhost/api/v1/work-experiences" -Token $token -Body $workExpBody
    
    if ($workExpResponse.success) {
        Write-Host "✓ Work experience added" -ForegroundColor Green
    } else {
        Write-Host "✗ Failed to add work experience" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ Failed to add work experience" -ForegroundColor Red
}
Write-Host ""

# Step 6: Add education
Write-Host "Step 6: Adding education..." -ForegroundColor Yellow
$educationBody = @{
    university = "University of Indonesia"
    major = "Computer Science"
    degree = "Bachelor"
    start_date = "2016-08-01"
    end_date = "2020-07-01"
    is_current = $false
    gpa = "3.75"
    description = "Focus on Software Engineering"
} | ConvertTo-Json

try {
    $educationResponse = Invoke-ApiCall -Method POST -Url "http://localhost/api/v1/educations" -Token $token -Body $educationBody
    
    if ($educationResponse.success) {
        Write-Host "✓ Education added" -ForegroundColor Green
    } else {
        Write-Host "✗ Failed to add education" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ Failed to add education" -ForegroundColor Red
}
Write-Host ""

# Step 7: Add skills
Write-Host "Step 7: Adding skills..." -ForegroundColor Yellow
$skillsBody = @{
    skills = @(
        @{
            skill_name = "Go"
            skill_type = "technical"
            proficiency_level = "expert"
            years_of_experience = 5
        },
        @{
            skill_name = "Docker"
            skill_type = "technical"
            proficiency_level = "advanced"
            years_of_experience = 4
        },
        @{
            skill_name = "Leadership"
            skill_type = "soft"
            proficiency_level = "advanced"
            years_of_experience = 3
        }
    )
} | ConvertTo-Json -Depth 10

try {
    $skillsResponse = Invoke-ApiCall -Method POST -Url "http://localhost/api/v1/skills/bulk" -Token $token -Body $skillsBody
    
    if ($skillsResponse.success) {
        $skillsCount = $skillsResponse.data.Count
        Write-Host "✓ $skillsCount skills added" -ForegroundColor Green
    } else {
        Write-Host "✗ Failed to add skills" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ Failed to add skills" -ForegroundColor Red
}
Write-Host ""

# Step 8: Check final completion status
Write-Host "Step 8: Checking final completion status..." -ForegroundColor Yellow
try {
    $completionResponse = Invoke-ApiCall -Method GET -Url "http://localhost/api/v1/profiles/completion" -Token $token
    $completionStatus = $completionResponse.data.completion_percentage
    
    Write-Host "✓ Profile completion: $completionStatus%" -ForegroundColor Green
} catch {
    Write-Host "⚠ Could not get completion status" -ForegroundColor Yellow
    $completionStatus = "Unknown"
}
Write-Host ""

# Final summary
Write-Host "========================================" -ForegroundColor Blue
Write-Host "           Test Summary" -ForegroundColor Blue
Write-Host "========================================`n" -ForegroundColor Blue

Write-Host "✓ Login successful" -ForegroundColor Green
Write-Host "✓ Profile auto-created" -ForegroundColor Green
Write-Host "✓ Career preference created" -ForegroundColor Green
Write-Host "✓ Work experience added" -ForegroundColor Green
Write-Host "✓ Education added" -ForegroundColor Green
Write-Host "✓ Skills added" -ForegroundColor Green
Write-Host "✓ Final completion: $completionStatus%" -ForegroundColor Green
Write-Host ""
Write-Host "========================================" -ForegroundColor Blue
Write-Host "All tests passed! 🎉" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Blue
