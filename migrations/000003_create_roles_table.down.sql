BEGIN;

DROP TRIGGER IF EXISTS update_roles_updated_at_trigger ON roles;
DROP TABLE IF EXISTS roles;

COMMIT;
