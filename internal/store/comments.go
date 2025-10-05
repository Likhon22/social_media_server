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
        SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at,
       u.id, u.username, u.email, u.created_at, u.updated_at
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

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
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

func (s *CommentStore) CreateComment(ctx context.Context, comment *Comment) error {
	query := `
        INSERT INTO comments (post_id, user_id, content, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id
    `

	// If you want to get the generated ID
	err := s.db.QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content).Scan(&comment.ID)
	if err != nil {
		return err
	}

	return nil
}
