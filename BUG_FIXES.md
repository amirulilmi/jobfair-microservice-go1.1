# 🔧 Bug Fixes Documentation

## Ringkasan Perbaikan

Dokumen ini menjelaskan perbaikan untuk 5 masalah utama yang ditemukan pada sistem JobFair Microservice.

---

## ✅ 1. Fix Bulk Skills API (POST /api/v1/skills/bulk)

### **Masalah:**
- Data skills tidak tersimpan ke database
- Response mengembalikan `data: null` meskipun success = true
- Request JSON tidak sesuai dengan model

### **Request Format (Sebelum):**
```json
{
  "technical_skills": [...],
  "soft_skills": [...]
}
```

### **Request Format (Setelah - FIXED):**
```json
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
    }
  ]
}
```

### **Perubahan:**
1. **File:** `internal/models/skill.go`
   - Update `BulkSkillRequest` struct dari `TechnicalSkills` dan `SoftSkills` menjadi `Skills`

2. **File:** `internal/services/skill_service.go`
   - Update logic `CreateBulk()` untuk iterate melalui `req.Skills`
   - Setiap skill menggunakan `skill_type` dari request (bisa "technical" atau "soft")
   - Add error handling jika tidak ada skills yang berhasil dibuat

### **Test:**
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

## ✅ 2. Fix Career Preference API (POST /api/v1/career-preference)

### **Masalah:**
- Foreign key constraint error: `career_preferences_profile_id_fkey`
- Request fields tidak sesuai dengan model

### **Request Format (Sebelum):**
```json
{
  "is_actively_looking": true,
  "expected_salary_min": 15000000,
  ...
}
```

### **Request Format (Setelah - FIXED):**
```json
{
  "job_type": "full_time",
  "work_location": "hybrid",
  "expected_salary_min": 15000000,
  "expected_salary_max": 25000000,
  "currency": "IDR",
  "willing_to_relocate": true,
  "available_from": "2025-11-01"
}
```

### **Perubahan:**
1. **File:** `internal/models/preference.go`
   - Update `CareerPreferenceRequest` struct dengan field yang sesuai:
     - `job_type` → maps to `PreferredWorkTypes`
     - `work_location` → maps to `PreferredWorkTypes`
     - `currency` → maps to `SalaryCurrency`
     - `available_from` → maps to `AvailableStartDate`

2. **File:** `internal/services/preference_service.go`
   - Update method signature dari banyak parameter menjadi satu `*models.CareerPreferenceRequest`
   - Combine `job_type` dan `work_location` menjadi comma-separated string
   - Ensure profile exists sebelum create/update preference

3. **File:** `internal/handlers/preference_handler.go`
   - Simplified handler untuk pass request object langsung ke service

### **Test:**
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

---

## ✅ 3. Fix Position Preferences API (POST /api/v1/position-preferences)

### **Masalah:**
- JSON unmarshal error: expects array of objects, received array of strings

### **Request Format (Sebelum):**
```json
{
  "positions": [
    {"position_name": "Software Engineer", "priority": 1},
    {"position_name": "Backend Developer", "priority": 2}
  ]
}
```

### **Request Format (Setelah - FIXED):**
```json
{
  "positions": [
    "Software Engineer",
    "Backend Developer",
    "Full Stack Developer"
  ]
}
```

### **Perubahan:**
1. **File:** `internal/models/preference.go`
   - Update `BulkPositionPreferenceRequest.Positions` dari `[]PositionPreferenceRequest` ke `[]string`

2. **File:** `internal/services/preference_service.go`
   - Update method signature `CreatePositionPreferences(userID uint, positions []string)`
   - Auto-generate priority berdasarkan urutan array (index + 1)

3. **File:** `internal/handlers/preference_handler.go`
   - Pass `req.Positions` (array of strings) ke service

### **Test:**
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

---

## ✅ 4. Fix CV Upload API (POST /api/v1/cv)

### **Masalah:**
- Error "No file uploaded" meskipun file sudah di-upload
- Routing issue dengan trailing slash

### **Perubahan:**
1. **File:** `cmd/main.go`
   - Update CV routes untuk handle dengan dan tanpa trailing slash:
   ```go
   // Before (using route group)
   cv := v1.Group("/cv")
   {
     cv.POST("", cvHandler.Upload)
   }
   
   // After (direct routes)
   v1.POST("/cv", cvHandler.Upload)
   v1.POST("/cv/", cvHandler.Upload)
   ```

### **Test dengan cURL:**
```bash
# Test upload PDF
curl -X POST http://localhost/api/v1/cv \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/your/resume.pdf"

# Test upload DOCX
curl -X POST http://localhost/api/v1/cv \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/your/resume.docx"
```

### **Test dengan Postman:**
1. Method: POST
2. URL: `http://localhost/api/v1/cv`
3. Headers: `Authorization: Bearer YOUR_TOKEN`
4. Body: 
   - Type: form-data
   - Key: `file` (type: File)
   - Value: Select your PDF/DOCX file

