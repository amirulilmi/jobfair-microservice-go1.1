# 📮 JobFair Microservices - Postman Collection

Koleksi lengkap API untuk JobFair Microservices yang mencakup **Auth Service, Company Service, Job Service, dan User Profile Service**.

## 📥 Cara Import ke Postman

### 1. Import Collection & Environment

1. Buka Postman
2. Klik **Import** (tombol di kiri atas)
3. Drag & drop atau pilih kedua file ini:
   - `JobFair_Microservices.postman_collection.json`
   - `JobFair_Microservices.postman_environment.json`
4. Klik **Import**

### 2. Set Environment

1. Di pojok kanan atas Postman, pilih dropdown environment
2. Pilih **"JobFair Microservices - Local"**
3. Environment sudah siap digunakan!

## 🚀 Quick Start Guide

### Flow untuk Job Seeker (Complete)

#### 1️⃣ **Register & Login**
```
1. Auth > Register - Step 1 (Init) [job_seeker]
2. Auth > Complete Profile (Job Seeker)
3. Auth > Send OTP
4. Auth > Verify OTP
5. Auth > Set Employment Status
6. Auth > Set Job Preferences
7. Auth > Upload Profile Photo (optional)
```

#### 2️⃣ **Build Detailed Profile**
```
1. User Profile > Profile > Create Profile
2. User Profile > Work Experience > Add Work Experience
3. User Profile > Education > Add Education
4. User Profile > Certifications > Add Certification
5. User Profile > Skills > Add Multiple Skills (Bulk)
6. User Profile > Career Preferences > Create/Update
7. User Profile > Position Preferences > Create
8. User Profile > CV Management > Upload CV
9. User Profile > Profile > Get Profile Completion Status
```

#### 3️⃣ **Browse & Apply Jobs**
```
1. Job Service > List Jobs (Public)
2. Job Service > Get Job by ID
3. Job Service > Apply to Job
4. Job Service > Get My Applications
```

#### 4️⃣ **Save Favorite Jobs**
```
1. Job Service > Save Job (Bookmark)
2. Job Service > Get Saved Jobs
```

### Flow untuk Company

#### 1️⃣ **Register & Login**
```
1. Auth > Register - Step 1 (Company)
2. Auth > Complete Profile (Company)
3. Auth > Send OTP
4. Auth > Verify OTP
5. Auth > Upload Profile Photo (Logo)
```

#### 2️⃣ **Setup Company Profile**
```
1. Company Service > Get My Company
2. Company Service > Create Company (if not exists)
3. Company Service > Update Company
4. Company Service > Upload Company Logo
5. Company Service > Get Company Analytics
```

#### 3️⃣ **Post & Manage Jobs**
```
1. Job Service > Create Job
2. Job Service > Publish Job
3. Job Service > Get My Jobs
4. Job Service > Get Applications by Job ID
5. Job Service > Update Application Status
```

## 🔑 Authentication

Collection ini sudah dilengkapi dengan **auto-save tokens**!

- Setelah Login/Register berhasil, `access_token` dan `refresh_token` otomatis tersimpan
- Token otomatis digunakan di semua request yang memerlukan authentication
- Jika token expired, gunakan endpoint **"Refresh Token"**

## 📂 Struktur Collection

```
📮 JobFair Microservices API
├── 🔐 Auth Service (Port 8080) - 12 endpoints
│   ├── Health Check
│   ├── Registration Flow (7 steps)
│   ├── Login
│   ├── Refresh Token
│   └── Get All Users (Debug)
│
├── 🏢 Company Service (Port 8081) - 8 endpoints
│   ├── Health Check
│   ├── Public Routes (List/Get Companies)
│   ├── Company Management (CRUD)
│   ├── File Uploads (Logo)
│   └── Analytics
│
├── 💼 Job Service (Port 8082) - 8 endpoints
│   ├── Health Check
│   ├── Public Routes (List/Get Jobs)
│   ├── Job Management (Create, Publish)
│   ├── Job Applications (Apply, Get My Applications)
│   └── Saved Jobs (Bookmark)
│
└── 👤 User Profile Service (Port 8083) - 35+ endpoints
    ├── Health Check
    ├── 📄 Profile Management (5 endpoints)
    │   ├── Create/Get/Update Profile
    │   ├── Get Full Profile with All Relations
    │   └── Get Profile Completion Status
    │
    ├── 💼 Work Experience (5 endpoints)
    │   └── Full CRUD operations
    │
    ├── 🎓 Education (4 endpoints)
    │   └── Full CRUD operations
    │
    ├── 🏆 Certifications (4 endpoints)
    │   └── Full CRUD operations
    │
    ├── 🔧 Skills (5 endpoints)
    │   ├── Add Single/Bulk Skills
    │   └── Full CRUD operations
    │
    ├── ⚙️ Career Preferences (2 endpoints)
    │   └── Create/Update & Get
    │
    ├── 🎯 Position Preferences (3 endpoints)
    │   └── Create, Get, Delete
    │
    └── 📄 CV Management (3 endpoints)
        └── Upload, Get, Delete CV
```

## 🔧 Environment Variables

