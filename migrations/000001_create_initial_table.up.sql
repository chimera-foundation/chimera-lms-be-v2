CREATE TABLE "organizations" (
  "id" uuid PRIMARY KEY,
  "name" varchar,
  "slug" varchar UNIQUE,
  "type" enum,
  "settings" jsonb,
  "created_at" timestamp DEFAULT (now())
  "created_by" FOREIGN KEY
  "updated_at" timestamp
  "updated_by" FOREIGN KEY
  "deleted_at" timestamp
);

CREATE TABLE "academic_periods" (
  "id" uuid PRIMARY KEY,
  "organization_id" uuid,
  "name" varchar,
  "slug" varchar,
  "start_date" date,
  "end_date" date,
  "is_active" bool DEFAULT true
);

CREATE TABLE "education_levels" (
  "id" uuid PRIMARY KEY,
  "organization_id" uuid,
  "name" varchar,
  "code" varchar
);

CREATE TABLE "subjects" (
  "id" uuid PRIMARY KEY,
  "organization_id" uuid,
  "education_level_id" uuid,
  "name" varchar,
  "code" varchar
);

CREATE TABLE "programs" (
  "id" uuid PRIMARY KEY,
  "organization_id" uuid,
  "name" varchar,
  "description" text
);

CREATE TABLE "courses" (
  "id" uuid PRIMARY KEY,
  "organization_id" uuid,
  "instructor_id" uuid,
  "subject_id" uuid,
  "education_level_id" uuid,
  "grade_level" int,
  "credits" int DEFAULT 0,
  "title" varchar,
  "description" text,
  "status" enum,
  "price" decimal
);

CREATE TABLE "program_courses" (
  "program_id" uuid,
  "course_id" uuid,
  "order_index" int,
  "primary" key(program_id,course_id)
);

CREATE TABLE "modules" (
  "id" uuid PRIMARY KEY,
  "course_id" uuid,
  "title" varchar,
  "order_index" int
);

CREATE TABLE "lessons" (
  "id" uuid PRIMARY KEY,
  "module_id" uuid,
  "title" varchar,
  "order_index" int
);

CREATE TABLE "contents" (
  "id" uuid PRIMARY KEY,
  "lesson_id" uuid,
  "assessment_id" uuid,
  "content_type" enum,
  "content_data" jsonb
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "organization_id" uuid,
  "email" varchar UNIQUE,
  "password_hash" varchar,
  "is_superuser" bool DEFAULT false
);

CREATE TABLE "roles" (
  "id" uuid PRIMARY KEY,
  "name" varchar,
  "permissions" jsonb
);

CREATE TABLE "user_roles" (
  "user_id" uuid,
  "role_id" uuid,
  "primary" key(user_id,role_id)
);

CREATE TABLE "cohorts" (
  "id" uuid PRIMARY KEY,
  "organization_id" uuid,
  "academic_period_id" uuid,
  "education_level_id" uuid,
  "name" varchar
);

CREATE TABLE "sections" (
  "id" uuid PRIMARY KEY,
  "cohort_id" uuid,
  "name" varchar,
  "room_number" varchar,
  "capacity" int
);

CREATE TABLE "section_members" (
  "section_id" uuid,
  "user_id" uuid,
  "role_type" enum,
  "primary" key(section_id,user_id)
);

CREATE TABLE "enrollments" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid,
  "course_id" uuid,
  "section_id" uuid,
  "academic_period_id" uuid,
  "status" enum,
  "enrolled_at" timestamp DEFAULT (now())
);

CREATE TABLE "assessments" (
  "id" uuid PRIMARY KEY,
  "organization_id" uuid,
  "title" varchar,
  "assessment_type" enum,
  "assessment_subtype" enum,
  "due_date" timestamp
);

CREATE TABLE "submissions" (
  "id" uuid PRIMARY KEY,
  "assessment_id" uuid,
  "user_id" uuid,
  "enrollment_id" uuid,
  "final_score" float,
  "submitted_at" timestamp
);

CREATE TABLE "progress_trackers" (
  "id" uuid PRIMARY KEY,
  "enrollment_id" uuid,
  "content_id" uuid,
  "is_completed" bool DEFAULT false,
  "updated_at" timestamp
);

ALTER TABLE "academic_periods" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "education_levels" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "subjects" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "subjects" ADD FOREIGN KEY ("education_level_id") REFERENCES "education_levels" ("id");

ALTER TABLE "programs" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "courses" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "courses" ADD FOREIGN KEY ("instructor_id") REFERENCES "users" ("id");

ALTER TABLE "courses" ADD FOREIGN KEY ("subject_id") REFERENCES "subjects" ("id");

ALTER TABLE "courses" ADD FOREIGN KEY ("education_level_id") REFERENCES "education_levels" ("id");

ALTER TABLE "program_courses" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id");

ALTER TABLE "program_courses" ADD FOREIGN KEY ("course_id") REFERENCES "courses" ("id");

ALTER TABLE "modules" ADD FOREIGN KEY ("course_id") REFERENCES "courses" ("id");

ALTER TABLE "lessons" ADD FOREIGN KEY ("module_id") REFERENCES "modules" ("id");

ALTER TABLE "contents" ADD FOREIGN KEY ("lesson_id") REFERENCES "lessons" ("id");

ALTER TABLE "contents" ADD FOREIGN KEY ("assessment_id") REFERENCES "assessments" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "user_roles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "cohorts" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "cohorts" ADD FOREIGN KEY ("academic_period_id") REFERENCES "academic_periods" ("id");

ALTER TABLE "cohorts" ADD FOREIGN KEY ("education_level_id") REFERENCES "education_levels" ("id");

ALTER TABLE "sections" ADD FOREIGN KEY ("cohort_id") REFERENCES "cohorts" ("id");

ALTER TABLE "section_members" ADD FOREIGN KEY ("section_id") REFERENCES "sections" ("id");

ALTER TABLE "section_members" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "enrollments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "enrollments" ADD FOREIGN KEY ("course_id") REFERENCES "courses" ("id");

ALTER TABLE "enrollments" ADD FOREIGN KEY ("section_id") REFERENCES "sections" ("id");

ALTER TABLE "enrollments" ADD FOREIGN KEY ("academic_period_id") REFERENCES "academic_periods" ("id");

ALTER TABLE "assessments" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "submissions" ADD FOREIGN KEY ("assessment_id") REFERENCES "assessments" ("id");

ALTER TABLE "submissions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "submissions" ADD FOREIGN KEY ("enrollment_id") REFERENCES "enrollments" ("id");

ALTER TABLE "progress_trackers" ADD FOREIGN KEY ("enrollment_id") REFERENCES "enrollments" ("id");

ALTER TABLE "progress_trackers" ADD FOREIGN KEY ("content_id") REFERENCES "contents" ("id");
