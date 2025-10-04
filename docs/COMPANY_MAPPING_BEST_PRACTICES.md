# Company Mapping System - Best Practices Guide

## Overview

Company mappings synchronize company data between `company-service` and `job-service` to enable job posting functionality. This document outlines the architecture, data flow, and best practices for maintaining data consistency.

## Architecture

### Services Involved
- **company-service**: Manages company profiles (PostgreSQL: `jobfair_companies`)
- **job-service**: Manages job postings (PostgreSQL: `jobfair_jobs`)
- **RabbitMQ**: Message broker for event-driven communication

### Data Flow

```
1. Company Registration
   company-service → RabbitMQ → job-service
   
2. Automatic Mapping Creation
   job-service creates mapping in company_mappings table
   
3. Job Creation
   job-service validates company exists via mapping
```

## Tables

### company-service: `companies`
```sql
- id (PK)
- user_id (FK to users)
- name
- description
- ...
```

### job-service: `company_mappings`
```sql
- id (PK)
- user_id (FK to users) UNIQUE
- company_id (references companies.id from company-service)
- company_name (denormalized for quick access)
- created_at
- updated_at
```

## Event-Driven Sync

### Events Published by company-service

1. **company.registered**
```json
{
  "event_type": "company.registered",
  "data": {
    "user_id": 123,
    "company_id": 456,
    "company_name": "Tech Corp",
    "email": "contact@techcorp.com"
  }
}
```

2. **company.updated**
```json
{
  "event_type": "company.updated",
  "data": {
    "user_id": 123,
    "updated_fields": {
      "company_id": 456,
      "company_name": "Tech Corp Updated"
    }
  }
}
```

3. **company.deleted**
```json
{
  "event_type": "company.deleted",
  "data": {
    "user_id": 123,
    "company_id": 456
  }
}
```

### Consumer in job-service

Located at: `internal/consumers/company_event_consumer.go`

- Listens to company events
- Auto-creates/updates/deletes company_mappings
- Uses `UpsertCompanyMapping()` for idempotency

## Data Consistency Tools

### 1. Health Check API

Check for data inconsistencies:

```bash
GET /api/v1/admin/health/data-consistency
Authorization: Bearer <ADMIN_TOKEN>

Response:
{
  "success": true,
  "status": "healthy|degraded|unhealthy",
  "checks": {
    "jobs_without_mapping": 0,
    "orphaned_applications": 0
  },
  "issues": [],
  "warnings": []
}
```

### 2. Manual Sync API

Manually sync a company mapping:

```bash
POST /api/v1/admin/sync-company-mapping
Authorization: Bearer <ADMIN_TOKEN>
Content-Type: application/json

{
  "user_id": 123,
  "company_id": 456,
  "company_name": "Tech Corp"
}
```

### 3. List All Mappings

View all company mappings:

```bash
GET /api/v1/admin/company-mappings
Authorization: Bearer <ADMIN_TOKEN>
```

### 4. Safe Sync Script

Automated bulk sync using API:

```bash
export ADMIN_TOKEN="your-admin-jwt-token"
./scripts/sync-company-mappings-safe.sh
```

## Best Practices

### ✅ DO

1. **Use Event-Driven Sync**: Let RabbitMQ events handle synchronization automatically
2. **Use Admin APIs**: For manual interventions, use the provided admin endpoints
3. **Monitor Health**: Regularly check `/admin/health/data-consistency`
4. **Use Idempotent Operations**: `UpsertCompanyMapping` handles duplicates gracefully
5. **Log Everything**: Admin actions are logged for audit trail
6. **Validate Permissions**: Admin endpoints check `user_type == "admin"`

### ❌ DON'T

1. **Direct Database Access**: Never manually INSERT/UPDATE database records
2. **Skip Validation**: Always validate data before syncing
3. **Ignore Errors**: Monitor consumer logs for event processing failures
4. **Production SQL**: Never run SQL commands in production without review
5. **Bypass Business Logic**: Always use services/repositories, not raw SQL

## Troubleshooting

### Issue: "Company profile not found"

**Symptom**: Job creation fails with "Company profile not found"

**Diagnosis**:
```bash
# Check if mapping exists
GET /api/v1/admin/company-mappings

# Check health
GET /api/v1/admin/health/data-consistency
```

**Solution**:
1. Check if company exists in company-service
2. Run safe sync script
3. If recurring, check RabbitMQ consumer logs

### Issue: Orphaned Jobs

**Symptom**: Jobs exist but company was deleted

**Diagnosis**:
```bash
GET /api/v1/admin/health/data-consistency
# Check "warnings" section
```

**Solution**:
- These are expected when companies are deleted
- Jobs remain for historical record
- Consider soft-delete strategy

### Issue: Missing Mappings After Migration

**Symptom**: Existing companies can't create jobs

**Solution**:
```bash
# Run migration repair
docker exec -it postgres-job psql -U jobfair_user -d jobfair_jobs \
  -f /app/migrations/0005_repair_company_mappings.up.sql

# Or use safe sync script
export ADMIN_TOKEN="..."
./scripts/sync-company-mappings-safe.sh
```

## Development Workflow

### Adding New Company
1. User registers as company → company-service
2. Create company profile → company-service publishes event
3. job-service consumer creates mapping automatically
4. Company can now create jobs ✅

### Testing Event Flow
```bash
# 1. Create company
POST /api/v1/companies
Authorization: Bearer <COMPANY_TOKEN>

# 2. Check mapping was created (may take 1-2 seconds)
GET /api/v1/admin/company-mappings
Authorization: Bearer <ADMIN_TOKEN>

# 3. Create job (should work)
POST /api/v1/jobs
Authorization: Bearer <COMPANY_TOKEN>
```

## Monitoring

### Key Metrics to Track
- Event processing latency
- Failed event processing count
- Mapping inconsistency count
- Admin API usage

### Logs to Monitor
```bash
# Consumer logs
docker logs -f jobfair-job-service | grep "JOB-SERVICE"

# Admin actions
docker logs -f jobfair-job-service | grep "ADMIN"
```

## Security

### Admin Endpoints
- Require JWT with `user_type: "admin"`
- All actions are logged
- Rate-limited (if configured)

### Event Security
- RabbitMQ uses authentication
- Events validated before processing
- Failed events logged but don't crash service

## Migration Guide

### From Manual Sync to Event-Driven

1. **Enable Consumer**: Ensure RabbitMQ is configured
2. **Run Repair Migration**: Sync existing data
3. **Test Event Flow**: Create test company
4. **Monitor**: Watch logs for 24-48 hours
5. **Cleanup**: Remove manual sync scripts (keep for backup)

## FAQ

**Q: What if RabbitMQ is down?**
A: Services continue working. Events queue up and process when RabbitMQ recovers.

**Q: Can I use direct SQL in emergency?**
A: Only in development. Production requires admin API or documented repair script.

**Q: How often should I check health?**
A: Run health check in CI/CD pipeline or daily cron job.

**Q: What's the recovery time for sync issues?**
A: Safe sync script typically completes in < 1 minute for 1000 companies.

## References

- Event Schema: `shared/events/company_events.go`
- Consumer: `internal/consumers/company_event_consumer.go`
- Admin Handler: `internal/handlers/admin_handler.go`
- Repository: `internal/repository/company_repository.go`

---

**Last Updated**: October 2025
**Version**: 1.0
**Maintained by**: Platform Team
