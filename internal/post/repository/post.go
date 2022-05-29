package repository

import (
	"context"
	"database/sql"
	"forum/internal/post"
	"forum/models"
	"log"
)

type Repo struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) post.Repository {
	return &Repo{
		db: db,
	}
}

func (pr Repo) CreatePost(ctx context.Context, title, author, content string) error {
	post := models.Post{
		Title:   title,
		Author:  author,
		Content: content,
	}

	_, err := pr.db.Exec(`CREATE TABLE IF NOT EXISTS posts(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		heading TEXT,
		author TEXT,
		content TEXT
	  );`)
	if err != nil {
		log.Fatalf("cannot exec file: %v", err.Error())
	}
	sqlStatement := `insert into posts (heading, author, content) values ($1, $2, $3)`
	_, err = pr.db.Exec(sqlStatement, post.Title, post.Author, post.Content)
	if err != nil {
		log.Fatalf("Insert error -> : %v", err)
	}

	return nil
}

func (pr Repo) GetAllPosts(ctx context.Context) []models.Post {
	posts := []models.Post{}
	sqlQuery := `SELECT * FROM posts`
	rows, err := pr.db.Query(sqlQuery)
	if err != nil {
		log.Fatalf("Select query %v", err)
	}
	for rows.Next() {
		post := models.Post{}
		if err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.Content); err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	return posts
}

func (pr Repo) GetLikedPosts(ctx context.Context) {
}

func (pr Repo) GetUnlikedPosts(ctx context.Context) {
}
