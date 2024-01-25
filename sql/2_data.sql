-- +migrate Up
CREATE TABLE lists (
  id BIGINT UNSIGNED PRIMARY KEY UNIQUE,
  parent_id BIGINT UNSIGNED,
  user_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(255) NOT NULL,
  is_shared BOOLEAN DEFAULT FALSE NOT NULL,
  is_favorite BOOLEAN DEFAULT FALSE NOT NULL,
  is_inbox_project BOOLEAN DEFAULT FALSE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (parent_id) REFERENCES lists(id) ON DELETE CASCADE
);

-- +migrate StatementBegin
USE `lists-backend`;
CREATE TRIGGER `lists-backend`.`before_insert_lists`
BEFORE INSERT ON `lists`
FOR EACH ROW
BEGIN
  IF NEW.id IS NULL THEN
    SET NEW.id = uuid_short();
  END IF;
END;
-- +migrate StatementEnd

CREATE INDEX idx_lists_user_id ON lists(user_id);

-- +migrate Down

DROP TABLE IF EXISTS lists;