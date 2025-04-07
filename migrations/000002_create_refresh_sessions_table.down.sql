BEGIN;

DROP TRIGGER IF EXISTS update_refresh_sessions_expire_at_trigger ON refresh_sessions;
DROP TRIGGER IF EXISTS update_refresh_sessions_updated_at_trigger ON refresh_sessions;

DROP FUNCTION IF EXISTS update_expire_at();

DROP INDEX IF EXISTS idx_refresh_sessions_user_id;
DROP INDEX IF EXISTS idx_refresh_sessions_refresh_token;

DROP TABLE IF EXISTS refresh_sessions;

COMMIT;
