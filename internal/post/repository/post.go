package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/post"
	"forum/internal/postmodels"
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
	sqlStatement := `insert into posts (title, author, content, author_id) values ($1, $2, $3, $4)`
	_, err := pr.db.Exec(sqlStatement, post.Title, post.Author, post.Content, post.AuthorId)
	if err != nil {
		log.Fatalf("Insert error -> ss: %v", err)
	}
	return nil
}

func (pr Repo) GetAllPosts(ctx context.Context) []postmodels.Post {
	_, err := pr.db.Exec(`CREATE TABLE IF NOT EXISTS posts(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		content TEXT NOT NULL,
		author_id INTEGER NOT NULL
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
	posts_client := []postmodels.Post{}

	sqlQuery := `SELECT * FROM posts`
	rows, err := pr.db.Query(sqlQuery)
	if err != nil {
		log.Fatalf("Select query %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		count_like := 0
		count_dislike := 0
		post := models.Post{}
		if err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.Content, &post.AuthorId); err != nil {
			return posts_client
		}
		sqlQueryEmotionLike := `SELECT * FROM emotion WHERE post_id = $1 AND like = $2 AND dislike = $3`
		rows2, err := pr.db.Query(sqlQueryEmotionLike, post.Id, 1, 0)
		if err != nil {
			return posts_client
		}
		defer rows2.Close()
		for rows2.Next() {
			count_like++
		}

		sqlQueryEmotionDislike := `SELECT * FROM emotion WHERE post_id = $1 AND like = $2 AND dislike = $3`
		rows3, err := pr.db.Query(sqlQueryEmotionDislike, post.Id, 0, 1)
		if err != nil {
			return posts_client
		}
		defer rows3.Close()
		for rows3.Next() {
			count_dislike++
		}

		post_client := postmodels.Post{
			Id:       post.Id,
			Title:    post.Title,
			Author:   post.Author,
			Content:  post.Content,
			AuthorId: post.AuthorId,
			Likes:    count_like,
			Dislikes: count_dislike,
		}
		posts_client = append(posts_client, post_client)
	}

	return posts_client
}

func (pr Repo) GetLikedPosts(ctx context.Context, user_id int) ([]postmodels.Post, error) {
	sqlQuery := `SELECT * FROM emotion WHERE user_id = $1 AND like = 1`
	rows, err := pr.db.Query(sqlQuery, user_id)
	if err != nil {
		return nil, fmt.Errorf("repository GetLikedPosts() -> USER_ID: %d: %v", user_id, err)
	}

	posts := []postmodels.Post{}
	posts_id := []int{}
	for rows.Next() {
		var post_id int
		rows.Scan(nil, nil, nil, &post_id, nil)
		posts_id = append(posts_id, post_id)
	}
	rows.Close()

	sqlQueryLikedPosts := `SELECT * FROM posts WHERE id = $1`
	for _, elem := range posts_id {
		post := models.Post{}
		if err := pr.db.QueryRow(sqlQueryLikedPosts, elem); err != nil {
			fmt.Printf("Repository GetLikedPosts() -> Not such row with ID:%d error: %v", elem, err)
		}
	}

	return posts, nil
}

func (pr Repo) GetUnlikedPosts(ctx context.Context, user_id int) ([]postmodels.Post, error) {
	return nil, nil
}

func (pr Repo) GetPost(ctx context.Context, post_id, user_id int) (postmodels.Post, error) {
	var err error
	post := models.Post{}
	sqlQuery := `SELECT * FROM posts WHERE id = $1`
	if err = pr.db.QueryRow(sqlQuery, post_id).Scan(&post.Id, &post.Title, &post.Author, &post.Content, &post.AuthorId); err != nil {
		return postmodels.Post{}, fmt.Errorf("get Post() -> Post scan error: POST_id: %d :%v", post_id, err)
	}

	post_client := postmodels.Post{}

	sqlQueryEmotion := `SELECT * FROM emotion WHERE user_id = $1 AND post_id = $2`
	emotion := models.Emotion{}

	if err = pr.db.QueryRow(sqlQueryEmotion, user_id, post_id).Scan(&emotion.Id, &emotion.Like, &emotion.Dislike, &emotion.PostId, &emotion.UserId); err != nil {
		fmt.Printf("Get Post() -> Emotion scan error: USER_ID: %d, POST_ID: %d: %v\n", user_id, post_id, err)
	}

	post_client = postmodels.Post{Id: post.Id, Title: post.Title, Author: post.Author, Content: post.Content, AuthorId: post.AuthorId, Likes: emotion.Like, Dislikes: emotion.Dislike}
	return post_client, nil
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

func (pr Repo) CreateEmotion(ctx context.Context, post_id, user_id, like, dislike int) error {
	emotion := models.Emotion{
		Like:    like,
		Dislike: dislike,
		PostId:  post_id,
		UserId:  user_id,
	}

	sqlQuery := `SELECT * FROM emotion WHERE user_id = $1 AND post_id = $2`
	emotion_two := models.Emotion{}
	if err := pr.db.QueryRow(sqlQuery, user_id, post_id).Scan(&emotion_two.Id, &emotion_two.Like, &emotion_two.Dislike, &emotion_two.PostId, &emotion_two.UserId); err != nil {
		sqlStatement := `insert into emotion (like, dislike, post_id, user_id) values ($1, $2, $3, $4)`
		_, err = pr.db.Exec(sqlStatement, emotion.Like, emotion.Dislike, emotion.PostId, emotion.UserId)
		if err != nil {
			return fmt.Errorf("create Emotion() -> insert into emotion error: %v", err)
		}
		return nil
	}
	ctxt := context.Background()
	tx, err := pr.db.BeginTx(ctxt, nil)
	if err != nil {
		return err
	}
	if emotion_two.Like == emotion.Like && emotion_two.Dislike == emotion.Dislike {
		updateQuery := `UPDATE emotion SET like = $1, dislike = $2 WHERE user_id = $3 AND post_id = $4;`
		_, err = tx.ExecContext(ctxt, updateQuery, 0, 0, user_id, post_id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("create Emotion() -> update emotion error like: 0, dislike: 0: %v", err)
		}
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("create Emotion() -> commit update error like: 0,dislike: 0: %v", err)
		}
		return nil
	} else {
		updateQuery := `UPDATE emotion SET like = $1, dislike = $2 WHERE user_id = $3 AND post_id = $4;`
		_, err = tx.ExecContext(ctxt, updateQuery, like, dislike, user_id, post_id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("create Emotion() -> rollback update emotion error like: 1-1, dislike: 1-1: %v", err)
		}
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("create Emotion() -> commit update error like: 1-1,dislike: 1-1: %v", err)
		}
		return nil
	}
}
