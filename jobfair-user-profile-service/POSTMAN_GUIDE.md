# 📮 Postman Collection Guide

Panduan lengkap cara import dan menggunakan Postman Collection untuk testing API.

## 📁 Files yang Tersedia

1. **`postman_collection.json`** - Collection dengan 35+ API requests
2. **`postman_environment.json`** - Environment variables (token, base_url, etc.)

---

## 🚀 Quick Start - Import ke Postman

### Step 1: Install Postman

Download dan install:
- **Desktop:** https://www.postman.com/downloads/
- **Web:** https://web.postman.com/

### Step 2: Import Collection

1. Buka Postman
2. Click **"Import"** button (top left)
3. Click **"Choose Files"**
4. Select `postman_collection.json`
5. Click **"Import"**

✅ Collection "JobFair - User Profile Service" akan muncul di sidebar!

### Step 3: Import Environment

1. Click **"Import"** button lagi
2. Select `postman_environment.json`
3. Click **"Import"**

✅ Environment "JobFair - User Profile Service (Local)" akan muncul!

### Step 4: Set Active Environment

1. Click dropdown di top-right corner (shows "No Environment")
2. Select **"JobFair - User Profile Service (Local)"**

✅ Environment aktif!

---

## 🔑 Setup JWT Token

### Option 1: Get Token from Auth Service

**Menggunakan Postman:**

1. Buat request baru atau gunakan auth service collection
2. POST `http://localhost:8080/api/v1/auth/login`
3. Body (JSON):
   ```json
   {
     "email": "test@example.com",
     "password": "password123"
   }
   ```
4. Send request
5. Copy `token` dari response
6. Paste ke environment variable `token`

**Set Token di Environment:**

1. Click **"Environments"** (left sidebar)
2. Click **"JobFair - User Profile Service (Local)"**
3. Paste token ke kolom **"Current Value"** dari variable `token`
4. Click **"Save"** (Ctrl+S)

### Option 2: Get Token via curl (Copy ke Postman)

```bash
# Login via curl
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"test@example.com\",\"password\":\"password123\"}"

# Copy token dari response
# Paste ke Postman environment variable "token"
```

✅ **Token sudah siap digunakan!**

---

## 🎯 Testing dengan Postman

### Test Sequence (Recommended Order)

1. **Health Check** ✅ (No auth needed)
2. **Create Profile** ✅
3. **Update Profile** ✅
4. **Add Work Experience** ✅
5. **Add Education** ✅
6. **Add Certification** ✅
7. **Add Skills (Bulk)** ✅
8. **Set Career Preference** ✅
9. **Set Position Preferences** ✅
10. **Upload CV** ✅
11. **Get Complete Profile** ✅
12. **Check Completion Status** ✅

### How to Run a Request

1. **Expand folder** (e.g., "1. Profile Management")
2. **Click request** (e.g., "Create Profile")
3. **Review request** (Method, URL, Headers, Body)
4. **Click "Send"** button
5. **See response** in bottom panel

### Example: Create Profile

1. Expand **"1. Profile Management"**
2. Click **"Create Profile"**
3. Lihat Body (sudah ada example JSON):
   ```json
   {
     "full_name": "John Doe",
     "phone_number": "081234567890"
   }
   ```
4. Click **"Send"**
5. Response akan muncul:
   ```json
   {
     "success": true,
     "message": "Profile created successfully",
     "data": {
       "id": "uuid-here",
       "full_name": "John Doe",
       "completion_status": 13
     }
   }
   ```

---

## 📊 Collection Structure

```
JobFair - User Profile Service/
├── 0. Health Check
│   └── Health Check (No auth)
│
├── 1. Profile Management (5 requests)
│   ├── Create Profile
│   ├── Get Profile
│   ├── Get Complete Profile (with relations)
│   ├── Update Profile
│   └── Get Profile Completion Status
│
├── 2. Work Experience (5 requests)
│   ├── Create Work Experience
│   ├── Get All Work Experiences
│   ├── Get Work Experience by ID
│   ├── Update Work Experience
│   └── Delete Work Experience
│
├── 3. Education (5 requests)
│   ├── Create Education
│   ├── Get All Educations
│   ├── Get Education by ID
│   ├── Update Education
│   └── Delete Education
│
├── 4. Certifications (5 requests)
│   ├── Create Certification
│   ├── Get All Certifications
│   ├── Get Certification by ID
│   ├── Update Certification
│   └── Delete Certification
│
├── 5. Skills (6 requests)
│   ├── Create Single Skill
│   ├── Create Multiple Skills (Bulk)
│   ├── Get All Skills
│   ├── Get Skill by ID
│   ├── Update Skill
│   └── Delete Skill
│
├── 6. Career Preference (2 requests)
│   ├── Create/Update Career Preference
│   └── Get Career Preference
│
├── 7. Position Preferences (3 requests)
│   ├── Create Position Preferences (Bulk)
│   ├── Get Position Preferences
│   └── Delete Position Preference
│
├── 8. CV Upload (3 requests)
│   ├── Upload CV
│   ├── Get CV Info
│   └── Delete CV
│
└── 9. Badges (1 request)
    └── Get User Badges
```

**Total: 35+ requests**

---

## 💡 Tips & Tricks

### 1. Auto-Set Variables from Response

**Extract ID dari response:**

Di tab **"Tests"** pada request, tambahkan:

```javascript
// Auto-save work experience ID
if (pm.response.code === 200 || pm.response.code === 201) {
    var jsonData = pm.response.json();
    if (jsonData.data && jsonData.data.id) {
        pm.environment.set("work_exp_id", jsonData.data.id);
    }
}
```

