-- 1. Remove the index
DROP INDEX IF EXISTS idx_assessments_course_id;

-- 2. Remove the column (Postgres automatically drops the constraint)
ALTER TABLE "assessments" DROP COLUMN IF EXISTS "course_id";