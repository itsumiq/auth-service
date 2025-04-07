BEGIN;

DROP TRIGGER IF EXISTS update_permissions_updated_at_trigger ON permissions;
DROP TABLE IF EXISTS permissions;

COMMIT;
