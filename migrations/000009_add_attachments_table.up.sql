-- Attachments table for assessment and submission file uploads
CREATE TABLE "attachments" (
    "id"              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "organization_id" uuid NOT NULL REFERENCES organizations(id),
    "uploaded_by"     uuid NOT NULL REFERENCES users(id),
    "assessment_id"   uuid REFERENCES assessments(id),
    "submission_id"   uuid REFERENCES submissions(id),
    "file_name"       varchar NOT NULL,
    "file_url"        varchar NOT NULL,
    "file_size"       bigint NOT NULL,
    "mime_type"       varchar NOT NULL,
    "created_at"      timestamptz DEFAULT now(),
    "deleted_at"      timestamptz
);

-- Partial indexes for fast lookups (exclude soft-deleted)
CREATE INDEX idx_attachments_assessment ON attachments(assessment_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_attachments_submission ON attachments(submission_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_attachments_uploaded_by ON attachments(uploaded_by) WHERE deleted_at IS NULL;
