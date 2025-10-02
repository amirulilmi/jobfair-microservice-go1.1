# 🧪 Manual API Testing Guide

Panduan step-by-step untuk testing API secara manual.

## Prerequisites

1. ✅ Service berjalan di `http://localhost:8083`
2. ✅ Auth service berjalan di `http://localhost:8080`
3. ✅ Sudah punya user account (email & password)

---

## 🚀 Quick Start - Copy & Paste Commands

### Step 1: Get JWT Token

```bash
# Login untuk mendapatkan token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"your-email@example.com\",\"password\":\"your-password\"}"
```

**Copy token dari response**, simpan di variabel:

```bash
# Untuk bash/Linux/Mac:
export TOKEN="your-jwt-token-here"

# Untuk Windows CMD:
set TOKEN=your-jwt-token-here

# Untuk Windows PowerShell:
$TOKEN="your-jwt-token-here"
```

---

## 📝 Test Endpoints

### 1. Health Check ✅

```bash
curl http://localhost:8083/health
```

**Expected:**
```json
{
  "status": "ok",
  "service": "user-profile-service"
}
```

---

### 2. Create Profile 👤

```bash
curl -X POST http://localhost:8083/api/v1/profiles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"full_name\":\"John Doe\",\"phone_number\":\"081234567890\"}"
```

**Expected:**
```json
{
  "success": true,
  "message": "Profile created successfully",
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "full_name": "John Doe",
    "phone_number": "081234567890",
    "completion_status": 13
  }
}
```

---

### 3. Get Profile 📄

```bash
curl -X GET http://localhost:8083/api/v1/profiles \
  -H "Authorization: Bearer $TOKEN"
```

---

### 4. Update Profile ✏️

```bash
curl -X PUT http://localhost:8083/api/v1/profiles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"bio\":\"Experienced Software Engineer\",\"city\":\"Jakarta\",\"province\":\"DKI Jakarta\"}"
```

---

### 5. Add Work Experience 💼

```bash
curl -X POST http://localhost:8083/api/v1/work-experiences \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"company_name\":\"PT Bumi Resources\",\"job_position\":\"Senior Engineer\",\"start_date\":\"2020-01-01T00:00:00Z\",\"is_current_job\":true,\"job_description\":\"Leading engineering team\"}"
```

---

### 6. Get Work Experiences 📋

```bash
curl -X GET http://localhost:8083/api/v1/work-experiences \
  -H "Authorization: Bearer $TOKEN"
```

---

### 7. Add Education 🎓

```bash
curl -X POST http://localhost:8083/api/v1/educations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"university\":\"Institut Teknologi Bandung\",\"major\":\"Teknik Pertambangan\",\"degree\":\"Bachelor\",\"start_date\":\"2014-08-01T00:00:00Z\",\"end_date\":\"2018-07-31T00:00:00Z\",\"gpa\":3.98}"
```

---

### 8. Get Educations 📚

```bash
curl -X GET http://localhost:8083/api/v1/educations \
  -H "Authorization: Bearer $TOKEN"
```

---

### 9. Add Certification 🏆

```bash
curl -X POST http://localhost:8083/api/v1/certifications \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"certification_name\":\"Certified Mine Manager\",\"organizer\":\"Ministry of Energy\",\"issue_date\":\"2022-03-15T00:00:00Z\",\"credential_id\":\"CMM-2022-12345\"}"
```

---

### 10. Get Certifications 📜

```bash
curl -X GET http://localhost:8083/api/v1/certifications \
  -H "Authorization: Bearer $TOKEN"
```

---

### 11. Add Skills (Bulk) 💪

```bash
curl -X POST http://localhost:8083/api/v1/skills/bulk \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"technical_skills\":[{\"skill_name\":\"Mine Planning\",\"skill_type\":\"technical\",\"proficiency_level\":\"expert\",\"years_of_experience\":5},{\"skill_name\":\"Drilling & Blasting\",\"skill_type\":\"technical\",\"proficiency_level\":\"advanced\",\"years_of_experience\":4}],\"soft_skills\":[{\"skill_name\":\"Leadership\",\"skill_type\":\"soft\",\"proficiency_level\":\"advanced\"},{\"skill_name\":\"Problem Solving\",\"skill_type\":\"soft\",\"proficiency_level\":\"expert\"}]}"
```

