-- +goose Up
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL ,
    description TEXT,
    published_at date not null ,
    feed_id UUID NOT NULL ,
    CONSTRAINT fk_feed_id FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;