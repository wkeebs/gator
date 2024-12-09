-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: ResetFeeds :exec
DELETE FROM feeds;

-- name: GetFeeds :many
SELECT name, 
    url, 
    user_id,
    (SELECT users.name FROM users WHERE users.id = user_id) as user_name
FROM feeds;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE feeds.url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at 
NULLS FIRST 
LIMIT 1;