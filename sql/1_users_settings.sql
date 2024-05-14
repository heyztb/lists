-- +migrate Up
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
  identifier VARCHAR(255) UNIQUE NOT NULL,
  salt VARCHAR(32) NOT NULL,
  verifier VARCHAR(768) NOT NULL,
  mfa_secret VARCHAR(255),
  mfa_recovery_codes TEXT[],
  name VARCHAR(255) DEFAULT 'John Doe',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID UNIQUE NOT NULL,
    session_duration INT DEFAULT 3600 NOT NULL,
    mfa_enabled BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +migrate StatementBegin
CREATE FUNCTION after_insert_users_create_settings_func()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO settings (user_id)
    VALUES (NEW.id);
    RETURN NULL;
END;
$$;

CREATE TRIGGER after_insert_users_create_settings
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION after_insert_users_create_settings_func();
-- +migrate StatementEnd

CREATE INDEX idx_users_identifier ON users(identifier);
CREATE INDEX idx_settings_user_id ON settings(user_id);

-- +migrate Down
DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS users;