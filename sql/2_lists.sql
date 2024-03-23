-- +migrate Up
CREATE TABLE lists (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
  parent_id UUID,
  user_id UUID NOT NULL,
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
CREATE FUNCTION after_insert_users_create_inbox_func()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO lists (parent_id, user_id, name, is_inbox_project)
    VALUES (NULL, NEW.id, 'Inbox', TRUE);
    RETURN NULL;
END;
$$;

CREATE TRIGGER after_insert_users_create_inbox
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION after_insert_users_create_inbox_func();
-- +migrate StatementEnd

CREATE INDEX idx_lists_parent_id on lists(parent_id);
CREATE INDEX idx_lists_user_id ON lists(user_id);

-- +migrate Down

DROP TABLE IF EXISTS lists;