Environment ini sudah include variable berikut:

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `auth_url` | http://localhost:8080 | Auth Service URL |
| `company_url` | http://localhost:8081 | Company Service URL |
| `job_url` | http://localhost:8082 | Job Service URL |
| `profile_url` | http://localhost:8083 | User Profile Service URL |
| `gateway_url` | http://localhost | API Gateway URL |
| `access_token` | (auto-saved) | JWT Access Token |
| `refresh_token` | (auto-saved) | JWT Refresh Token |
| `user_id` | (auto-saved) | Current User ID |
| `company_id` | (auto-saved) | Company ID |
| `job_id` | (auto-saved) | Job ID |
| `profile_id` | (auto-saved) | Profile ID |
| `work_exp_id` | (auto-saved) | Work Experience ID |
| `education_id` | (auto-saved) | Education ID |
| `cert_id` | (auto-saved) | Certification ID |
| `skill_id` | (auto-saved) | Skill ID |

## 📝 Tips Penggunaan

### 1. Auto-Save IDs
Collection ini dilengkapi dengan **Test Scripts** yang otomatis menyimpan ID penting setelah create operation.

### 2. File Upload
Untuk endpoint yang memerlukan file upload:
1. Pilih tab **Body** 
2. Pilih **form-data**
3. Klik dropdown di kolom **Key**, pilih **File**
4. Klik **Select Files** untuk upload file

**CV Upload Requirements:**
- Format: PDF only
- Max size: 5MB
- Endpoint: `POST /api/v1/cv`

### 3. Bulk Operations
User Profile Service mendukung bulk operations untuk Skills:
- Gunakan endpoint **"Add Multiple Skills (Bulk)"**
- Kirim array skills dalam satu request
- Lebih efisien daripada add satu per satu

### 4. Profile Completion Tracking
Gunakan endpoint **"Get Profile Completion Status"** untuk:
- Check progress profile completion
- Identifikasi section yang belum dilengkapi
- Meningkatkan profile visibility

### 5. Testing Different User Types
Untuk test flow berbeda:
1. **Job Seeker**: Login sebagai job seeker
2. **Company**: Login sebagai company
3. Clear token antar test untuk switch user

## 🎯 User Profile Service Features

### Profile Management
- ✅ Create complete professional profile
- ✅ Track profile completion percentage
- ✅ Get profile with all relations in one call
- ✅ Update profile information

### Work Experience
- ✅ Add multiple work experiences
- ✅ Mark current position
- ✅ Include employment type & dates
- ✅ Detailed descriptions

### Education
- ✅ Academic background
- ✅ Degree & field of study
- ✅ Grade/GPA tracking
- ✅ Duration tracking

### Certifications
- ✅ Professional certifications
- ✅ Credential ID & URL
- ✅ Issue & expiry dates
- ✅ Verification links

### Skills
- ✅ Add single or bulk skills
- ✅ Proficiency levels (beginner/intermediate/advanced/expert)
- ✅ Years of experience per skill
- ✅ Comprehensive skill portfolio

### Career Preferences
- ✅ Job type preferences
- ✅ Work location (remote/onsite/hybrid)
- ✅ Salary expectations
- ✅ Relocation willingness
- ✅ Availability date

### Position Preferences
- ✅ Desired job titles
- ✅ Multiple position preferences
- ✅ Easy management

### CV Management
- ✅ Upload PDF CV
- ✅ View CV information
- ✅ Update/replace CV
- ✅ Delete CV

## 🐛 Troubleshooting

### Token Expired
```
Error: 401 Unauthorized
Solution: Gunakan endpoint "Refresh Token"
```

### File Upload Failed
```
Error: File too large / Invalid format
Solution: 
- CV: Max 5MB, PDF only
- Check file format & size before upload
```

### Profile Not Found
```
Error: 404 Profile not found
Solution: Create profile first using "Create Profile" endpoint
```

### Permission Denied
```
Error: 403 Forbidden - Only companies can...
Solution: Pastikan login dengan user type yang benar
```

### Server Not Running
```
Error: Could not get any response / ECONNREFUSED
Solution: 
1. docker compose ps (check all services running)
2. docker compose logs -f [service-name]
3. Check port tidak bentrok
```

## 🎯 Testing Checklist

### Job Seeker Complete Flow
- [ ] Register new job seeker account
- [ ] Complete basic profile & verify OTP
- [ ] Create detailed profile in User Profile Service
- [ ] Add work experiences (min 2)
- [ ] Add education background
- [ ] Add certifications (if any)
- [ ] Add skills (use bulk add for efficiency)
- [ ] Set career preferences
- [ ] Set position preferences
- [ ] Upload CV
- [ ] Check profile completion status (should be 100%)
- [ ] Browse available jobs
- [ ] Apply to jobs
- [ ] Save favorite jobs
- [ ] View application status

### Company Flow
- [ ] Register new company account
- [ ] Complete company profile
- [ ] Upload logo
- [ ] Create job postings
- [ ] Publish jobs
- [ ] View received applications
- [ ] Update application status
- [ ] View analytics

## 📊 Profile Completion Guide

Untuk mendapatkan **100% profile completion**:

1. ✅ Basic Profile (Profile Service)
2. ✅ Work Experience (min 1)
3. ✅ Education (min 1)
4. ✅ Skills (min 3)
5. ✅ Career Preferences
6. ✅ Position Preferences (min 1)
7. ✅ CV Upload

**Pro Tip:** Gunakan endpoint **"Get Profile with All Relations"** untuk melihat semua data profile dalam satu response!

## 📞 Support

Jika ada pertanyaan atau issue:
1. Check endpoint documentation di Postman
2. Check service logs: `docker compose logs -f [service-name]`
3. Verify environment variables sudah di-set dengan benar
4. Test dengan health check endpoints terlebih dahulu

---

**Updated with User Profile Service! 🎉**

**Total Endpoints: 70+ across 4 services**
**Ready for Production Testing! 🚀**
