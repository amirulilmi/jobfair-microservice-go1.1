# 📄 CV Upload Troubleshooting Guide

## 🔧 Common Issues & Solutions

### ❌ Issue: "No file uploaded" Error

**Symptoms:**
```json
{
    "success": false,
    "message": "No file uploaded",
    "code": "NO_FILE"
}
```

**Root Causes & Solutions:**

#### 1️⃣ **Incorrect Form Field Name**
The API expects the field name to be `file` (lowercase).

✅ **Correct:**
```bash
curl -X POST http://localhost/api/v1/cv \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/cv.pdf"
```

❌ **Incorrect:**
```bash
-F "cv=@/path/to/cv.pdf"      # Wrong field name
-F "document=@/path/to/cv.pdf" # Wrong field name
-F "File=@/path/to/cv.pdf"     # Wrong case (uppercase F)
```

#### 2️⃣ **Missing or Incorrect Content-Type**
The request must be `multipart/form-data`.

✅ **Correct (curl does this automatically with -F):**
```bash
curl -F "file=@cv.pdf" ...
```

❌ **Incorrect:**
```bash
curl -H "Content-Type: application/json" -d '{"file":"..."}' ...
```

#### 3️⃣ **File Path Issues**
Make sure the file exists and path is correct.

```bash
# Check if file exists
ls -la /path/to/cv.pdf

# Use absolute path
curl -F "file=@/Users/username/Documents/cv.pdf" ...

# Or relative path from current directory
curl -F "file=@./cv.pdf" ...
```

#### 4️⃣ **Authentication Token Missing**
```bash
# ✅ Correct
curl -H "Authorization: Bearer eyJhbGc..." -F "file=@cv.pdf" ...

# ❌ Missing
curl -F "file=@cv.pdf" ...
```

---

## 🧪 Testing Methods

### Method 1: Using Test Scripts (Recommended)

**Bash (Linux/Mac):**
```bash
chmod +x test-upload-cv.sh
./test-upload-cv.sh "YOUR_JWT_TOKEN"
```

**PowerShell (Windows):**
```powershell
.\test-upload-cv.ps1 -Token "YOUR_JWT_TOKEN"
```

### Method 2: Using cURL

**Step 1: Get your JWT token**
```bash
# Login first
TOKEN=$(curl -X POST http://localhost/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }' | jq -r '.data.access_token')

echo $TOKEN
```

**Step 2: Upload CV**
```bash
curl -X POST http://localhost/api/v1/cv \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/your/cv.pdf" \
  | jq
```

### Method 3: Using Postman

1. **Set Request Type:** `POST`
2. **URL:** `http://localhost/api/v1/cv`
3. **Headers:**
   - Key: `Authorization`
   - Value: `Bearer YOUR_JWT_TOKEN`
4. **Body:**
   - Select `form-data`
   - Add key: `file` (change type to `File`)
   - Click "Select Files" and choose your CV

### Method 4: Using HTTPie (if installed)

```bash
http -f POST http://localhost/api/v1/cv \
  Authorization:"Bearer YOUR_TOKEN" \
  file@/path/to/cv.pdf
```

---

## 📋 Supported File Types & Size Limits

### Allowed File Types
- `.pdf` - PDF documents
- `.doc` - Microsoft Word (old format)
- `.docx` - Microsoft Word (new format)

### File Size Limits
- **Maximum:** 10MB (configurable)
- **Recommended:** Under 5MB for best performance

---

## 🔍 Debugging Steps

### 1. Check Service Logs
```bash
# View profile service logs
docker logs jobfair-user-profile-service -f

# View API gateway logs
docker logs jobfair-api-gateway -f
```

### 2. Verify Token
```bash
# Decode JWT to check if it's valid
echo "YOUR_TOKEN" | cut -d. -f2 | base64 -d | jq
```

### 3. Test Health Endpoint
```bash
# Check if service is running
curl http://localhost/api/v1/profiles/health
```

### 4. Check Upload Directory Permissions
```bash
# Inside the container
docker exec -it jobfair-user-profile-service ls -la /root/uploads
```

---

## 🐛 Common Error Messages

### Error: "Unauthorized"
**Cause:** Missing or invalid JWT token
**Solution:** Login again to get a fresh token

### Error: "File type not allowed"
**Cause:** Trying to upload unsupported file format
**Solution:** Convert file to .pdf, .doc, or .docx

### Error: "File size exceeds maximum"
**Cause:** File is too large
**Solution:** Compress the PDF or reduce file size

### Error: "Profile not found"
**Cause:** User profile doesn't exist yet
**Solution:** Create profile first via `/api/v1/profiles`

---

## ✅ Complete Working Example

```bash
#!/bin/bash

# Step 1: Register & Login
echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jobseeker@example.com",
    "password": "SecurePass123!"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ Login failed"
    exit 1
fi

echo "✅ Login successful"
echo "Token: ${TOKEN:0:20}..."

# Step 2: Create/Check Profile
echo -e "\nChecking profile..."
curl -s -X GET http://localhost/api/v1/profiles \
  -H "Authorization: Bearer $TOKEN" | jq

# Step 3: Upload CV
echo -e "\nUploading CV..."
UPLOAD_RESPONSE=$(curl -s -X POST http://localhost/api/v1/cv \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@./my-cv.pdf")

echo $UPLOAD_RESPONSE | jq

# Check if upload was successful
SUCCESS=$(echo $UPLOAD_RESPONSE | jq -r '.success')

if [ "$SUCCESS" = "true" ]; then
    echo "✅ CV uploaded successfully!"
else
    echo "❌ Upload failed"
    echo $UPLOAD_RESPONSE | jq
fi
```

---

## 🔄 API Endpoints Reference

### Upload CV
```
POST /api/v1/cv
Authorization: Bearer {token}
Content-Type: multipart/form-data

Body:
  file: [PDF/DOC/DOCX file]
```

### Get Current CV
```
GET /api/v1/cv
Authorization: Bearer {token}
```

### Delete CV
```
DELETE /api/v1/cv
Authorization: Bearer {token}
```

---

## 💡 Tips for Success

1. **Always use the `-F` flag with curl** for file uploads (not `-d`)
2. **Field name must be exactly `file`** (lowercase)
3. **Include `@` before the file path** in curl
4. **Check file exists** before uploading
5. **Use fresh JWT tokens** (they expire)
6. **Monitor service logs** when debugging
7. **Test with small files first** (< 1MB)

---

## 📞 Still Having Issues?

If you're still experiencing problems:

1. Check the logs: `docker logs jobfair-user-profile-service -f`
2. Verify API Gateway is running: `docker ps | grep api-gateway`
3. Test with the provided test scripts
4. Create an issue with:
   - Full curl command you're using
   - Complete error response
   - Service logs
   - File type and size

---

## 🎯 Quick Test

Run this one-liner to test if CV upload is working:

```bash
# Create a tiny test PDF and upload it
echo '%PDF-1.4
%%EOF' > test.pdf && \
curl -X POST http://localhost/api/v1/cv \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@test.pdf" && \
rm test.pdf
```

Replace `YOUR_TOKEN` with your actual JWT token.
