package store

import (
	"context"
	"database/sql"
	"time"
)

type FollowerStore struct {
	db *sql.DB
}

type Follower struct {
	UserId     int64     `db:"user_id" json:"user_id"`
	FOllowerId int64     `db:"follower_id" json:"follower_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

func (s *FollowerStore) Follow(ctx context.Context, userId int64, followerId int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `INSERT INTO followers (user_id,follower_id) VALUES ($1,$2)`
	_, err := s.db.ExecContext(ctx, query, userId, followerId)
	if err != nil {
		return err
	}
	return nil
}
func (s *FollowerStore) UnFOllow(ctx context.Context, userId int64, unFOllowerId int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`
	_, err := s.db.ExecContext(ctx, query, userId, unFOllowerId)
	if err != nil {
		return err
	}
	return nil
}
