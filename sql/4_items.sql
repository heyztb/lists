-- +migrate Up
create table items (
  id uuid primary key default uuid_generate_v7(),
  user_id uuid not null,
  list_id uuid not null,
  parent_id uuid,
  section_id uuid,
  content text not null,
  description text,
  is_completed boolean default false not null,
  labels json,
  priority int not null,
  due timestamp,
  duration int,
  created_at timestamp default current_timestamp not null,
  updated_at timestamp default current_timestamp not null,
  foreign key (list_id) references lists(id) on delete cascade,
  foreign key (section_id) references sections(id) on delete cascade,
  foreign key (user_id) references users(id) on delete cascade,
  foreign key (parent_id) references items(id) on delete cascade
);

create index idx_items_list_id on items(list_id);
create index idx_items_section_id on items(section_id);
create index idx_items_user_id on items(user_id);
create index idx_items_parent_id on items(parent_id);

-- +migrate Down
drop table if exists items;