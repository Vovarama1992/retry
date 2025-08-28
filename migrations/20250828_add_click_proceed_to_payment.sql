BEGIN;

INSERT INTO action_types (name)
VALUES ('click_proceed_to_payment')
ON CONFLICT (name) DO NOTHING;

COMMIT;