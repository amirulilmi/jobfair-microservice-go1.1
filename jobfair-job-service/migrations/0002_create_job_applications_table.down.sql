DROP TRIGGER IF EXISTS update_job_applications_updated_at ON job_applications;
DROP INDEX IF EXISTS idx_job_applications_deleted_at;
DROP INDEX IF EXISTS idx_job_applications_created_at;
DROP INDEX IF EXISTS idx_job_applications_status;
DROP INDEX IF EXISTS idx_job_applications_user_id;
DROP INDEX IF EXISTS idx_job_applications_job_id;
DROP TABLE IF EXISTS job_applications;
