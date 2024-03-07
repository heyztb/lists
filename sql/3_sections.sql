-- +migrate Up
CREATE TABLE sections (
  id bigint unsigned primary key unique,
  user_id bigint unsigned not null,
  list_id bigint unsigned not null,
  name varchar(255) not null,
  created_at timestamp default current_timestamp not null,
  updated_at timestamp default current_timestamp not null,
  foreign key (list_id) references lists(id) on delete cascade
);
-- +migrate StatementBegin
USE `lists-backend`;
CREATE TRIGGER `lists-backend`.`before_insert_sections`
BEFORE INSERT ON `sections`
FOR EACH ROW
BEGIN
  IF NEW.id IS NULL THEN
    SET NEW.id = uuid_short();
  END IF;
END;
-- +migrate StatementEnd

create index idx_sections_list_id on sections(list_id);

-- +migrate Down
drop table if exists sections;