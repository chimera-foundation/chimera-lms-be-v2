-- Remove the "System" Organization using its fixed UUID
DELETE FROM organizations 
WHERE id = '00000000-0000-4000-a000-000000000000';
