-- +migrate Up
CREATE TABLE users (
    id BIGINT UNSIGNED PRIMARY KEY UNIQUE,
    identifier VARCHAR(255) UNIQUE NOT NULL,
    salt VARCHAR(32) NOT NULL,
    verifier VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE settings (
  id BIGINT UNSIGNED PRIMARY KEY UNIQUE,
  user_id BIGINT UNSIGNED UNIQUE NOT NULL,
  session_duration INT DEFAULT 3600 NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +migrate StatementBegin
USE `lists-backend`;
CREATE TRIGGER `lists-backend`.`before_insert_users`
BEFORE INSERT ON `users`
FOR EACH ROW
BEGIN
  IF NEW.id IS NULL THEN
    SET NEW.id = uuid_short();
  END IF;
END;
-- +migrate StatementEnd

-- +migrate StatementBegin
USE `lists-backend`;
CREATE TRIGGER `lists-backend`.`after_insert_users_create_settings`
AFTER INSERT ON `users`
FOR EACH ROW
BEGIN
  INSERT INTO settings(id, user_id) VALUES (uuid_short(), NEW.id);
END;
-- +migrate StatementEnd

CREATE INDEX idx_users_identifier ON users(identifier);
CREATE INDEX idx_settings_user_id ON settings(user_id);

-- +migrate Down
DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS users;