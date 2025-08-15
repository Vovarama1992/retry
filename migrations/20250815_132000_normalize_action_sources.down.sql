BEGIN;
DROP INDEX IF EXISTS idx_actions_source;
TRUNCATE actions;
INSERT INTO actions SELECT * FROM actions_backup_20250812;
DROP TABLE IF EXISTS actions_backup_20250812;
COMMIT;
