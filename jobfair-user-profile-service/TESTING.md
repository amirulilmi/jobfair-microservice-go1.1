# 🧪 Testing Guide - User Profile Service

Panduan lengkap untuk testing API endpoints.

## 📋 Pilihan Testing

Anda punya 3 cara untuk test API:

### 1️⃣ **Automated Testing (PowerShell)** ⭐ RECOMMENDED untuk Windows

Script otomatis yang test semua endpoints sekaligus.

```powershell
# Jalankan di PowerShell
cd C:\laragon\www\jobfair-microservice\jobfair-user-profile-service
.\test_api.ps1
```

**Kelebihan:**
- ✅ Test 18 endpoints otomatis
- ✅ Colored output (hijau = pass, merah = fail)
- ✅ Summary di akhir
- ✅ Tidak perlu input manual (kecuali login)

---

### 2️⃣ **Automated Testing (Bash)** untuk Linux/Mac/Git Bash

```bash
# Jalankan di terminal
cd /c/laragon/www/jobfair-microservice/jobfair-user-profile-service
chmod +x test_api.sh
./test_api.sh
```

---

### 3️⃣ **Manual Testing** 

Test satu-satu endpoint dengan copy-paste command.

**Buka:** `MANUAL_TESTING.md`

**Cara:**
1. Copy command dari file
2. Paste di terminal/PowerShell
3. Lihat response

**Kelebihan:**
- ✅ Kontrol penuh
- ✅ Bisa fokus ke specific endpoint
- ✅ Debugging lebih mudah

---

## 🚀 Quick Start - PowerShell Testing

### Step 1: Pastikan Service Running

```powershell
# Terminal 1: Jalankan Auth Service
cd C:\laragon\www\jobfair-microservice\jobfair-auth-service
go run cmd/main.go

# Terminal 2: Jalankan User Profile Service  
cd C:\laragon\www\jobfair-microservice\jobfair-user-profile-service
go run cmd/main.go
```

### Step 2: Run Test Script

```powershell
# Terminal 3: Run tests
cd C:\laragon\www\jobfair-microservice\jobfair-user-profile-service
.\test_api.ps1
```

### Step 3: Input Credentials

Script akan minta:
```
Enter email: your-email@example.com
Enter password: ********
```

### Step 4: Lihat Hasil

**Output example:**
```
============================================
JobFair User Profile Service - API Testing
============================================

[1] Testing Health Check...
✓ Test 1: Health Check PASSED

[SETUP] Getting JWT Token from Auth Service...
Enter email: test@example.com
Enter password: ********
✓ JWT Token obtained successfully
Token: eyJhbGciOiJIUzI1Ni...

[2] Testing Create Profile...
✓ Test 2: Create Profile PASSED

[3] Testing Get Profile...
✓ Test 3: Get Profile PASSED

...

============================================
Test Summary
============================================
Total Tests: 18
Passed: 18
Failed: 0

🎉 All tests passed! Service is working correctly.
```

---

## 📊 Test Coverage

Script akan test **18 endpoints**:

### Profile (5 tests)
1. ✅ Create Profile
2. ✅ Get Profile
3. ✅ Update Profile
4. ✅ Get Complete Profile
5. ✅ Get Completion Status

### Work Experience (2 tests)
6. ✅ Create Work Experience
7. ✅ Get All Work Experiences

### Education (2 tests)
8. ✅ Create Education
9. ✅ Get All Educations

### Certifications (2 tests)
10. ✅ Create Certification
11. ✅ Get All Certifications

### Skills (2 tests)
12. ✅ Create Skills (Bulk)
13. ✅ Get All Skills

### Preferences (4 tests)
14. ✅ Create Career Preference
15. ✅ Get Career Preference
16. ✅ Create Position Preferences
17. ✅ Get Position Preferences

### Health Check (1 test)
18. ✅ Health Check

---

## 🎯 Expected Results

Jika **semua berjalan baik:**

```
Total Tests: 18
Passed: 18
Failed: 0

🎉 All tests passed! Service is working correctly.
```

Profile completion seharusnya: **80-93%** (karena belum upload CV)

---

## 🔧 Troubleshooting

### Problem: "Cannot run script"

**PowerShell Execution Policy Issue**

```powershell
# Set execution policy
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Lalu run lagi
.\test_api.ps1
```

### Problem: "Failed to get JWT token"

**Solutions:**

1. **Check Auth Service Running**
   ```powershell
   # Test auth service
   curl http://localhost:8080/health
   ```

2. **Check User Exists**
   ```powershell
   # Register new user dulu
   curl -X POST http://localhost:8080/api/v1/auth/register `
     -H "Content-Type: application/json" `
     -d '{\"email\":\"test@example.com\",\"password\":\"password123\",\"full_name\":\"Test User\"}'
   ```

3. **Check Credentials Correct**
   - Email harus valid
   - Password harus sesuai

### Problem: "Some tests failed"

**Check mana yang fail:**

Script akan show:
```
✗ Test 5: Create Work Experience FAILED
Response: {...}
```

**Common issues:**

1. **Database not migrated**
   ```powershell
   migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/jobfair_user_profile?sslmode=disable" up
   ```

2. **JWT Secret mismatch**
   - Check `.env` di auth-service
   - Check `.env` di user-profile-service
   - Harus SAMA PERSIS!

3. **Service not running**
   ```powershell
   # Check service
   curl http://localhost:8083/health
   ```

---

## 📱 Testing Individual Endpoints

Jika mau test endpoint tertentu saja:

```powershell
# Set token dulu
$TOKEN = "your-jwt-token-here"

# Test specific endpoint
curl -X GET http://localhost:8083/api/v1/profiles `
  -H "Authorization: Bearer $TOKEN"
```

Lihat `MANUAL_TESTING.md` untuk semua curl commands.

---

## 🎨 Postman Collection (Coming Soon)

Untuk testing visual, import Postman collection:

1. Buka Postman
2. Import → File → `postman_collection.json`
3. Set environment variable `token`
4. Run collection

---

## 💡 Tips

1. **Run test setelah perubahan code** untuk ensure tidak break
2. **Save test results** untuk comparison
3. **Test di environment berbeda** (dev, staging, prod)
4. **Automate testing** dalam CI/CD pipeline

---

## 📈 Next Steps

Setelah manual testing berhasil:

1. ✅ Write unit tests
2. ✅ Write integration tests
3. ✅ Setup CI/CD dengan automated testing
4. ✅ Add performance testing
5. ✅ Add load testing

---

## 🆘 Need Help?

Jika masih ada error:

1. Check service logs
2. Check database connection
3. Check JWT token validity
4. Lihat `MANUAL_TESTING.md` untuk test manual
5. Lihat `TROUBLESHOOTING.md` untuk common issues

---

**Happy Testing! 🎉**
