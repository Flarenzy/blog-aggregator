-- +goose Up
CREATE TABLE IF NOT EXISTS feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL ,
    feed_id UUID NOT NULL ,
    FOREIGN KEY (user_id) references users(id) on delete cascade,
    FOREIGN KEY (feed_id) references feeds(id) ON DELETE CASCADE,
    CONSTRAINT unique_fks UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;