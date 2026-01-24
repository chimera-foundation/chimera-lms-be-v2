-- Add first_name and last_name to users table
ALTER TABLE "users" 
ADD COLUMN "first_name" varchar NOT NULL DEFAULT '',
ADD COLUMN "last_name" varchar NOT NULL DEFAULT '';