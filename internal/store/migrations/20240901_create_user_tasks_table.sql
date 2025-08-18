-- +goose Up
create table if not exists user_tasks (
    user_id INTEGER primary key,
    tasks TEXT not null
);

-- +goose Down
drop table if exists user_tasks;
