# 🚀 Quick Start Guide - JobFair User Profile Service

## ⚡ Menjalankan Service di Local (Tanpa Docker)

### Step 1: Setup Database

```bash
# 1. Pastikan PostgreSQL sudah running (Laragon)
# 2. Buat database
psql -U postgres -c "CREATE DATABASE jobfair_user_profile;"
```

### Step 2: Setup Environment

```bash
# 1. Buat file .env (sudah ada, tinggal edit jika perlu)
# 2. Pastikan JWT_SECRET SAMA dengan auth-service!

# Cek isi .env:
PORT=8083
DATABASE_URL=postgres://postgres:postgres@localhost:5432/jobfair_user_profile?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-12345
```

### Step 3: Install Dependencies

```bash
go mod download
```

### Step 4: Run Migrations

```bash
# Install migrate tool (one time only)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/jobfair_user_profile?sslmode=disable" up
```

### Step 5: Run Service

```bash
go run cmd/main.go
```

**✅ Service berjalan di: http://localhost:8083**

---

## 🧪 Testing Service

### 1. Health Check (No Auth)

```bash
curl http://localhost:8083/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "user-profile-service"
}
```

### 2. Test dengan JWT Token

**Cara mendapatkan JWT Token:**
1. Jalankan `jobfair-auth-service` di port 8080
2. Register/Login untuk mendapat token
3. Copy token yang didapat

**Test Create Profile:**

```bash
curl -X POST http://localhost:8083/api/v1/profiles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "phone_number": "081234567890"
  }'
```

Expected response:
```json
{
  "success": true,
  "message": "Profile created successfully",
  "data": {
    "id": "uuid-here",
    "user_id": "uuid-here",
    "full_name": "John Doe",
    "phone_number": "081234567890",
    "completion_status": 0,
    "created_at": "2025-10-01T...",
    "updated_at": "2025-10-01T..."
  }
}
```

**Test Get Profile:**

```bash
curl -X GET http://localhost:8083/api/v1/profiles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

---

## 🔗 Integration dengan Auth Service

### Flow Lengkap:

```
1. User → Auth Service (POST /api/v1/auth/register)
   Response: { "token": "eyJhbGc..." }

2. User → User Profile Service (POST /api/v1/profiles)
   Header: Authorization: Bearer eyJhbGc...
   Response: Profile created!

3. User → User Profile Service (GET /api/v1/profiles/full)
   Header: Authorization: Bearer eyJhbGc...
   Response: Complete profile data
```

### PENTING: Sinkronisasi JWT Secret

Pastikan JWT_SECRET di kedua service **SAMA PERSIS**:

**Auth Service (.env):**
```env
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-12345
```

**User Profile Service (.env):**
```env
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-12345
```

---

## 📝 Common Tasks

### Add Work Experience

```bash
curl -X POST http://localhost:8083/api/v1/work-experiences \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "PT Example",
    "job_position": "Software Engineer",
    "start_date": "2020-01-01T00:00:00Z",
    "is_current_job": true,
    "job_description": "Developing microservices..."
  }'
```

### Add Education

```bash
curl -X POST http://localhost:8083/api/v1/educations \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "university": "Universitas Indonesia",
    "major": "Computer Science",
    "degree": "Bachelor",
    "start_date": "2016-08-01T00:00:00Z",
    "end_date": "2020-07-31T00:00:00Z",
    "gpa": 3.75
  }'
```

### Upload CV

```bash
curl -X POST http://localhost:8083/api/v1/cv \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@C:/path/to/your/cv.pdf"
```

### Get Complete Profile

```bash
curl -X GET http://localhost:8083/api/v1/profiles/full \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## ⚠️ Troubleshooting

### Error: "database connection failed"

**Solution:**
```bash
# Pastikan PostgreSQL running
# Cek dengan:
psql -U postgres -h localhost

# Jika tidak bisa connect, start PostgreSQL di Laragon
```

### Error: "migration failed"

**Solution:**
```bash
# Reset migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/jobfair_user_profile?sslmode=disable" down
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/jobfair_user_profile?sslmode=disable" up
```

### Error: "unauthorized" / "invalid token"

**Solution:**
1. Pastikan JWT_SECRET sama di auth-service dan user-profile-service
2. Pastikan token belum expired
3. Cek format header: `Authorization: Bearer <token>` (ada spasi setelah Bearer)

### Error: "profile not found"

**Solution:**
Profile dibuat otomatis saat first access. Pastikan:
1. Token valid
2. Auth service sudah running
3. User sudah register/login

---

## 🎯 Next Steps

Setelah service berjalan, Anda bisa:

1. ✅ Integrate dengan frontend mobile app
2. ✅ Test semua endpoints
3. ✅ Add unit tests
4. ✅ Setup monitoring & logging
5. ✅ Deploy ke production

---

## 📞 Need Help?

- Check **README.md** untuk dokumentasi lengkap
- Check **Makefile** untuk available commands
- Run `make help` untuk list commands

**Happy Coding! 🚀**
