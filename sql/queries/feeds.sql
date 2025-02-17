-- name: CreateFeed :one
INSERT INTO "feeds" (id, created_at, updated_at, name, url, user_id)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
       ) RETURNING *;

-- name: GetUsersFeeds :many
SELECT *
FROM "feeds"
WHERE user_id = $1;

-- name: GetFeedByUrl :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: GetAllFeeds :many
SELECT feeds.name, feeds.url, users.name
FROM feeds
INNER JOIN users
ON feeds.user_id = users.id;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $2, last_fetched_at = $3
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;
