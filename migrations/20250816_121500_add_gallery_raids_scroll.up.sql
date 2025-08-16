BEGIN;

INSERT INTO action_types (name)
VALUES ('gallery_raids_scroll')
ON CONFLICT (name) DO NOTHING;

COMMIT;

