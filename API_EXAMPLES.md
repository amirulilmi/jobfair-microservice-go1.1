# 📘 API Request & Response Examples

Dokumen ini berisi contoh request dan response untuk semua endpoint yang telah diperbaiki.

---

## 1. Bulk Skills API

### ✅ Request (FIXED)
```http
POST http://localhost/api/v1/skills/bulk
Authorization: Bearer YOUR_TOKEN
Content-Type: application/json

{
  "skills": [
    {
      "skill_name": "PostgreSQL",
      "skill_type": "technical",
      "proficiency_level": "advanced",
      "years_of_experience": 4
    },
    {
      "skill_name": "Docker",
      "skill_type": "technical",
      "proficiency_level": "intermediate",
      "years_of_experience": 3
    },
    {
      "skill_name": "Kubernetes",
      "skill_type": "technical",
      "proficiency_level": "intermediate",
      "years_of_experience": 2
    },
    {
      "skill_name": "Communication",
      "skill_type": "soft",
      "proficiency_level": "expert",
      "years_of_experience": 5
    }
  ]
}
```

### ✅ Response (Success)
```json
{
  "success": true,
  "message": "Skills created successfully",
  "data": [
    {
      "id": 1,
      "profile_id": 1,
      "skill_name": "PostgreSQL",
      "skill_type": "technical",
      "proficiency_level": "advanced",
      "years_of_experience": 4,
      "created_at": "2025-10-03T10:00:00Z",
      "updated_at": "2025-10-03T10:00:00Z"
    },
    {
      "id": 2,
      "profile_id": 1,
      "skill_name": "Docker",
      "skill_type": "technical",
      "proficiency_level": "intermediate",
      "years_of_experience": 3,
      "created_at": "2025-10-03T10:00:00Z",
      "updated_at": "2025-10-03T10:00:00Z"
    },
    {
      "id": 3,
      "profile_id": 1,
      "skill_name": "Kubernetes",
      "skill_type": "technical",
      "proficiency_level": "intermediate",
      "years_of_experience": 2,
      "created_at": "2025-10-03T10:00:00Z",
      "updated_at": "2025-10-03T10:00:00Z"
    },
    {
      "id": 4,
      "profile_id": 1,
      "skill_name": "Communication",
      "skill_type": "soft",
      "proficiency_level": "expert",
      "years_of_experience": 5,
      "created_at": "2025-10-03T10:00:00Z",
      "updated_at": "2025-10-03T10:00:00Z"
    }
  ]
}
```

### cURL Example
```bash
curl -X POST http://localhost/api/v1/skills/bulk \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "skills": [
      {
        "skill_name": "PostgreSQL",
        "skill_type": "technical",
        "proficiency_level": "advanced",
        "years_of_experience": 4
      },
      {
        "skill_name": "Docker",
        "skill_type": "technical",
        "proficiency_level": "intermediate",
        "years_of_experience": 3
      }
    ]
  }'
```

---

## 2. Career Preference API

### ✅ Request (FIXED)
```http
POST http://localhost/api/v1/career-preference
Authorization: Bearer YOUR_TOKEN
Content-Type: application/json

{
  "job_type": "full_time",
  "work_location": "hybrid",
  "expected_salary_min": 15000000,
  "expected_salary_max": 25000000,
  "currency": "IDR",
  "willing_to_relocate": true,
  "available_from": "2025-11-01T00:00:00Z"
}
```

### ✅ Response (Success)
```json
{
  "success": true,
  "message": "Career preference saved successfully",
  "data": {
    "id": 1,
    "profile_id": 1,
    "is_actively_looking": true,
    "expected_salary_min": 15000000,
    "expected_salary_max": 25000000,
    "salary_currency": "IDR",
    "is_negotiable": true,
    "preferred_work_types": "full_time,hybrid",
    "preferred_locations": "",
    "willing_to_relocate": true,
    "available_start_date": "2025-11-01T00:00:00Z",
    "created_at": "2025-10-03T10:00:00Z",
    "updated_at": "2025-10-03T10:00:00Z"
  }
}
```

