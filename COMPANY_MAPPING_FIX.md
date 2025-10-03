# Fix: Company Mapping Error di Job Service

## 🔍 Masalah

Error ketika POST `/api/v1/jobs`:
```
ERROR: relation "companies" does not exist (SQLSTATE 42P01)
SELECT "id" FROM "companies" WHERE user_id = 2
```

**Root Cause**: Job-service mencoba query ke table `companies` di database sendiri, padahal table tersebut ada di database company-service (arsitektur microservices yang benar - setiap service punya database terpisah).

## ✅ Solusi yang Diterapkan

### 1. **Event-Driven Architecture**
- Company-service akan **re-publish** event `company.registered` dengan `company_id` setelah company dibuat
- Job-service akan **consume** event tersebut dan menyimpan mapping `user_id -> company_id` di table lokal

### 2. **Database Denormalization**
- Job-service punya table `company_mappings` untuk menyimpan mapping sederhana
- Table ini di-populate melalui event consumer

### 3. **Perubahan yang Dibuat**

#### A. Shared Events Library
- ✅ Update `CompanyRegisteredData` untuk include `company_id`

#### B. Company Service
- ✅ Consumer re-publish event dengan `company_id` setelah company dibuat
- ✅ Publisher ditambahkan ke consumer

#### C. Job Service
- ✅ Consumer baru untuk consume company events
- ✅ Repository method untuk upsert dan delete mapping
- ✅ Struct `CompanyMapping` dengan timestamps yang benar
- ✅ TableName() method untuk explicit table mapping
- ✅ Main.go diupdate untuk start consumer

## 🚀 Cara Perbaikan

### Step 1: Run Migration
Migration sudah ada di `jobfair-job-service/migrations/0004_create_company_mappings_table.up.sql`

```bash
# Masuk ke container job-service
docker exec -it jobfair-job-service bash

# Run migration
migrate -path=/app/migrations -database="$DATABASE_URL" up
```

Atau via docker-compose:
```bash
cd jobfair-job-service
docker-compose down
docker-compose up -d --build
```

### Step 2: Populate Existing Companies

Untuk companies yang sudah ada sebelum fix ini, jalankan script populate:

```bash
# Pastikan chmod +x dulu
chmod +x scripts/populate-company-mappings.sh

# Run script
./scripts/populate-company-mappings.sh
```

Script ini akan:
1. Ambil semua companies dari `jobfair_company` database
2. Insert mapping ke `company_mappings` table di `jobfair_jobs` database

### Step 3: Restart Services

```bash
# Restart semua services untuk apply changes
docker-compose down
docker-compose up -d --build
```

### Step 4: Test

1. **Login sebagai company**
```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "company@example.com",
    "password": "password123"
  }'
```

2. **Create Job** (gunakan token dari login)
```bash
curl -X POST http://localhost:8000/api/v1/jobs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Senior Software Engineer",
    "description": "Looking for talented engineers",
    "employment_type": "full_time",
    "work_type": "remote",
    "experience_level": "senior",
    "location": "Jakarta, Indonesia",
    "requirements": ["5+ years experience", "Go/Python"],
    "receive_method": "internal"
  }'
```

## 🔄 Flow Setelah Fix

### New Company Registration:
1. User register sebagai company di **auth-service**
2. Auth-service publish `company.registered` event (without company_id)
3. **Company-service** consume event → create company → **re-publish event WITH company_id**
4. **Job-service** consume re-published event → save mapping di `company_mappings`

### Create Job:
1. User (company) call POST `/api/v1/jobs`
2. Job-service lookup `company_mappings` table: `user_id → company_id`
3. Job-service create job dengan `company_id` yang ditemukan

## 📋 Verification Checklist

Setelah fix, pastikan:

- [ ] Migration `0004_create_company_mappings_table` sudah jalan
- [ ] Table `company_mappings` ada di database `jobfair_jobs`
- [ ] Existing companies sudah ter-populate via script
- [ ] Consumer job-service berjalan (check logs: "✅ Company event consumer started")
- [ ] Consumer company-service berjalan
- [ ] Bisa create job tanpa error "relation companies does not exist"

## 🐛 Troubleshooting

### Error: "Company profile not found"
**Penyebab**: Mapping belum ada di `company_mappings` table

**Solusi**:
1. Check apakah consumer berjalan:
```bash
docker logs jobfair-job-service | grep "Company event consumer"
```

2. Manual insert untuk testing:
```sql
-- Connect ke jobfair_jobs database
INSERT INTO company_mappings (user_id, company_id, company_name, created_at, updated_at)
VALUES (2, 1, 'Test Company', NOW(), NOW())
ON CONFLICT (user_id) DO UPDATE SET company_id = EXCLUDED.company_id;
```

3. Run populate script untuk existing companies

### Consumer Tidak Start
**Check**:
```bash
docker logs jobfair-job-service
```

Cari log:
- ✅ "Company event consumer started" → OK
- ⚠️ "Warning: Failed to initialize company event consumer" → RabbitMQ issue

**Fix RabbitMQ**:
```bash
docker-compose restart rabbitmq
docker-compose restart jobfair-job-service
```

### Event Tidak Ter-consume
**Check RabbitMQ**:
```bash
# Access RabbitMQ Management UI
open http://localhost:15672
# Login: guest/guest

# Check queues:
# - job-service-company-queue should exist
# - Should have consumers connected
```

## 📊 Database Schema

### company_mappings Table
```sql
CREATE TABLE company_mappings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    company_id INTEGER NOT NULL,
    company_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_company_mappings_user_id ON company_mappings(user_id);
CREATE INDEX idx_company_mappings_company_id ON company_mappings(company_id);
```

## 📁 File Changes Summary

1. `jobfair-shared-libs/go/events/models.go` - Added `CompanyID` field
2. `jobfair-company-service/internal/consumers/company_event_consumer.go` - Re-publish event with company_id
3. `jobfair-job-service/internal/consumers/company_event_consumer.go` - NEW consumer
4. `jobfair-job-service/internal/repository/company_repository.go` - Updated struct & added methods
5. `jobfair-job-service/cmd/main.go` - Initialize & start consumer
6. `scripts/populate-company-mappings.sh` - NEW populate script

## 🎯 Next Steps

Untuk companies baru, semuanya otomatis via event system.

Untuk existing companies, jalankan populate script sekali saja.

Happy coding! 🚀
