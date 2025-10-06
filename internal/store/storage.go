package store

import (
	"context"
	"database/sql"
	"time"
)

type Posts interface {
	Create(ctx context.Context, post *Post) error
	GetAll(ctx context.Context) ([]*Post, error)
	GetByID(ctx context.Context, id int64) (*Post, error)
	Delete(ctx context.Context, postID int64) error
	Update(ctx context.Context, postID int64, post *Post) error
	GetUserFeed(ctx context.Context, userId int64, fq PaginatedFeedQuery) (*[]PostWithMetaData, error)
}
type Users interface {
	Create(ctx context.Context, user *User) error
	GetUsers(ctx context.Context) (*[]User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserById(ctx context.Context, id int64) (*User, error)
}
type Comments interface {
	GetCommentsWithPost(ctx context.Context, postID int64) (*[]Comment, error)
	CreateComment(ctx context.Context, comments *Comment) error
}
type Followers interface {
	Follow(ctx context.Context, userId int64, followerId int64) error
	UnFOllow(ctx context.Context, userId int64, unFOllowerId int64) error
}
type Storage struct {
	Posts     Posts
	Users     Users
	Comments  Comments
	Followers Followers
}

var (
	QueryTimeoutDuration = time.Second * 5
)

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Posts:     &PostStore{db: db},
		Users:     &UserStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}
