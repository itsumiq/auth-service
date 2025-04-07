BEGIN;

CREATE TABLE refresh_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    refresh_token CHAR(36) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    expire_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC' + INTERVAL '30 days'),
    CONSTRAINT fk_refresh_sessions_users
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE INDEX idx_refresh_sessions_user_id ON refresh_sessions(user_id);
CREATE INDEX idx_refresh_sessions_refresh_token ON refresh_sessions(refresh_token);

CREATE OR REPLACE FUNCTION update_expire_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.expire_at = NOW() AT TIME ZONE 'UTC' + INTERVAL '30 days';
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_refresh_sessions_updated_at_trigger
BEFORE UPDATE ON refresh_sessions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_refresh_sessions_expire_at_trigger
BEFORE UPDATE ON refresh_sessions
FOR EACH ROW
EXECUTE FUNCTION update_expire_at();

COMMIT;
