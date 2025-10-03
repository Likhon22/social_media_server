package store

import (
	"context"
	"database/sql"
)

type Posts interface {
	Create(context.Context) error
}
type Users interface {
	Create(context.Context) error
}
type Storage struct {
	Posts Posts
	Users Users
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Posts: &PostsStore{db: db},
		Users: &UsersStore{db: db},
	}
}