### cURL Example
```bash
curl -X POST http://localhost/api/v1/career-preference \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_type": "full_time",
    "work_location": "hybrid",
    "expected_salary_min": 15000000,
    "expected_salary_max": 25000000,
    "currency": "IDR",
    "willing_to_relocate": true,
    "available_from": "2025-11-01T00:00:00Z"
  }'
```

### GET Career Preference
```http
GET http://localhost/api/v1/career-preference
Authorization: Bearer YOUR_TOKEN
```

---

## 3. Position Preferences API

### ✅ Request (FIXED)
```http
POST http://localhost/api/v1/position-preferences
Authorization: Bearer YOUR_TOKEN
Content-Type: application/json

{
  "positions": [
    "Software Engineer",
    "Backend Developer",
    "Full Stack Developer"
  ]
}
```

### ✅ Response (Success)
```json
{
  "success": true,
  "message": "Position preferences created successfully",
  "data": [
    {
      "id": 1,
      "profile_id": 1,
      "position_name": "Software Engineer",
      "priority": 1,
      "created_at": "2025-10-03T10:00:00Z",
      "updated_at": "2025-10-03T10:00:00Z"
    },
    {
      "id": 2,
      "profile_id": 1,
      "position_name": "Backend Developer",
      "priority": 2,
      "created_at": "2025-10-03T10:00:00Z",
      "updated_at": "2025-10-03T10:00:00Z"
    },
    {
      "id": 3,
      "profile_id": 1,
      "position_name": "Full Stack Developer",
      "priority": 3,
      "created_at": "2025-10-03T10:00:00Z",
      "updated_at": "2025-10-03T10:00:00Z"
    }
  ]
}
```

### cURL Example
```bash
curl -X POST http://localhost/api/v1/position-preferences \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "positions": [
      "Software Engineer",
      "Backend Developer",
      "Full Stack Developer"
    ]
  }'
```

### GET Position Preferences
```http
GET http://localhost/api/v1/position-preferences
Authorization: Bearer YOUR_TOKEN
```

### DELETE Position Preference
```http
DELETE http://localhost/api/v1/position-preferences/1
Authorization: Bearer YOUR_TOKEN
```

---

## 4. CV Upload API

### ✅ Request (FIXED)
```http
POST http://localhost/api/v1/cv
Authorization: Bearer YOUR_TOKEN
Content-Type: multipart/form-data

file: [binary file data - PDF or DOCX]
```

### ✅ Response (Success)
```json
{
  "success": true,
  "message": "CV uploaded successfully",
  "data": {
    "id": 1,
    "profile_id": 1,
    "file_name": "resume.pdf",
    "file_url": "/root/uploads/1_1696320000.pdf",
    "file_size": 245678,
    "file_type": ".pdf",
    "is_verified": false,
    "uploaded_at": "2025-10-03T10:00:00Z",
    "created_at": "2025-10-03T10:00:00Z",
    "updated_at": "2025-10-03T10:00:00Z"
  }
}
```

### cURL Example (Linux/Mac)
```bash
curl -X POST http://localhost/api/v1/cv \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/your/resume.pdf"
```

### cURL Example (Windows PowerShell)
```powershell
$file = Get-Item "C:\path\to\your\resume.pdf"
Invoke-RestMethod -Uri "http://localhost/api/v1/cv" `
  -Method Post `
  -Headers @{ Authorization = "Bearer YOUR_TOKEN" } `
  -Form @{ file = $file }
```

### Postman Setup
1. Method: `POST`
2. URL: `http://localhost/api/v1/cv`
3. Headers: 
   - `Authorization: Bearer YOUR_TOKEN`
4. Body:
   - Select `form-data`
   - Key: `file` (change type to File)
   - Value: Select your PDF/DOCX file

### GET CV
```http
GET http://localhost/api/v1/cv
Authorization: Bearer YOUR_TOKEN
```

### DELETE CV
```http
DELETE http://localhost/api/v1/cv
Authorization: Bearer YOUR_TOKEN
```

---

## 5. Job Service API (via Gateway)

### List Jobs
```http
GET http://localhost/api/v1/jobs?page=1&limit=10
```

**Response:**
```json
{
  "success": true,
  "data": {
    "jobs": [...],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 50
    }
  }
}
```

### Get Job Detail
```http
GET http://localhost/api/v1/jobs/1
```

