-- 1. Create Custom Types for Enums
CREATE TYPE event_type AS ENUM (
    'holiday', 'deadline', 'session', 'vanilla', 'meeting', 'schedule'
);

CREATE TYPE event_scope AS ENUM (
    'global', 'cohort', 'section', 'personal'
);

-- 2. Create the Events Table
CREATE TABLE IF NOT EXISTS events (
    -- Base Fields (from shared.Base)
    id UUID PRIMARY KEY,
    created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at timestamp WITH TIME ZONE,

    -- Core Fields
    organization_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    event_type event_type NOT NULL,
    color VARCHAR(7) NOT NULL DEFAULT '#3B82F6', -- Hex format

    -- Timing Fields
    start_at timestamp WITH TIME ZONE,
    end_at timestamp WITH TIME ZONE,
    is_all_day BOOLEAN NOT NULL DEFAULT FALSE,
    recurrence_rule TEXT,

    -- Scoping Fields
    scope event_scope NOT NULL DEFAULT 'global',
    cohort_id UUID,
    section_id UUID,
    user_id UUID,

    -- Polymorphic Source Links
    source_id UUID,
    source_type VARCHAR(50),

    -- Media
    image_url TEXT
);

-- 3. Optimization: Composite Index for the 'Find' Method
-- This index covers Organization + Scope + Time, which is your most common query path.
CREATE INDEX idx_events_lookup ON events (organization_id, scope, start_at) 
WHERE deleted_at IS NULL;

-- 4. Specific Indexes for Foreign Keys (Standard Practice)
CREATE INDEX idx_events_section_id ON events (section_id) WHERE section_id IS NOT NULL;
CREATE INDEX idx_events_cohort_id ON events (cohort_id) WHERE cohort_id IS NOT NULL;
CREATE INDEX idx_events_user_id ON events (user_id) WHERE user_id IS NOT NULL;