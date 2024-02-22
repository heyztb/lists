-- +migrate Up
create table items (
  id bigint unsigned primary key unique,
  list_id bigint unsigned not null,
  section_id bigint unsigned,
  creator_id bigint unsigned not null,
  content text not null,
  description text,
  is_completed boolean default false not null,
  labels json,
  parent_id bigint unsigned,
  position int not null,
  priority int not null,
  due timestamp,
  duration int,
  created_at timestamp default current_timestamp not null,
  updated_at timestamp default current_timestamp not null,
  foreign key (list_id) references lists(id) on delete cascade,
  foreign key (section_id) references sections(id) on delete cascade,
  foreign key (creator_id) references users(id) on delete cascade,
  foreign key (parent_id) references items(id) on delete cascade
);

-- +migrate StatementBegin
USE `lists-backend`;
CREATE TRIGGER `lists-backend`.`before_insert_items`
BEFORE INSERT ON `items`
FOR EACH ROW
BEGIN
  IF NEW.id IS NULL THEN
    SET NEW.id = uuid_short();
  END IF;
END;
-- +migrate StatementEnd

create index idx_items_list_id on items(list_id);
create index idx_items_section_id on items(section_id);
create index idx_items_creator_id on items(creator_id);
create index idx_items_parent_id on items(parent_id);

-- +migrate Down
drop table if exists items;