-- 1. Add the course_id column
ALTER TABLE "assessments" ADD COLUMN "course_id" uuid;

-- 2. Establish the Foreign Key relationship
ALTER TABLE "assessments" 
ADD CONSTRAINT fk_assessments_course 
FOREIGN KEY ("course_id") REFERENCES "courses" ("id") 
ON DELETE CASCADE;

-- 3. Create an index for performance
-- This optimizes the lookup for the subject/course details in your summary endpoint
CREATE INDEX idx_assessments_course_id ON "assessments" ("course_id");