BEGIN;

DELETE FROM action_types
WHERE name = 'gallery_raids_scroll';

COMMIT;