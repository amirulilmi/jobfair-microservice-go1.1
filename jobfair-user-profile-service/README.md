# 🎯 JobFair User Profile Service

Microservice untuk mengelola profil pengguna dalam ekosistem JobFair. Service ini menangani data komprehensif pencari kerja termasuk informasi pribadi, pengalaman kerja, pendidikan, sertifikasi, skills, preferensi karir, CV, dan badge system.

## 🚀 Tech Stack

- **Language:** Go 1.23
- **Framework:** Gin
- **ORM:** GORM
- **Database:** PostgreSQL
- **Authentication:** JWT
- **Port:** 8083

## 📋 Prerequisites

Pastikan Anda sudah menginstall:

- [Go 1.23+](https://golang.org/dl/)
- [PostgreSQL 14+](https://www.postgresql.org/download/)
- [Make](https://www.gnu.org/software/make/) (optional, untuk commands)
- [golang-migrate](https://github.com/golang-migrate/migrate) (untuk database migrations)

## 🛠️ Setup & Installation

### 1. Clone Repository

```bash
cd C:\laragon\www\jobfair-microservice\jobfair-user-profile-service
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

Atau menggunakan Make:
```bash
make deps
```

### 3. Setup Environment Variables

Copy file `.env.example` ke `.env`:

```bash
cp .env.example .env
```

Edit `.env` dan sesuaikan dengan konfigurasi Anda:

```env
# Server Configuration
PORT=8083
SERVICE_NAME=user-profile-service

# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/jobfair_user_profile?sslmode=disable

# JWT Configuration (HARUS SAMA dengan auth-service!)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-12345

# File Upload Configuration
MAX_FILE_SIZE=5242880  # 5MB
ALLOWED_FILE_TYPES=.pdf,.doc,.docx
UPLOAD_DIR=./uploads/cv
```

**⚠️ PENTING:** Pastikan `JWT_SECRET` **SAMA PERSIS** dengan yang digunakan di `jobfair-auth-service`!

### 4. Create Database

Buat database PostgreSQL:

```bash
# Menggunakan psql
psql -U postgres -c "CREATE DATABASE jobfair_user_profile;"
```

Atau menggunakan Make:
```bash
make db-create
```

### 5. Run Database Migrations

Install golang-migrate tool:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Atau:
```bash
make migrate-install
```

Jalankan migrations:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/jobfair_user_profile?sslmode=disable" up
```

Atau menggunakan Make:
```bash
make migrate-up
```

### 6. Run Service

```bash
go run cmd/main.go
```

Atau menggunakan Make:
```bash
make run
```

Service akan berjalan di `http://localhost:8083`

## 🎯 Available Make Commands

```bash
make help              # Show all available commands
make run               # Run the application
make build             # Build the application
make test              # Run tests
make clean             # Clean build artifacts
make deps              # Install dependencies
make dev               # Run with auto-reload (requires air)

# Database commands
make db-create         # Create database
make db-drop           # Drop database
make db-reset          # Reset database (drop, create, migrate)

# Migration commands
make migrate-install   # Install golang-migrate tool
make migrate-up        # Run migrations
make migrate-down      # Rollback migrations
make migrate-create NAME=migration_name  # Create new migration

# Docker commands
make docker-build      # Build Docker image
make docker-run        # Run Docker container

# Setup
make setup             # Complete setup (env, deps, migrate tool)
```

## 📡 API Endpoints

### Health Check (No Auth Required)

```
GET /health
```

### Profile Endpoints

```
POST   /api/v1/profiles              # Create profile
GET    /api/v1/profiles              # Get own profile
GET    /api/v1/profiles/full         # Get profile with all relations
PUT    /api/v1/profiles              # Update profile
GET    /api/v1/profiles/completion   # Get completion status
```

### Work Experience Endpoints

```
POST   /api/v1/work-experiences      # Add work experience
GET    /api/v1/work-experiences      # List work experiences
GET    /api/v1/work-experiences/:id  # Get single work experience
PUT    /api/v1/work-experiences/:id  # Update work experience
DELETE /api/v1/work-experiences/:id  # Delete work experience
```

### Education Endpoints

```
POST   /api/v1/educations            # Add education
GET    /api/v1/educations            # List educations
GET    /api/v1/educations/:id        # Get single education
PUT    /api/v1/educations/:id        # Update education
DELETE /api/v1/educations/:id        # Delete education
```

### Certification Endpoints

```
POST   /api/v1/certifications        # Add certification
GET    /api/v1/certifications        # List certifications
GET    /api/v1/certifications/:id    # Get single certification
PUT    /api/v1/certifications/:id    # Update certification
DELETE /api/v1/certifications/:id    # Delete certification
```

### Skill Endpoints

```
POST   /api/v1/skills                # Add skill
POST   /api/v1/skills/bulk           # Add multiple skills
GET    /api/v1/skills                # List skills
GET    /api/v1/skills/:id            # Get single skill
PUT    /api/v1/skills/:id            # Update skill
DELETE /api/v1/skills/:id            # Delete skill
```

### Career Preference Endpoints

```
POST   /api/v1/career-preference     # Create/update career preference
GET    /api/v1/career-preference     # Get career preference
```

### Position Preference Endpoints

```
POST   /api/v1/position-preferences  # Add position preferences (bulk)
GET    /api/v1/position-preferences  # List position preferences
DELETE /api/v1/position-preferences/:id # Delete position preference
```

### CV Endpoints

```
POST   /api/v1/cv                    # Upload CV (multipart/form-data)
GET    /api/v1/cv                    # Get CV info
DELETE /api/v1/cv                    # Delete CV
```

### Badge Endpoints

```
GET    /api/v1/badges                # List user's badges
```

## 🔐 Authentication

Semua endpoints (kecuali `/health`) memerlukan JWT token dalam header:

```
Authorization: Bearer <your-jwt-token>
```

Token didapat dari `jobfair-auth-service` setelah login/register.

## 📝 Request Examples

### 1. Create Profile

```bash
curl -X POST http://localhost:8083/api/v1/profiles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "phone_number": "081234567890"
  }'
```

### 2. Add Work Experience

```bash
curl -X POST http://localhost:8083/api/v1/work-experiences \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "PT Bumi Mineral Sejahtera",
    "job_position": "Senior Mining Engineer",
    "start_date": "2020-01-15T00:00:00Z",
    "end_date": "2024-07-31T00:00:00Z",
    "is_current_job": false,
    "job_description": "Led a team of 15 engineers..."
  }'
```

### 3. Add Education

```bash
curl -X POST http://localhost:8083/api/v1/educations \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "university": "Institut Teknologi Bandung",
    "major": "Teknik Pertambangan",
    "degree": "Bachelor",
    "start_date": "2014-08-01T00:00:00Z",
    "end_date": "2018-07-31T00:00:00Z",
    "is_current": false,
    "gpa": 3.98
  }'
```

### 4. Add Skills (Bulk)

```bash
curl -X POST http://localhost:8083/api/v1/skills/bulk \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "technical_skills": [
      {
        "skill_name": "Mine Planning & Design",
        "skill_type": "technical",
        "proficiency_level": "expert",
        "years_of_experience": 5
      }
    ],
    "soft_skills": [
      {
        "skill_name": "Leadership",
        "skill_type": "soft",
        "proficiency_level": "advanced"
      }
    ]
  }'
```

### 5. Upload CV

```bash
curl -X POST http://localhost:8083/api/v1/cv \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/your/cv.pdf"
```

### 6. Get Complete Profile

```bash
curl -X GET http://localhost:8083/api/v1/profiles/full \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🔄 Integration dengan Auth Service

Service ini terintegrasi dengan `jobfair-auth-service` melalui JWT authentication:

1. User login/register di auth service → mendapat JWT token
2. JWT token berisi `user_id` dalam claims
3. User Profile Service membaca `user_id` dari JWT token
4. Data profile di-associate dengan `user_id` tersebut

**Flow:**
```
User → Auth Service (Login) → JWT Token
     → User Profile Service (dengan JWT Token)
     → Create/Read/Update Profile
```

## 🗄️ Database Schema

Service ini menggunakan 9 tabel utama:

1. **profiles** - Data profil utama
2. **work_experiences** - Pengalaman kerja
3. **educations** - Riwayat pendidikan
4. **certifications** - Sertifikasi
5. **skills** - Technical & soft skills
6. **career_preferences** - Preferensi karir (salary, work type, dll)
7. **position_preferences** - Posisi yang diminati
8. **cv_documents** - Dokumen CV
9. **badges** - Badge gamification
10. **profile_badges** - Many-to-many junction table

## 📊 Profile Completion Tracking

Service menghitung profile completion (0-100%) berdasarkan 15 kriteria:

- ✅ Full name (required)
- ✅ Phone number (required)
- ✅ Bio
- ✅ Date of birth
- ✅ City
- ✅ At least 1 work experience
- ✅ At least 1 education
- ✅ At least 1 certification
- ✅ At least 1 technical skill
- ✅ At least 1 soft skill
- ✅ Career preference set
- ✅ Position preferences (at least 1)
- ✅ CV uploaded
- ✅ Profile picture uploaded
- ✅ LinkedIn URL added

**Formula:** `(completed_fields / 15) * 100`

## 🐳 Docker Deployment

Build Docker image:

```bash
docker build -t jobfair-user-profile-service:latest .
```

Run container:

```bash
docker run -p 8083:8083 --env-file .env jobfair-user-profile-service:latest
```

## 🧪 Testing

Run tests:

```bash
go test -v ./...
```

Atau:
```bash
make test
```

## 📁 Project Structure

```
jobfair-user-profile-service/
├── cmd/
│   └── main.go                    # Entry point
├── internal/
│   ├── config/                    # Configuration
│   ├── handlers/                  # HTTP Handlers
│   ├── services/                  # Business Logic
│   ├── repository/                # Data Access Layer
│   ├── models/                    # Domain Models
│   ├── middleware/                # JWT Middleware
│   └── utils/                     # Helper functions
├── pkg/
│   └── database/                  # Database connection
├── migrations/                    # Database migrations
├── deployments/
│   ├── docker/
│   └── k8s/
├── .env                          # Environment variables
├── .env.example                  # Example environment
├── Dockerfile                    # Docker configuration
├── Makefile                      # Build automation
└── go.mod                        # Go dependencies
```

## 🔧 Troubleshooting

### Database Connection Failed

```bash
# Pastikan PostgreSQL running
# Windows (dengan Laragon):
# Buka Laragon → Start PostgreSQL

# Cek koneksi:
psql -U postgres -h localhost -p 5432
```

### Migration Failed

```bash
# Reset database
make db-reset

# Atau manual:
make db-drop
make db-create
make migrate-up
```

### JWT Token Invalid

Pastikan:
1. JWT_SECRET di `.env` sama dengan auth service
2. Token belum expired
3. Format header: `Authorization: Bearer <token>`

### File Upload Failed

Pastikan:
1. Folder `uploads/cv` exists dan writable
2. File size < 5MB
3. File type: .pdf, .doc, atau .docx

## 📞 Support

Jika ada masalah, silakan buat issue di repository atau hubungi tim development.

## 📄 License

[Your License Here]

---

**Made with ❤️ for JobFair Microservices**
