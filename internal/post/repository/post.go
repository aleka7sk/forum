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
	_, err = pr.db.Exec(`CREATE TABLE IF NOT EXISTS emotion(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		like  INTEGER, 
		dislike INTEGER,
		post_id INTEGER,
		user_id INTEGER,
		FOREIGN KEY(post_id) REFERENCES posts(author_id),
		FOREIGN KEY(user_id) REFERENCES users(id)
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

func (pr Repo) CreateEmotion(ctx context.Context, post, user_id, like, dislike int) error {
	emotion := models.Emotion{
		Like:    like,
		Dislike: dislike,
		PostId:  post,
		UserId:  user_id,
	}
	var Id int
	sqlQuery := `SELECT * FROM emotion WHERE user_id = $1 AND post_id = $2`
	rows, err := pr.db.Query(sqlQuery, user_id, post)
	if err != nil {
		return err
	}

	for rows.Next() {
		emotion_two := models.Emotion{}
		if err := rows.Scan(&Id, &emotion_two.Like, &emotion_two.Dislike, &emotion_two.PostId, &emotion_two.UserId); err != nil {
			return err
		}

		if emotion_two == emotion {
			rows.Close()
			updateQuery := `UPDATE emotion SET like = $1, dislike = $2 WHERE user_id = $3 AND post_id = $4;`
			_, err = pr.db.Exec(updateQuery, 0, 0, post, user_id)
			if err != nil {
				return err
			}

			return nil
		} else {
			rows.Close()
			updateQuery := `UPDATE emotion SET like = $1, dislike = $2 WHERE user_id = $3 AND post_id = $4;`
			_, err = pr.db.Exec(updateQuery, like, dislike, post, user_id)
			if err != nil {
				return err
			}

			return nil
		}

	}
	sqlStatement := `insert into emotion (like, dislike, post_id, user_id) values ($1, $2, $3, $4)`

	_, err = pr.db.Exec(sqlStatement, emotion.Like, emotion.Dislike, emotion.PostId, emotion.UserId)
	if err != nil {
		return err
	}

	return nil
}
