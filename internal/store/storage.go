package store

import (
	"context"
	"database/sql"
)

type Posts interface {
	Create(ctx context.Context, post *Post) error
}
type Users interface {
	Create(ctx context.Context, user *User) error
}
type Storage struct {
	Posts Posts
	Users Users
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Posts: &PostStore{db: db},
		Users: &UserStore{db: db},
	}
}
