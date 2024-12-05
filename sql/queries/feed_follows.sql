-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1, 
        $2, 
        $3,  
        (SELECT users.id FROM users WHERE users.id = $4),
        (SELECT feeds.id FROM feeds WHERE feeds.id = $5)
    )
    RETURNING *
)
SELECT 
    inserted.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM 
    inserted
    INNER JOIN users ON inserted.user_id = users.id
    INNER JOIN feeds ON inserted.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT *, 
    (SELECT users.name FROM users WHERE users.id = feed_follows.user_id) as user_name,
    (SELECT feeds.name FROM feeds WHERE feeds.id = feed_follows.feed_id) as feed_name
FROM feed_follows
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1
AND (SELECT feeds.id FROM feeds WHERE feeds.url = $2) = feed_follows.feed_id;