---

## ✅ 5. Fix API Gateway Routing

### **Masalah:**
- Job Service tidak bisa diakses via API Gateway
- User Profile Service routing tidak optimal
- Company & Auth Service perlu dicek

### **Perubahan:**

#### **File:** `jobfair-api-gateway/cmd/main.go`
1. **Tambah Job Service Environment Variable:**
   ```go
   jobServiceURL := getEnv("JOB_SERVICE_URL", "http://localhost:8082")
   ```

2. **Tambah Job Service Proxy:**
   ```go
   jobProxy := createReverseProxy(jobServiceURL, "job-service")
   ```

3. **Tambah Job Service Routes:**
   ```go
   // Jobs
   router.Any("/api/v1/jobs/*proxyPath", proxyHandler(jobProxy, "/api/v1/jobs"))
   router.Any("/api/v1/jobs", proxyHandler(jobProxy, "/api/v1/jobs"))
   
   // Applications
   router.Any("/api/v1/applications/*proxyPath", proxyHandler(jobProxy, "/api/v1/applications"))
   router.Any("/api/v1/applications", proxyHandler(jobProxy, "/api/v1/applications"))
   ```

4. **Update All Service Routes** untuk handle dengan dan tanpa trailing slash

#### **File:** `docker-compose.yml`
1. Update API Gateway environment variable:
   ```yaml
   USER_PROFILE_SERVICE_URL: http://user-profile-service:8083
   JOB_SERVICE_URL: http://job-service:8082
   PORT: 8000
   ```

2. Update port mapping:
   ```yaml
   ports:
     - "80:8000"
   ```

### **Test Job Service via Gateway:**
```bash
# List jobs
curl http://localhost/api/v1/jobs?page=1&limit=10

# Get job detail
curl http://localhost/api/v1/jobs/1

# Create job (requires auth)
curl -X POST http://localhost/api/v1/jobs \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{...}'
```

### **Test Auth Service via Gateway:**
```bash
# Register
curl -X POST http://localhost/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123",
    "role": "jobseeker"
  }'

# Login
curl -X POST http://localhost/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123"
  }'
```

### **Test Company Service via Gateway:**
```bash
# Get companies
curl http://localhost/api/v1/companies?page=1&limit=10

# Get company detail
curl http://localhost/api/v1/companies/1
```

---

## 🚀 Cara Menjalankan Aplikasi

### **1. Start semua services dengan Docker Compose:**
```bash
# Build dan start semua services
docker-compose up --build

# Atau run di background
docker-compose up -d --build
```

### **2. Stop services:**
```bash
docker-compose down
```

### **3. Rebuild specific service:**
```bash
docker-compose up -d --build user-profile-service
docker-compose up -d --build api-gateway
```

### **4. View logs:**
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f user-profile-service
docker-compose logs -f api-gateway
```

---

## 🔍 Troubleshooting

### **Issue: Connection refused to service**
```bash
# Check if service is running
docker-compose ps

# Restart specific service
docker-compose restart user-profile-service

# Check service logs
docker-compose logs user-profile-service
```

### **Issue: Profile not found error**
```bash
# Pastikan user sudah membuat profile terlebih dahulu
curl -X POST http://localhost/api/v1/profiles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "phone_number": "081234567890"
  }'
```

### **Issue: Database migration error**
```bash
# Run migration manually
docker-compose up profile-migration

# Check database
docker exec -it postgres-profile psql -U jobfair_user -d jobfair_profiles
```

---

## 📝 Summary of Changes

| Issue | File(s) Changed | Status |
|-------|----------------|--------|
| Bulk Skills tidak menyimpan | `models/skill.go`, `services/skill_service.go` | ✅ Fixed |
| Career Preference foreign key error | `models/preference.go`, `services/preference_service.go`, `handlers/preference_handler.go` | ✅ Fixed |
| Position Preferences unmarshal error | `models/preference.go`, `services/preference_service.go` | ✅ Fixed |
| CV Upload "no file" error | `cmd/main.go` | ✅ Fixed |
| API Gateway routing | `jobfair-api-gateway/cmd/main.go`, `docker-compose.yml` | ✅ Fixed |

---

## 🎯 Testing Checklist

- [ ] Bulk skills API berhasil menyimpan data
- [ ] Career preference API berhasil create/update
- [ ] Position preferences API berhasil dengan array of strings
- [ ] CV upload berhasil untuk PDF dan DOCX
- [ ] Job Service dapat diakses via API Gateway
- [ ] Auth Service dapat diakses via API Gateway
- [ ] Company Service dapat diakses via API Gateway
- [ ] User Profile Service dapat diakses via API Gateway

---

## 📞 Support

Jika masih ada masalah, silakan cek:
1. Docker logs: `docker-compose logs -f [service-name]`
2. Database connection: pastikan semua database services healthy
3. RabbitMQ connection: `http://localhost:15672` (user: jobfair_user, pass: jobfair_pass)
