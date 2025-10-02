# 🎯 Job Service

Job Service adalah microservice yang mengelola job postings, job applications, dan saved jobs dalam platform jobfair.

## 📋 Features

### For Companies
- ✅ **Post Jobs** - Create job postings with detailed information
- ✅ **Manage Jobs** - Update, publish, close, and delete job postings
- ✅ **Track Applications** - View and manage applications per job
- ✅ **Application Status** - Update application status (shortlisted, interview, hired, rejected)
- ✅ **Statistics** - View application statistics by status

### For Job Seekers
- ✅ **Browse Jobs** - Search and filter available jobs
- ✅ **Job Details** - View detailed job information with company info
- ✅ **Apply to Jobs** - Submit applications with CV and cover letter
- ✅ **Bulk Apply** - Apply to multiple jobs at once (One-Click Apply)
- ✅ **Save Jobs** - Bookmark jobs for later
- ✅ **Track Applications** - Monitor application status
- ✅ **Withdraw Applications** - Cancel submitted applications

## 🏗️ Architecture

```
┌─────────────────────────────────────────────┐
│           Job Service API                   │
├─────────────────────────────────────────────┤
│                                             │
│  ┌──────────────┐    ┌──────────────┐     │
│  │   Jobs       │    │ Applications │     │
│  │   Handler    │    │   Handler    │     │
│  └──────┬───────┘    └──────┬───────┘     │
│         │                    │              │
│  ┌──────▼───────┐    ┌──────▼───────┐     │
│  │   Job        │    │ Application  │     │
│  │   Service    │    │   Service    │     │
│  └──────┬───────┘    └──────┬───────┘     │
│         │                    │              │
│  ┌──────▼───────────────────▼───────┐     │
│  │      Repository Layer            │     │
│  │  - JobRepository                 │     │
│  │  - ApplicationRepository         │     │
│  │  - SavedJobRepository            │     │
│  └──────────────┬───────────────────┘     │
│                 │                          │
└─────────────────┼──────────────────────────┘
                  │
         ┌────────▼────────┐
         │   PostgreSQL    │
         │  (jobfair_jobs) │
         └─────────────────┘
```

## 🗄️ Database Schema

### Tables
1. **jobs** - Job postings
2. **job_applications** - Job applications from users
3. **saved_jobs** - Bookmarked jobs

### Key Fields

**Jobs:**
- Basic info: title, description, slug
- Employment: type, work type, experience level
- Location: location, city, country, remote
- Salary: min, max, currency, period
- Requirements, responsibilities, skills, benefits
- Status: draft, published, closed, archived
- Tracking: views, applications count

**Applications:**
- Job and user references
- CV URL, cover letter
- Status: applied, reviewing, shortlisted, interview, hired, rejected
- Tracking timestamps: viewed, reviewed, interview, responded

## 🚀 API Endpoints

### Public Endpoints
```
GET    /api/v1/jobs              # List jobs with filters
GET    /api/v1/jobs/:id          # Get job detail
```

### Company Endpoints
```
POST   /api/v1/jobs              # Create job
PUT    /api/v1/jobs/:id          # Update job
DELETE /api/v1/jobs/:id          # Delete job
POST   /api/v1/jobs/:id/publish  # Publish job
POST   /api/v1/jobs/:id/close    # Close job
GET    /api/v1/jobs/my           # Get my posted jobs

GET    /api/v1/jobs/:job_id/applications  # Get applications for job
PUT    /api/v1/applications/:id/status    # Update application status
GET    /api/v1/applications/stats         # Get application statistics
```

### Job Seeker Endpoints
```
POST   /api/v1/jobs/:id/apply    # Apply to job
POST   /api/v1/jobs/bulk-apply   # Apply to multiple jobs
POST   /api/v1/jobs/:id/save     # Save job
DELETE /api/v1/jobs/:id/save     # Unsave job
GET    /api/v1/jobs/saved        # Get saved jobs

GET    /api/v1/applications/my   # Get my applications
GET    /api/v1/applications/:id  # Get application detail
DELETE /api/v1/applications/:id  # Withdraw application
```

## 📦 Installation

### Prerequisites
- Go 1.23+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Local Development

