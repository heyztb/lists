-- +migrate Up
create table labels (
  id bigint unsigned primary key unique,
  user_id bigint unsigned not null,
  name text not null,
  is_favorite boolean default false,
  created_at timestamp default current_timestamp not null,
  updated_at timestamp default current_timestamp not null,
  foreign key (user_id) references users(id) on delete cascade
);

-- +migrate StatementBegin
USE `lists-backend`;
CREATE TRIGGER `lists-backend`.`before_insert_labels`
BEFORE INSERT ON `labels`
FOR EACH ROW
BEGIN
  IF NEW.id IS NULL THEN
    SET NEW.id = uuid_short();
  END IF;
END;
-- +migrate StatementEnd

create index idx_labels_user_id on labels(user_id);

-- +migrate Down
drop table if exists labels;