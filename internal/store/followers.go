package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserId     int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
	CreatedAt  int64 `json:"created_at"`
}

type FollowerStore struct {
	*sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, userID int64) error {
	query := `
	INSERT INTO followers (user_id, follower_id) VALUES($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := s.DB.ExecContext(ctx, query, userID, followerID)
	return err
}
func (s *FollowerStore) Unfollow(ctx context.Context, followerID, userID int64) error {
	query := `
	DELETE FROM followers 
	WHERE user_id=$1
	AND follower_id=$2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := s.DB.ExecContext(ctx, query, userID, followerID)
	return err
}