---

### 12. Get Skills ⚡

```bash
curl -X GET http://localhost:8083/api/v1/skills \
  -H "Authorization: Bearer $TOKEN"
```

---

### 13. Set Career Preference 🎯

```bash
curl -X POST http://localhost:8083/api/v1/career-preference \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"is_actively_looking\":true,\"expected_salary_min\":15000000,\"expected_salary_max\":25000000,\"salary_currency\":\"IDR\",\"is_negotiable\":true,\"preferred_work_types\":[\"onsite\",\"hybrid\"],\"preferred_locations\":[\"Jakarta\",\"Bandung\"],\"willing_to_relocate\":false}"
```

---

### 14. Get Career Preference 🔍

```bash
curl -X GET http://localhost:8083/api/v1/career-preference \
  -H "Authorization: Bearer $TOKEN"
```

---

### 15. Set Position Preferences 📍

```bash
curl -X POST http://localhost:8083/api/v1/position-preferences \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"positions\":[{\"position_name\":\"Senior Mining Engineer\",\"priority\":1},{\"position_name\":\"Mine Operations Manager\",\"priority\":2},{\"position_name\":\"Health & Safety Supervisor\",\"priority\":3}]}"
```

---

### 16. Get Position Preferences 📌

```bash
curl -X GET http://localhost:8083/api/v1/position-preferences \
  -H "Authorization: Bearer $TOKEN"
```

---

### 17. Upload CV 📄

```bash
curl -X POST http://localhost:8083/api/v1/cv \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@C:/path/to/your/cv.pdf"
```

**Note:** Ganti path dengan lokasi file CV Anda.

---

### 18. Get CV Info 📥

```bash
curl -X GET http://localhost:8083/api/v1/cv \
  -H "Authorization: Bearer $TOKEN"
```

---

### 19. Get Complete Profile 🎊

```bash
curl -X GET http://localhost:8083/api/v1/profiles/full \
  -H "Authorization: Bearer $TOKEN"
```

**Response:** Profile lengkap dengan semua relasi (work exp, education, skills, dll)

---

### 20. Check Profile Completion 📊

```bash
curl -X GET http://localhost:8083/api/v1/profiles/completion \
  -H "Authorization: Bearer $TOKEN"
```

**Expected:**
```json
{
  "success": true,
  "message": "Completion status retrieved",
  "data": {
    "completion_status": 93,
    "is_complete": false
  }
}
```

---

## 🎯 Testing Sequence

Ikuti urutan ini untuk test lengkap:

1. ✅ Health Check
2. ✅ Get Token (login)
3. ✅ Create Profile
4. ✅ Update Profile
5. ✅ Add Work Experience
6. ✅ Add Education
7. ✅ Add Certification
8. ✅ Add Skills (bulk)
9. ✅ Set Career Preference
10. ✅ Set Position Preferences
11. ✅ Upload CV
12. ✅ Get Complete Profile
13. ✅ Check Completion Status (should be ~90-100%)

---

## 🔧 Troubleshooting

### Error: "Unauthorized"

```bash
# Check token format
echo $TOKEN

# Token harus dimulai dengan: eyJ...
# Jika kosong atau salah, login ulang
```

### Error: "Profile not found"

```bash
# Create profile dulu
curl -X POST http://localhost:8083/api/v1/profiles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"full_name\":\"Your Name\",\"phone_number\":\"08123456789\"}"
```

### Error: "Database connection failed"

```bash
# Check PostgreSQL running
psql -U postgres -h localhost

# Check migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/jobfair_user_profile?sslmode=disable" up
```

---

## 💡 Tips

1. **Gunakan Postman/Thunder Client** untuk testing yang lebih visual
2. **Save token** dalam environment variable untuk kemudahan
3. **Test satu-satu** endpoint untuk memahami response
4. **Check completion status** setelah menambah data untuk lihat progress

---

## 📱 Testing dari Mobile/Frontend

Format request sama persis, hanya perlu:

1. Set header `Authorization: Bearer <token>`
2. Set header `Content-Type: application/json`
3. Parse JSON response

**Example (JavaScript/React):**

```javascript
const response = await fetch('http://localhost:8083/api/v1/profiles', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});

const data = await response.json();
console.log(data);
```

---

**Happy Testing! 🚀**
