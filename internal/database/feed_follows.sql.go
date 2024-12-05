// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1, 
        $2, 
        $3,  
        (SELECT users.id FROM users WHERE users.id = $4),
        (SELECT feeds.id FROM feeds WHERE feeds.id = $5)
    )
    RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT 
    inserted.id, inserted.created_at, inserted.updated_at, inserted.user_id, inserted.feed_id,
    users.name AS user_name,
    feeds.name AS feed_name
FROM 
    inserted
    INNER JOIN users ON inserted.user_id = users.id
    INNER JOIN feeds ON inserted.feed_id = feeds.id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	ID_2      uuid.UUID
	ID_3      uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	UserName  string
	FeedName  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.ID_2,
		arg.ID_3,
	)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.UserName,
		&i.FeedName,
	)
	return i, err
}