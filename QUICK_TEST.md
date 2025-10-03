# 🚀 Quick Test Guide - CV Upload

## 📝 Prerequisites
- Docker services running: `docker-compose ps`
- Valid JWT token (get from login)

---

## ⚡ Quick Test (Copy & Paste)

### 1️⃣ Get JWT Token
```bash
export TOKEN=$(curl -s -X POST http://localhost/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jobseeker@example.com",
    "password": "YourPassword123"
  }' | jq -r '.data.access_token')

echo "Token: ${TOKEN:0:30}..."
```

### 2️⃣ Upload CV (Auto Test)
```bash
# Bash/Linux/Mac
./test-upload-cv.sh "$TOKEN"

# PowerShell/Windows
.\test-upload-cv.ps1 -Token "$TOKEN"
```

### 3️⃣ Upload CV (Manual)
```bash
curl -X POST http://localhost/api/v1/cv \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/your/cv.pdf" \
  | jq
```

---

## 🔧 Common Commands

### Check Services
```bash
docker-compose ps                           # List all services
docker-compose logs -f api-gateway          # View gateway logs
docker-compose logs -f user-profile-service # View profile service logs
```

### Restart Services
```bash
docker-compose restart api-gateway user-profile-service
```

### Rebuild After Code Changes
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

---

## 🎯 Test Different File Types

```bash
# PDF (recommended)
curl -H "Authorization: Bearer $TOKEN" \
  -F "file=@cv.pdf" \
  http://localhost/api/v1/cv

# Word Document
curl -H "Authorization: Bearer $TOKEN" \
  -F "file=@cv.docx" \
  http://localhost/api/v1/cv

# DOC (old Word format)
curl -H "Authorization: Bearer $TOKEN" \
  -F "file=@cv.doc" \
  http://localhost/api/v1/cv
```

---

## ✅ Expected Success Response

```json
{
  "success": true,
  "message": "CV uploaded successfully",
  "data": {
    "id": 1,
    "profile_id": 1,
    "file_name": "cv.pdf",
    "file_size": 123456,
    "file_type": ".pdf",
    "is_verified": false
  }
}
```

---

## ❌ Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| "No file uploaded" | Wrong field name or missing file | Use `-F "file=@..."` |
| "Unauthorized" | Invalid/expired token | Get new token via login |
| "File type not allowed" | Unsupported format | Use .pdf, .doc, or .docx |
| "Profile not found" | No profile exists | Create profile first |

---

## 🔄 Complete Workflow

```bash
# 1. Login
export TOKEN=$(curl -s -X POST http://localhost/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"pass123"}' \
  | jq -r '.data.access_token')

# 2. Check Profile
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost/api/v1/profiles | jq

# 3. Upload CV
curl -s -H "Authorization: Bearer $TOKEN" \
  -F "file=@cv.pdf" \
  http://localhost/api/v1/cv | jq

# 4. Get CV Info
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost/api/v1/cv | jq

# 5. Delete CV (optional)
curl -s -X DELETE \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost/api/v1/cv | jq
```

---

## 🐛 Debug Commands

```bash
# Check if API Gateway is accessible
curl http://localhost/health

# Test specific service directly (bypass gateway)
curl http://localhost:8083/health

# View recent logs
docker-compose logs --tail=50 user-profile-service

# Check upload directory
docker exec jobfair-user-profile-service ls -la /root/uploads/cv
```

---

## 📱 Using Postman

1. **Request Type**: POST
2. **URL**: `http://localhost/api/v1/cv`
3. **Headers**:
   - Authorization: `Bearer YOUR_TOKEN`
4. **Body** → form-data:
   - Key: `file` (type: File)
   - Value: [Select your CV file]

---

## 💡 Pro Tips

✅ Always use lowercase `file` as field name  
✅ Use `@` before file path in curl  
✅ Keep CV files under 5MB  
✅ Test with small files first  
✅ Check logs when debugging  

---

## 📞 Need Help?

- Read full guide: `CV_UPLOAD_GUIDE.md`
- Check fix summary: `CV_UPLOAD_FIX_SUMMARY.md`
- View service logs for errors
