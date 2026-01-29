-- UP MIGRATION

-- 1. Optimized for filtering by Organization and Time (The primary filter)
CREATE INDEX idx_events_org_time_lookup 
ON events (organization_id, start_at ASC) 
WHERE deleted_at IS NULL;

-- 2. Optimized for the "Who" logic (Scoping)
-- This helps the planner jump quickly to Section, Cohort, or User specific events
CREATE INDEX idx_events_scoping_composite
ON events (organization_id, scope, user_id, section_id, cohort_id)
WHERE deleted_at IS NULL;

-- 3. Optimized for Type filtering (The "What")
-- We use a GIN index if you find you're filtering by many types at once, 
-- but for standard ENUM-like EventTypes, a B-Tree on the column is usually faster.
CREATE INDEX idx_events_type_lookup 
ON events (organization_id, event_type) 
WHERE deleted_at IS NULL;