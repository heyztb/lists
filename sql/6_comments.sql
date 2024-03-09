-- +migrate Up
create table comments (
  id bigint unsigned primary key unique,
  user_id bigint unsigned not null,
  item_id bigint unsigned,
  list_id bigint unsigned,
  content text not null,
  created_at timestamp default current_timestamp not null,
  updated_at timestamp default current_timestamp not null,
  foreign key (user_id) references users(id) on delete cascade
);

-- +migrate StatementBegin
USE `lists-backend`;
CREATE TRIGGER `lists-backend`.`before_insert_comments`
BEFORE INSERT ON `comments`
FOR EACH ROW
BEGIN
  IF NEW.id IS NULL THEN
    SET NEW.id = uuid_short();
  END IF;
END;
-- +migrate StatementEnd

create index idx_comments_user_id on comments(user_id);

-- +migrate Down
drop table if exists labels;
