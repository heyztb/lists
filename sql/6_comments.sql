-- +migrate Up
create table comments (
  id uuid primary key default uuid_generate_v7(),
  user_id uuid not null,
  item_id uuid,
  list_id uuid,
  content text not null,
  created_at timestamp default current_timestamp not null,
  updated_at timestamp default current_timestamp not null,
  foreign key (user_id) references users(id) on delete cascade
);

create index idx_comments_user_id on comments(user_id);

-- +migrate Down
drop table if exists comments;
