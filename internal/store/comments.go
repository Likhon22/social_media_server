package store

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID        int       `db:"id" json:"id"`
	PostID    int       `db:"post_id" json:"post_id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	User      User      `db:"user" json:"user"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetCommentsWithPost(ctx context.Context, postID int64) (*[]Comment, error) {
	query := `
        SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at, u.id, u.username
        FROM comments c
        JOIN users u ON u.id = c.user_id
        WHERE c.post_id = $1
        ORDER BY c.created_at DESC
    `

	var comments []Comment

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		var user User

		// Scan columns into Comment + nested User
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&user.ID,
			&user.Username,
		)
		if err != nil {
			return nil, err
		}

		comment.User = user
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &comments, nil
}
