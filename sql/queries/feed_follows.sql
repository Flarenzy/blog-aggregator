-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
INSERT INTO "feed_follows" (
                          id,
                          created_at,
                          updated_at,
                          user_id,
                          feed_id
) VALUES
      (
          $1,
          $2,
          $3,
          $4,
          $5
      )
       RETURNING *
         )
SELECT inserted_feed_follow.*
, feeds.name as feed_name
, users.name as user_name
FROM inserted_feed_follow
INNER JOIN feeds
ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users
ON users.id = inserted_feed_follow.user_id;

-- name: GetFeedFollowsForUser :many
WITH temp_users AS (
SELECT *
FROM users
WHERE users.name = $1
)
SELECT temp_users.name as username,
        feeds.name AS feed_name
FROM temp_users
INNER JOIN feed_follows
    ON feed_follows.user_id = temp_users.id
INNER JOIN feeds
    ON feeds.id = feed_follows.feed_id
;

-- name: DeleteUserAndFeed :exec
DELETE
FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;