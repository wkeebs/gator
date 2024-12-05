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
