BEGIN;

UPDATE action_types
SET name = 'external_link_garage_raid'
WHERE name = 'external_link_raid';

COMMIT;