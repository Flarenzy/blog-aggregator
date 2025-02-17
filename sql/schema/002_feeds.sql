-- +goose Up
CREATE TABLE IF NOT EXISTS feeds (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text not null,
    url text unique not null,
    user_id uuid not null,
    foreign key (user_id) references users(id) on delete cascade
);

-- +goose Down
DROP TABLE feeds;
