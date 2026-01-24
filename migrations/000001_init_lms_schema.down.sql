-- 1. Drop Foreign Key Constraints
-- Note: Dropping the tables with CASCADE would also work, 
-- but explicitly dropping constraints is cleaner for migration logs.

-- 2. Drop Junction Tables (Tables with multiple dependencies)
DROP TABLE IF EXISTS "section_members" CASCADE;
DROP TABLE IF EXISTS "user_roles" CASCADE;
DROP TABLE IF EXISTS "program_courses" CASCADE;

-- 3. Drop Data/Transaction Tables
DROP TABLE IF EXISTS "progress_trackers" CASCADE;
DROP TABLE IF EXISTS "submissions" CASCADE;
DROP TABLE IF EXISTS "enrollments" CASCADE;
DROP TABLE IF EXISTS "contents" CASCADE;
DROP TABLE IF EXISTS "assessments" CASCADE;
DROP TABLE IF EXISTS "lessons" CASCADE;
DROP TABLE IF EXISTS "modules" CASCADE;
DROP TABLE IF EXISTS "courses" CASCADE;

-- 4. Drop Structural/Organizational Tables
DROP TABLE IF EXISTS "sections" CASCADE;
DROP TABLE IF EXISTS "cohorts" CASCADE;
DROP TABLE IF EXISTS "programs" CASCADE;
DROP TABLE IF EXISTS "subjects" CASCADE;
DROP TABLE IF EXISTS "education_levels" CASCADE;
DROP TABLE IF EXISTS "academic_periods" CASCADE;
DROP TABLE IF EXISTS "roles" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "organizations" CASCADE;

-- 5. Drop Custom Types
DROP TYPE IF EXISTS role_type;
DROP TYPE IF EXISTS assessment_type;
DROP TYPE IF EXISTS enrollment_status;
DROP TYPE IF EXISTS content_type;
DROP TYPE IF EXISTS course_status;
DROP TYPE IF EXISTS organization_type;

-- 6. Optional: Drop Extensions
-- Only include this if the extension wasn't already in your DB before this migration
-- DROP EXTENSION IF EXISTS "pgcrypto";