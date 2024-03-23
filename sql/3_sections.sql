-- +migrate Up
CREATE TABLE sections (
  id uuid primary key default uuid_generate_v7(),
  user_id uuid not null,
  list_id uuid not null,
  name varchar(255) not null,
  created_at timestamp default current_timestamp not null,
  updated_at timestamp default current_timestamp not null,
  foreign key (list_id) references lists(id) on delete cascade
);

create index idx_sections_list_id on sections(list_id);

-- +migrate Down
drop table if exists sections;