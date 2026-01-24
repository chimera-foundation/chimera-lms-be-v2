CREATE TABLE "cohort_members" (
  "cohort_id" uuid REFERENCES "cohorts"("id") ON DELETE CASCADE,
  "user_id" uuid REFERENCES "users"("id") ON DELETE CASCADE,
  "created_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "updated_at" timestamp WITH TIME ZONE DEFAULT (now()),
  "deleted_at" timestamp WITH TIME ZONE,
  "created_by" uuid,
  "updated_by" uuid,
  "deleted_by" uuid,
  PRIMARY KEY (cohort_id, user_id)
);

-- Index for reverse lookups (finding all cohorts for a user)
CREATE INDEX idx_cohort_members_user_id ON cohort_members (user_id);