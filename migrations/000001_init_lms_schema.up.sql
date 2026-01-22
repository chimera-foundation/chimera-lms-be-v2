-- 1. Setup Extensions and Types
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

DO $$ BEGIN
    CREATE TYPE organization_type AS ENUM ('school', 'university', 'corporate', 'bootcamp');
    CREATE TYPE course_status AS ENUM ('draft', 'published', 'archived');
    CREATE TYPE content_type AS ENUM ('video', 'document', 'quiz', 'assignment');
    CREATE TYPE enrollment_status AS ENUM ('active', 'completed', 'dropped', 'pending');
    CREATE TYPE assessment_type AS ENUM ('quiz', 'exam', 'project');
    CREATE TYPE role_type AS ENUM ('student', 'teacher', 'admin');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- 2. Tables with Audit Columns
CREATE TABLE "organizations" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar NOT NULL,
  "slug" varchar UNIQUE NOT NULL,
  "type" organization_type,
  "settings" jsonb,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "updated_by" uuid
  "deleted_by" uuid
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "organization_id" uuid,
  "email" varchar UNIQUE NOT NULL,
  "password_hash" varchar NOT NULL,
  "is_superuser" bool DEFAULT false,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "updated_by" uuid
  "deleted_by" uuid
);

CREATE TABLE "academic_periods" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "organization_id" uuid,
  "name" varchar NOT NULL,
  "slug" varchar,
  "start_date" date,
  "end_date" date,
  "is_active" bool DEFAULT true,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "education_levels" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "organization_id" uuid,
  "name" varchar NOT NULL,
  "code" varchar,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "subjects" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "organization_id" uuid,
  "education_level_id" uuid,
  "name" varchar NOT NULL,
  "code" varchar,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "programs" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "organization_id" uuid,
  "name" varchar NOT NULL,
  "description" text,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "courses" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "organization_id" uuid,
  "instructor_id" uuid,
  "subject_id" uuid,
  "education_level_id" uuid,
  "grade_level" int,
  "credits" int DEFAULT 0,
  "title" varchar NOT NULL,
  "description" text,
  "status" course_status,
  "price" decimal(12,2),
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "modules" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "course_id" uuid,
  "title" varchar NOT NULL,
  "order_index" int,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "lessons" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "module_id" uuid,
  "title" varchar NOT NULL,
  "order_index" int,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "assessments" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "organization_id" uuid,
  "title" varchar NOT NULL,
  "assessment_type" assessment_type,
  "due_date" timestamp WITH TIME ZONE,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "contents" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "lesson_id" uuid,
  "assessment_id" uuid,
  "content_type" content_type,
  "content_data" jsonb,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "roles" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar NOT NULL,
  "permissions" jsonb,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "cohorts" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "organization_id" uuid,
  "academic_period_id" uuid,
  "education_level_id" uuid,
  "name" varchar NOT NULL,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "sections" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "cohort_id" uuid,
  "name" varchar NOT NULL,
  "room_number" varchar,
  "capacity" int,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "enrollments" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" uuid,
  "course_id" uuid,
  "section_id" uuid,
  "academic_period_id" uuid,
  "status" enrollment_status,
  "enrolled_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "submissions" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "assessment_id" uuid,
  "user_id" uuid,
  "enrollment_id" uuid,
  "final_score" float,
  "submitted_at" timestamp WITH TIME ZONE,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "deleted_by" uuid
  "updated_by" uuid
);

CREATE TABLE "progress_trackers" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "enrollment_id" uuid,
  "content_id" uuid,
  "is_completed" bool DEFAULT false,
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now())
);

-- 3. Junction Tables (Typically no audit columns needed unless history is required)
CREATE TABLE "program_courses" (
  "program_id" uuid REFERENCES "programs"("id"),
  "course_id" uuid REFERENCES "courses"("id"),
  "order_index" int,
  PRIMARY KEY (program_id, course_id)
);

CREATE TABLE "user_roles" (
  "user_id" uuid REFERENCES "users"("id"),
  "role_id" uuid REFERENCES "roles"("id"),
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE "section_members" (
  "section_id" uuid REFERENCES "sections"("id"),
  "user_id" uuid REFERENCES "users"("id"),
  "role_type" role_type,
  PRIMARY KEY (section_id, user_id)
);

-- 4. Foreign Key Constraints (Relationships)
ALTER TABLE "organizations" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "organizations" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");
ALTER TABLE "users" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "users" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "academic_periods" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");
ALTER TABLE "academic_periods" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "academic_periods" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "education_levels" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");
ALTER TABLE "education_levels" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "education_levels" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "subjects" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");
ALTER TABLE "subjects" ADD FOREIGN KEY ("education_level_id") REFERENCES "education_levels" ("id");
ALTER TABLE "subjects" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "subjects" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "programs" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");
ALTER TABLE "programs" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "programs" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "courses" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");
ALTER TABLE "courses" ADD FOREIGN KEY ("instructor_id") REFERENCES "users" ("id");
ALTER TABLE "courses" ADD FOREIGN KEY ("subject_id") REFERENCES "subjects" ("id");
ALTER TABLE "courses" ADD FOREIGN KEY ("education_level_id") REFERENCES "education_levels" ("id");
ALTER TABLE "courses" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "courses" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "modules" ADD FOREIGN KEY ("course_id") REFERENCES "courses" ("id");
ALTER TABLE "modules" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "modules" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "lessons" ADD FOREIGN KEY ("module_id") REFERENCES "modules" ("id");
ALTER TABLE "lessons" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "lessons" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "contents" ADD FOREIGN KEY ("lesson_id") REFERENCES "lessons" ("id");
ALTER TABLE "contents" ADD FOREIGN KEY ("assessment_id") REFERENCES "assessments" ("id");
ALTER TABLE "contents" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "contents" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "cohorts" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");
ALTER TABLE "cohorts" ADD FOREIGN KEY ("academic_period_id") REFERENCES "academic_periods" ("id");
ALTER TABLE "cohorts" ADD FOREIGN KEY ("education_level_id") REFERENCES "education_levels" ("id");
ALTER TABLE "cohorts" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "cohorts" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "sections" ADD FOREIGN KEY ("cohort_id") REFERENCES "cohorts" ("id");
ALTER TABLE "sections" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "sections" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "enrollments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "enrollments" ADD FOREIGN KEY ("course_id") REFERENCES "courses" ("id");
ALTER TABLE "enrollments" ADD FOREIGN KEY ("section_id") REFERENCES "sections" ("id");
ALTER TABLE "enrollments" ADD FOREIGN KEY ("academic_period_id") REFERENCES "academic_periods" ("id");
ALTER TABLE "enrollments" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "enrollments" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "assessments" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");
ALTER TABLE "assessments" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "assessments" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "submissions" ADD FOREIGN KEY ("assessment_id") REFERENCES "assessments" ("id");
ALTER TABLE "submissions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "submissions" ADD FOREIGN KEY ("enrollment_id") REFERENCES "enrollments" ("id");

ALTER TABLE "progress_trackers" ADD FOREIGN KEY ("enrollment_id") REFERENCES "enrollments" ("id");
ALTER TABLE "progress_trackers" ADD FOREIGN KEY ("content_id") REFERENCES "contents" ("id");