-- Drop the table first (removes indexes automatically)
DROP TABLE IF EXISTS events;

-- Drop the custom types
DROP TYPE IF EXISTS event_type;
DROP TYPE IF EXISTS event_scope;