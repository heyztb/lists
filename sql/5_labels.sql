-- +migrate Up
create table labels (
  id uuid primary key default uuid_generate_v7(),
  user_id uuid not null,
  name text not null,
  color text not null,
  is_favorite boolean default false,
  created_at timestamp default current_timestamp not null,
  updated_at timestamp default current_timestamp not null,
  foreign key (user_id) references users(id) on delete cascade
);

create index idx_labels_user_id on labels(user_id);

-- +migrate Down
drop table if exists labels;