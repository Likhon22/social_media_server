package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/likhon22/social/internal/store"
)

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()
	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			fmt.Println(err)
			return
		}

	}
	posts := generatePosts(200)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			fmt.Println(err)
			return
		}

	}
	comments := generateComments(200)
	for _, comment := range comments {
		if err := store.Comments.CreateComment(ctx, comment); err != nil {
			fmt.Println(err)
			return
		}

	}
	log.Println("seeding complete")

}
func generateUsers(num int) []*store.User {
	gofakeit.Seed(0)
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		pass := gofakeit.Password(true, true, true, true, false, 10)

		var password store.Password
		if err := password.Set(pass); err != nil {
			panic(err)
		}

		users[i] = &store.User{
			Username:  gofakeit.Username(),
			Email:     gofakeit.Email(),
			Password:  password, // ✅ hashed and ready
			CreatedAt: gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
			UpdatedAt: time.Now(),
		}
	}
	return users

}

func generatePosts(num int) []*store.Post {
	gofakeit.Seed(0) // ensure random data on each run
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		posts[i] = &store.Post{

			Title:     gofakeit.Sentence(3),
			Content:   gofakeit.Paragraph(1, 3, 10, " "),
			UserID:    int64(gofakeit.Number(20, 80)),
			Tags:      []string{gofakeit.Word(), gofakeit.Word()},
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}
	}

	return posts
}
func generateComments(num int) []*store.Comment {
	gofakeit.Seed(0)
	comments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		comments[i] = &store.Comment{
			PostID:    gofakeit.Number(20, 180),
			UserID:    gofakeit.Number(20, 80),
			Content:   gofakeit.Sentence(10),
			CreatedAt: gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
			UpdatedAt: time.Now(),
		}
	}

	return comments
}