1. **Clone and navigate:**
```bash
cd jobfair-job-service
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Setup environment:**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Run migrations:**
```bash
# Using migrate CLI
migrate -path ./migrations -database "postgres://user:pass@localhost:5435/jobfair_jobs?sslmode=disable" up
```

5. **Run service:**
```bash
go run cmd/main.go
```

### Docker

```bash
# From root project
docker-compose up -d job-service
```

## 🧪 Testing Examples

### Create Job (Company)
```bash
curl -X POST http://localhost:8082/api/v1/jobs \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mining Operations Manager",
    "description": "Manage mining operations...",
    "employment_type": "fulltime",
    "work_type": "onsite",
    "experience_level": "mid",
    "location": "Jakarta, Indonesia",
    "salary_min": 5000000,
    "salary_max": 8000000,
    "requirements": ["5+ years experience", "Mining certification"],
    "responsibilities": ["Oversee operations", "Manage team"],
    "skills": ["Leadership", "Mining Operations"],
    "receive_method": "email",
    "contact_email": "hr@company.com"
  }'
```

### List Jobs (Public)
```bash
curl "http://localhost:8082/api/v1/jobs?page=1&limit=10&employment_type=fulltime&location=Jakarta"
```

### Apply to Job (Job Seeker)
```bash
curl -X POST http://localhost:8082/api/v1/jobs/1/apply \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "cv_url": "/uploads/cv/user123.pdf",
    "cover_letter": "I am interested in this position..."
  }'
```

### Bulk Apply (Job Seeker)
```bash
curl -X POST http://localhost:8082/api/v1/jobs/bulk-apply \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "job_ids": [1, 2, 3, 4],
    "cv_url": "/uploads/cv/user123.pdf"
  }'
```

### Update Application Status (Company)
```bash
curl -X PUT http://localhost:8082/api/v1/applications/1/status \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "shortlisted",
    "status_note": "Impressive background"
  }'
```

## 🔍 Query Filters

### Job List Filters
- `search` - Search in title and description
- `employment_type` - fulltime, parttime, contract, freelance, intern
- `work_type` - onsite, remote, hybrid
- `experience_level` - entry, junior, mid, senior
- `location` - Filter by location or city
- `salary_min` - Minimum salary
- `salary_max` - Maximum salary
- `company_id` - Filter by company
- `tags` - Filter by tags
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)
- `order_by` - created_at, views, applications
- `order` - asc, desc

Example:
```
GET /api/v1/jobs?employment_type=fulltime&work_type=remote&experience_level=mid&salary_min=5000000&page=1&limit=10
```

## 📊 Application Status Flow

```
Applied → Reviewing → Shortlisted → Interview → Hired/Rejected
```

**Status Meanings:**
- `applied` - Initial submission
- `reviewing` - Company is reviewing
- `shortlisted` - Selected for next stage
- `interview` - Interview scheduled/completed
- `hired` - Offer accepted
- `rejected` - Application declined

## 🔐 Authentication

Service menggunakan JWT authentication. Token didapat dari Auth Service.

**Required Headers:**
```
Authorization: Bearer <jwt_token>
```

**Token Claims:**
- `user_id` - User ID
- `user_type` - "company" or "job_seeker"

## 🎨 Design Implementation

Service ini mengimplementasikan design dari PDF yang diberikan:

1. **Job List** - Browse dengan filter (employment type, work type, experience level, salary)
2. **Job Detail** - View lengkap dengan company info
3. **One-Click Apply** - Bulk apply ke multiple jobs
4. **Application Tracking** - Track status (Applied, Shortlisted, Interview, Hired, Rejected)
5. **Saved Jobs** - Bookmark functionality
6. **Company Dashboard** - Manage jobs dan applications

## 🚨 Error Handling

Service menggunakan standard HTTP status codes:

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Missing/invalid token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## 🔄 Integration

### With Auth Service
- JWT token validation
- User authentication

### With Company Service
- Company information retrieval
- Company profile data

## 📝 Environment Variables

```env
DATABASE_URL=postgres://user:pass@localhost:5435/jobfair_jobs?sslmode=disable
PORT=8082
JWT_SECRET=your-jwt-secret
RABBITMQ_URL=amqp://user:pass@localhost:5672/
AUTH_SERVICE_URL=http://localhost:8080
COMPANY_SERVICE_URL=http://localhost:8081
```

## 🐛 Troubleshooting

### Database connection failed
```bash
# Check database is running
docker-compose ps postgres-job

# Check connection
psql -h localhost -p 5435 -U jobfair_user -d jobfair_jobs
```

### JWT validation failed
- Ensure JWT_SECRET matches auth service
- Check token expiration
- Verify token format

### Application already exists
- User can only apply once per job
- Use GET to check if already applied

## 📚 Resources

- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)

---

**Version:** 1.0.0  
**Last Updated:** 2025-01-01