Sekarang ID otomatis tersimpan untuk request selanjutnya!

### 2. Run Collection

**Run semua requests sekaligus:**

1. Right-click collection → **"Run collection"**
2. Select requests yang mau di-run
3. Click **"Run JobFair - User Profile Service"**
4. Lihat summary report

### 3. Save Response Examples

**Save successful responses:**

1. Send request
2. Click **"Save as Example"** (di tab Response)
3. Example tersimpan untuk referensi

### 4. Use Pre-request Script

**Check token before request:**

Di tab **"Pre-request Script"** pada collection level:

```javascript
// Check if token exists
if (!pm.environment.get("token")) {
    console.log("⚠️ Token not set! Please login first.");
}
```

### 5. Share with Team

**Export dan share:**

1. Right-click collection → **"Export"**
2. Choose format: Collection v2.1
3. Share JSON file dengan tim

---

## 🔧 Troubleshooting

### Issue 1: "Unauthorized" Error

**Solution:**
- Check token di environment variable
- Token format harus: `eyJ...`
- Token mungkin expired, login ulang

### Issue 2: "Could not get response"

**Solution:**
- Check service running: `curl http://localhost:8083/health`
- Check firewall/antivirus
- Check base_url correct

### Issue 3: "404 Not Found"

**Solution:**
- Check endpoint URL correct
- Check base_url di environment
- Service mungkin tidak running

### Issue 4: CV Upload Failed

**Solution:**
- Pilih file di Body → form-data → file field
- File max 5MB
- Format: PDF, DOC, atau DOCX

---

## 📝 Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `base_url` | API base URL | `http://localhost:8083` |
| `auth_service_url` | Auth service URL | `http://localhost:8080` |
| `token` | JWT token | `eyJhbGc...` |
| `work_exp_id` | Work experience ID | `uuid` |
| `education_id` | Education ID | `uuid` |
| `cert_id` | Certification ID | `uuid` |
| `skill_id` | Skill ID | `uuid` |
| `position_pref_id` | Position preference ID | `uuid` |

**Cara edit variables:**

1. Click **"Environments"** (left sidebar)
2. Select environment
3. Edit **"Current Value"**
4. Click **"Save"**

---

## 🎨 Customize Collection

### Edit Request Body

1. Click request
2. Go to **"Body"** tab
3. Edit JSON
4. Click **"Send"**

### Add New Request

1. Right-click folder
2. Click **"Add Request"**
3. Set method, URL, body
4. Click **"Save"**

### Duplicate Request

1. Right-click request
2. Click **"Duplicate"**
3. Modify as needed

---

## 🚀 Advanced Features

### 1. Collection Runner

**Run all tests automatically:**

```javascript
// In Tests tab
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has success field", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.success).to.eql(true);
});
```

### 2. Newman (CLI Runner)

**Run collection from command line:**

```bash
# Install Newman
npm install -g newman

# Run collection
newman run postman_collection.json \
  -e postman_environment.json \
  --reporters cli,json
```

### 3. Monitor (Scheduled Tests)

1. Click **"Monitors"** (left sidebar)
2. Click **"Create a Monitor"**
3. Select collection
4. Set schedule (e.g., every hour)
5. Get notifications on failures

---

## 📊 Response Examples

### Successful Response

```json
{
  "success": true,
  "message": "Profile created successfully",
  "data": {
    "id": "uuid-here",
    "user_id": "uuid-here",
    "full_name": "John Doe",
    "phone_number": "081234567890",
    "completion_status": 13,
    "created_at": "2025-10-01T...",
    "updated_at": "2025-10-01T..."
  }
}
```

### Error Response

```json
{
  "success": false,
  "message": "Invalid request",
  "code": "VALIDATION_ERROR",
  "error": "full_name is required"
}
```

---

## 🎯 Testing Checklist

Use this checklist untuk ensure semua endpoints work:

### Profile Management
- [ ] Create Profile
- [ ] Get Profile
- [ ] Update Profile
- [ ] Get Complete Profile
- [ ] Get Completion Status

### Work Experience
- [ ] Create Work Experience
- [ ] Get All Work Experiences
- [ ] Update Work Experience
- [ ] Delete Work Experience

### Education
- [ ] Create Education
- [ ] Get All Educations
- [ ] Update Education
- [ ] Delete Education

### Certifications
- [ ] Create Certification
- [ ] Get All Certifications
- [ ] Update Certification
- [ ] Delete Certification

### Skills
- [ ] Create Skills (Bulk)
- [ ] Get All Skills
- [ ] Update Skill
- [ ] Delete Skill

### Preferences
- [ ] Set Career Preference
- [ ] Get Career Preference
- [ ] Set Position Preferences
- [ ] Get Position Preferences

### CV Upload
- [ ] Upload CV
- [ ] Get CV Info
- [ ] Delete CV

---

## 🆘 Need Help?

**Resources:**
- Postman Documentation: https://learning.postman.com/
- API Documentation: See `API_COLLECTION.md`
- Manual Testing: See `MANUAL_TESTING.md`

**Common Issues:**
- Service not running → Start service
- Token expired → Login again
- Wrong URL → Check environment variables

---

## 🎊 Summary

### ✅ What You Get:

1. **35+ ready-to-use API requests**
2. **Organized in folders** by feature
3. **Environment variables** for easy management
4. **Example bodies** for all requests
5. **Descriptions** for each endpoint

### 🚀 How to Use:

1. Import collection & environment
2. Set token (login first)
3. Click request → Send
4. See response

### 💡 Pro Tips:

- Use Collection Runner for bulk testing
- Save response examples
- Use pre-request scripts
- Share with team

---

**Happy Testing with Postman! 🎉**