### Create Job (Company only)
```http
POST http://localhost/api/v1/jobs
Authorization: Bearer COMPANY_TOKEN
Content-Type: application/json

{
  "title": "Senior Backend Developer",
  "description": "We are looking for...",
  "employment_type": "full_time",
  "location": "Jakarta",
  "salary_min": 15000000,
  "salary_max": 25000000,
  "requirements": ["3+ years experience", "Go/Python"],
  "benefits": ["Health insurance", "Remote work"]
}
```

### Apply to Job (Job Seeker only)
```http
POST http://localhost/api/v1/jobs/1/apply
Authorization: Bearer JOBSEEKER_TOKEN
Content-Type: application/json

{
  "cover_letter": "I am interested in this position..."
}
```

---

## 6. Auth Service API (via Gateway)

### Register
```http
POST http://localhost/api/v1/auth/register
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123!",
  "role": "jobseeker",
  "full_name": "John Doe"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "user": {
      "id": 1,
      "email": "john@example.com",
      "role": "jobseeker",
      "is_verified": false
    }
  }
}
```

### Login
```http
POST http://localhost/api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "email": "john@example.com",
      "role": "jobseeker"
    }
  }
}
```

---

## 7. Company Service API (via Gateway)

### List Companies
```http
GET http://localhost/api/v1/companies?page=1&limit=10
```

### Get Company Detail
```http
GET http://localhost/api/v1/companies/1
```

### Create Company Profile (Auth required)
```http
POST http://localhost/api/v1/companies
Authorization: Bearer COMPANY_TOKEN
Content-Type: application/json

{
  "company_name": "Tech Startup Inc",
  "industry": "Technology",
  "company_size": "50-100",
  "description": "We are a fast-growing tech startup...",
  "website": "https://techstartup.com",
  "location": "Jakarta"
}
```

---

## Complete Workflow Example

### 1. Register & Login
```bash
# Register
curl -X POST http://localhost/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test1234!",
    "role": "jobseeker",
    "full_name": "Test User"
  }'

# Login
TOKEN=$(curl -s -X POST http://localhost/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test1234!"
  }' | jq -r '.data.access_token')

echo "Token: $TOKEN"
```

### 2. Create Profile
```bash
curl -X POST http://localhost/api/v1/profiles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "phone_number": "081234567890",
    "date_of_birth": "1990-01-01",
    "gender": "male"
  }'
```

### 3. Add Skills
```bash
curl -X POST http://localhost/api/v1/skills/bulk \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "skills": [
      {"skill_name": "Go", "skill_type": "technical", "proficiency_level": "advanced", "years_of_experience": 3},
      {"skill_name": "PostgreSQL", "skill_type": "technical", "proficiency_level": "intermediate", "years_of_experience": 2}
    ]
  }'
```

### 4. Set Career Preferences
```bash
curl -X POST http://localhost/api/v1/career-preference \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_type": "full_time",
    "work_location": "remote",
    "expected_salary_min": 15000000,
    "expected_salary_max": 25000000,
    "currency": "IDR",
    "willing_to_relocate": false
  }'
```

### 5. Set Position Preferences
```bash
curl -X POST http://localhost/api/v1/position-preferences \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "positions": ["Backend Developer", "Full Stack Developer"]
  }'
```

### 6. Upload CV
```bash
curl -X POST http://localhost/api/v1/cv \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@resume.pdf"
```

### 7. Browse & Apply to Jobs
```bash
# List jobs
curl http://localhost/api/v1/jobs?page=1&limit=10

# Apply to job
curl -X POST http://localhost/api/v1/jobs/1/apply \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"cover_letter": "I am interested..."}'
```

---

## Error Responses

### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized",
  "code": "UNAUTHORIZED"
}
```

### 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid request",
  "error": "validation error details",
  "code": "VALIDATION_ERROR"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Resource not found",
  "code": "NOT_FOUND"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Internal server error",
  "error": "error details",
  "code": "INTERNAL_ERROR"
}
```

---

## Notes

- Replace `YOUR_TOKEN` dengan token yang didapat dari login
- Semua endpoint sekarang accessible via `http://localhost` (tanpa port)
- File upload maximum size: 10MB
- Allowed file types untuk CV: `.pdf`, `.docx`
- Semua timestamp dalam format ISO 8601 (UTC)
