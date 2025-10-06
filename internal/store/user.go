package store

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  Password  `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Password struct {
	text *string
	hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err

	}
	p.hash = hash

	p.text = &text
	return nil

}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.hash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetUsers(ctx context.Context) (*[]User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `SELECT id, username, email, created_at, updated_at 
	          FROM users 
	          ORDER BY created_at ASC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `SELECT id, username, email,  created_at, updated_at 
	          FROM users WHERE email = $1`
	user := User{}

	err := s.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (s *UserStore) GetUserById(ctx context.Context, id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `SELECT id, username, email,  created_at, updated_at 
	          FROM users WHERE id = $1`
	user := User{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		if err := s.createUserInvitation(ctx, tx, token, user.ID, invitationExp); err != nil {
			return err
		}
		return nil
	})

}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, userId int64, invitationExp time.Duration) error {
	query := `INSERT into user_invitations (token,user_id,expiry) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, token, userId, time.Now().Add(invitationExp))
	if err != nil {
		return err

	}
	return nil

}
