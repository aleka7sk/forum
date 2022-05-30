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

func (pr Repo) CreatePost(ctx context.Context, title, author, content, author_id string) error {
	post := models.Post{
		Title:    title,
		Author:   author,
		Content:  content,
		AuthorId: author_id,
	}

	sqlStatement := `insert into posts (heading, author, content, author_id) values ($1, $2, $3, $4)`
	_, err := pr.db.Exec(sqlStatement, post.Title, post.Author, post.Content, post.AuthorId)
	if err != nil {
		log.Fatalf("Insert error -> ss: %v", err)
	}

	return nil
}

func (pr Repo) GetAllPosts(ctx context.Context) []models.Post {
	_, err := pr.db.Exec(`CREATE TABLE IF NOT EXISTS posts(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		heading TEXT,
		author TEXT,
		content TEXT,
		author_id INTEGER
	  );`)
	if err != nil {
		log.Fatalf("cannot exec file: %v", err.Error())
	}
	posts := []models.Post{}
	sqlQuery := `SELECT * FROM posts`
	rows, err := pr.db.Query(sqlQuery)
	if err != nil {
		log.Fatalf("Select query %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		post := models.Post{}
		if err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.Content, &post.AuthorId); err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	return posts
}

func (pr Repo) GetLikedPosts(ctx context.Context) {
}

func (pr Repo) GetPost(ctx context.Context, id string) models.Post {
	post := models.Post{}
	sqlQuery := `SELECT * FROM posts WHERE id = $1`
	rows, err := pr.db.Query(sqlQuery, id)
	if err != nil {
		log.Fatalf("Select query %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.Content, &post.AuthorId); err != nil {
			panic(err)
		}
	}

	return post
}

func (pr Repo) GetUnlikedPosts(ctx context.Context) {
}

func (pr Repo) GetMyPosts(ctx context.Context, author_id string) []models.Post {
	posts := []models.Post{}
	sqlQuery := `SELECT * FROM posts WHERE author_id = $1`
	rows, err := pr.db.Query(sqlQuery, author_id)
	if err != nil {
		log.Fatalf("Select query %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		post := models.Post{}
		if err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.Content, &post.AuthorId); err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	return posts
}
