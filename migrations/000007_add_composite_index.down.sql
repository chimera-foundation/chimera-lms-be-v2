-- DOWN MIGRATION

DROP INDEX IF EXISTS idx_events_type_lookup;
DROP INDEX IF EXISTS idx_events_scoping_composite;
DROP INDEX IF EXISTS idx_events_org_time_lookup;