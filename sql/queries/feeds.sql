-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: ResetFeeds :exec
DELETE FROM feeds;

-- name: GetFeeds :many
SELECT name, 
    url, 
    (SELECT users.name FROM users WHERE users.id = user_id) as user_name
FROM feeds;
