package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	Title     string    `json:"title" db:"title"`
	Tags      []string  `json:"tags" db:"tags"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Comments  []Comment `json:"comments" db:"comments"`
	User      User      `json:"Users" db:"Users"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PostWithMetaData struct {
	Post
	CommentCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

// @Summary		Create a new post
// @Description	Creates a new post
// @Tags			Posts
// @Accept			json
// @Produce		json
// @Param			post	body		store.Post	true	"Post information"
// @Success		201		{object}	store.Post
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Router			/posts [post]
func (s *PostStore) Create(ctx context.Context, post *Post) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `INSERT INTO posts (content, title, tags, user_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, pq.Array(post.Tags), post.UserID).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// @Summary		Get all posts
// @Description	Retrieves a list of all posts
// @Tags			Posts
// @Produce		json
// @Success		200	{array}		store.Post
// @Failure		400	{object}	error
// @Failure		500	{object}	error
// @Router			/posts [get]
func (s *PostStore) GetAll(ctx context.Context) ([]*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `SELECT id, content, title, tags, user_id, created_at, updated_at FROM posts ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(&post.ID, &post.Content, &post.Title, pq.Array(&post.Tags), &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

// @Summary		Get a post by ID
// @Description	Retrieves a post and its comments by post ID
// @Tags			Posts
// @Produce		json
// @Param			postId	path		int	true	"Post ID"
// @Success		200		{object}	store.Post
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Failure		500		{object}	error
// @Router			/posts/{postId} [get]
func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id, content, title, tags, user_id, created_at, updated_at FROM posts WHERE id = $1`
	post := &Post{}
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Content, &post.Title, pq.Array(&post.Tags), &post.UserID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return post, nil
}

// @Summary		Delete a post by ID
// @Description	Deletes a post by its ID
// @Tags			Posts
// @Produce		json
// @Param			postId	path		int		true	"Post ID"
// @Success		200		{string}	string	"post deleted successfully"
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Router			/posts/{postId} [delete]
func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1` // Use ? if MySQL
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	// Execute the query
	result, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	// Optional: check if a row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no post found with id %d", postID)
	}

	return nil
}

// @Summary		Update a post
// @Description	Updates a post's content, title, or tags by post ID
// @Tags			Posts
// @Accept			json
// @Produce		json
// @Param			postId	path		int			true	"Post ID"
// @Param			post	body		store.Post	true	"Post data to update"
// @Success		200		{object}	store.Post
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Router			/posts/{postId} [patch]
func (s *PostStore) Update(ctx context.Context, postID int64, post *Post) error {
	setParts := []string{}
	args := []interface{}{}
	i := 1
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	if post.Title != "" {
		setParts = append(setParts, fmt.Sprintf("title = $%d", i))
		args = append(args, post.Title)
		i++
	}

	if post.Content != "" {
		setParts = append(setParts, fmt.Sprintf("content = $%d", i))
		args = append(args, post.Content)
		i++
	}

	if post.Tags != nil {
		setParts = append(setParts, fmt.Sprintf("tags = $%d", i))
		args = append(args, pq.Array(post.Tags))
		i++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Always update UpdatedAt
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", i))
	args = append(args, time.Now())
	i++

	// Build final query
	query := fmt.Sprintf("UPDATE posts SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), i)
	args = append(args, postID)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no post found with id %d", postID)
	}

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userId int64, fq PaginatedFeedQuery) (*[]PostWithMetaData, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
SELECT 
    p.id,
    p.user_id,
    p.title,
    p.content,
    p.created_at,
    p.updated_at,
    p.tags,
    u.username,
    COUNT(c.id) AS comments_count
FROM posts p
LEFT JOIN comments c ON c.post_id = p.id
LEFT JOIN users u ON u.id = p.user_id
WHERE (p.user_id = $1 OR p.user_id IN (
    SELECT f.follower_id
    FROM followers f
    WHERE f.user_id = $1
))
`

	args := []interface{}{userId}

	// Add search filter only if search is provided
	if fq.Search != "" {
		query += " AND (p.title ILIKE '%' || $2 || '%' OR p.content ILIKE '%' || $2 || '%')"
		args = append(args, fq.Search)
	}

	query += fmt.Sprintf(" GROUP BY p.id, u.username ORDER BY p.created_at %s LIMIT $%d OFFSET $%d",
		fq.Sort, len(args)+1, len(args)+2)

	args = append(args, fq.Limit, fq.Offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []PostWithMetaData{}
	for rows.Next() {
		post := PostWithMetaData{}
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &posts, nil
}
