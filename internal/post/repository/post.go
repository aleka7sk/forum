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
	CreateTables(db)
	return &Repo{
		db: db,
	}
}

func CreateTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS category(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		title TEXT NOT NULL
	  );`)
	if err != nil {
		log.Fatalf("cannot create category: %v", err.Error())
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS post(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		content TEXT NOT NULL,
		author_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		likes INTEGER NOT NULL,
		dislikes INTEGER NOT NULL,
		FOREIGN KEY(author_id) REFERENCES users(id),
		FOREIGN KEY(category_id) REFERENCES category(id)
	  );`)
	if err != nil {
		log.Fatalf("cannot create post: %v", err.Error())
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS vote(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		condition INTEGER,
		post_id INTEGER,
		user_id INTEGER,
		FOREIGN KEY(post_id) REFERENCES post(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	  );`)
	if err != nil {
		log.Fatalf("cannot create vote: %v", err.Error())
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comment(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		author  TEXT NOT NULL, 
		content TEXT NOT NULL,
		post_id INTEGER,
		user_id INTEGER,
		FOREIGN KEY(post_id) REFERENCES post(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	  );`)
	if err != nil {
		log.Fatalf("cannot create vote: %v", err.Error())
	}
	query := `INSERT into category (title) VALUES ('WEB')`
	_, err = db.Exec(query)
	fmt.Println(err)
}

func (pr Repo) CreatePost(ctx context.Context, title, author, content string, category_id, author_id int) error {
	post := models.Post{
		Title:      title,
		Author:     author,
		Content:    content,
		AuthorId:   author_id,
		CategoryId: category_id,
		Likes:      0,
		Dislikes:   0,
	}
	sqlStatement := `insert into post (title, author, content, author_id, category_id, likes, dislikes) values ($1, $2, $3, $4, $5, $6, $7)`
	fmt.Println(1)
	_, err := pr.db.Exec(sqlStatement, post.Title, post.Author, post.Content, post.AuthorId, post.CategoryId, post.Likes, post.Dislikes)
	if err != nil {
		log.Fatalf("Insert error -> ss: %v", err)
	}
	return nil
}

func (pr Repo) GetAllPosts(ctx context.Context) []postmodels.Post {
	sqlQuery := `SELECT * FROM post`
	rows, err := pr.db.Query(sqlQuery)
	if err != nil {
		log.Fatalf("Select query %v", err)
	}
	defer rows.Close()
	posts_client := countingLikesDislikes(pr.db, rows)
	// for rows.Next() {
	// 	count_like := 0
	// 	count_dislike := 0
	// 	post := models.Post{}
	// 	if err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.Content, &post.AuthorId); err != nil {
	// 		return posts_client
	// 	}
	// 	sqlQueryEmotionLike := `SELECT * FROM emotion WHERE post_id = $1 AND like = $2 AND dislike = $3`
	// 	rows2, err := pr.db.Query(sqlQueryEmotionLike, post.Id, 1, 0)
	// 	if err != nil {
	// 		return posts_client
	// 	}
	// 	defer rows2.Close()
	// 	for rows2.Next() {
	// 		count_like++
	// 	}

	// 	sqlQueryEmotionDislike := `SELECT * FROM emotion WHERE post_id = $1 AND like = $2 AND dislike = $3`
	// 	rows3, err := pr.db.Query(sqlQueryEmotionDislike, post.Id, 0, 1)
	// 	if err != nil {
	// 		return posts_client
	// 	}
	// 	defer rows3.Close()
	// 	for rows3.Next() {
	// 		count_dislike++
	// 	}

	// 	post_client := postmodels.Post{
	// 		Id:       post.Id,
	// 		Title:    post.Title,
	// 		Author:   post.Author,
	// 		Content:  post.Content,
	// 		AuthorId: post.AuthorId,
	// 		Likes:    count_like,
	// 		Dislikes: count_dislike,
	// 	}
	// 	posts_client = append(posts_client, post_client)
	// }

	return posts_client
}

func (pr Repo) GetLikedPosts(ctx context.Context, user_id int) ([]postmodels.Post, error) {
	// sqlQuery := `SELECT * FROM emotion WHERE user_id = $1 AND like = 1`
	// rows, err := pr.db.Query(sqlQuery, user_id)
	// if err != nil {
	// 	return nil, fmt.Errorf("repository GetLikedPosts() -> USER_ID: %d: %v", user_id, err)
	// }

	// posts_id := []int{}
	// for rows.Next() {
	// 	var post_id int
	// 	rows.Scan(nil, nil, nil, &post_id, nil)
	// 	posts_id = append(posts_id, post_id)
	// }
	// rows.Close()
	// posts := []postmodels.Post{}
	// sqlQueryLikedPosts := `SELECT * FROM posts WHERE id = $1`
	// // for _, elem := range posts_id {
	// // 	post := models.Post{}
	// // 	if err := pr.db.QueryRow(sqlQueryLikedPosts, elem); err != nil {
	// // 		fmt.Printf("Repository GetLikedPosts() -> Not such row with ID:%d error: %v", elem, err)
	// // 	}
	// // }
	// return nil, nil
	return nil, nil
}

func (pr Repo) GetUnlikedPosts(ctx context.Context, user_id int) ([]postmodels.Post, error) {
	return nil, nil
}

func (pr Repo) GetPost(ctx context.Context, post_id, user_id int) (postmodels.Post, error) {
	var err error

	sqlQuery := `SELECT * FROM post WHERE id = $1`
	rows, err := pr.db.Query(sqlQuery, post_id)
	post_client := countingLikesDislikes(pr.db, rows)
	if err != nil {
		return postmodels.Post{}, fmt.Errorf("get Post() -> Post scan error: POST_id: %d :%v", post_id, err)
	}

	sqlQueryEmotion := `SELECT * FROM vote WHERE user_id = $1 AND post_id = $2`
	vote := models.Vote{}

	if err = pr.db.QueryRow(sqlQueryEmotion, user_id, post_id).Scan(&vote.Id, &vote.Condition, &vote.PostId, &vote.UserId); err != nil {
		fmt.Printf("Get Post() -> Emotion scan error: USER_ID: %d, POST_ID: %d: %v\n", user_id, post_id, err)
	}
	if post_client != nil {
		post_client[0].Condition = vote.Condition

		return post_client[0], nil
	}
	return postmodels.Post{}, nil
}

func (pr Repo) GetMyPosts(ctx context.Context, author_id string) []models.Post {
	posts := []models.Post{}
	sqlQuery := `SELECT * FROM post WHERE author_id = $1`
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

func (pr Repo) CreateVote(ctx context.Context, post_id, user_id, condition int) error {
	vote := models.Vote{
		Condition: condition,
		PostId:    post_id,
		UserId:    user_id,
	}
	sqlQuery := `SELECT * FROM vote WHERE user_id = $1 AND post_id = $2`
	vote_two := models.Vote{}
	if err := pr.db.QueryRow(sqlQuery, user_id, post_id).Scan(&vote_two.Id, &vote_two.Condition, &vote_two.PostId, &vote_two.UserId); err != nil {
		sqlStatement := `insert into vote (condition, post_id, user_id) values ($1, $2, $3)`
		_, err = pr.db.Exec(sqlStatement, vote.Condition, vote.PostId, vote.UserId)
		if err != nil {
			return fmt.Errorf("CreateVote() -> insert into emotion error: %v", err)
		}
		ctxt := context.Background()
		tx, err := pr.db.BeginTx(ctxt, nil)
		if err != nil {
			return err
		}
		incrementLikesDislikes(ctxt, tx, vote.Condition, post_id)
		return nil
	}
	ctxt := context.Background()
	tx, err := pr.db.BeginTx(ctxt, nil)
	if err != nil {
		return err
	}
	if vote_two.Condition == vote.Condition {
		updateQuery := `UPDATE vote SET condition = 0 WHERE user_id = $1 AND post_id = $2;`
		_, err = tx.ExecContext(ctxt, updateQuery, user_id, post_id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("CreateVote() -> update emotion error like: 0, dislike: 0: %v", err)
		}
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("CreateVote() -> commit update error like: 0,dislike: 0: %v", err)
		}
		return nil
	} else {
		updateQuery := `UPDATE vote SET condition = $1 WHERE user_id = $2 AND post_id = $3;`
		_, err = tx.ExecContext(ctxt, updateQuery, condition, user_id, post_id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("CreateVote() -> rollback update emotion error like: 1-1, dislike: 1-1: %v", err)
		}
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("CreateVote() -> commit update error like: 1-1,dislike: 1-1: %v", err)
		}
		return nil
	}
}

func countingLikesDislikes(db *sql.DB, rows *sql.Rows) []postmodels.Post {
	posts_client := []postmodels.Post{}

	for rows.Next() {
		count_like := 0
		count_dislike := 0
		post := models.Post{}

		if err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.Content, &post.AuthorId, &post.CategoryId, &post.Likes, &post.Dislikes); err != nil {
			return posts_client
		}
		sqlQueryEmotionLike := `SELECT * FROM vote WHERE post_id = $1 AND condition = 1`
		rows2, err := db.Query(sqlQueryEmotionLike, post.Id)
		if err != nil {
			fmt.Printf("Not likes in post: %v", err)
		}

		defer rows2.Close()
		for rows2.Next() {
			count_like++
		}

		sqlQueryEmotionDislike := `SELECT * FROM vote WHERE post_id = $1 AND condition = 2`
		rows3, err := db.Query(sqlQueryEmotionDislike, post.Id)
		if err != nil {
			fmt.Printf("Not dislikes in post: %v", err)
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
	rows.Close()
	return posts_client
}

func incrementLikesDislikes(ctxt context.Context, tx *sql.Tx, condition, post_id int) error {
	var err error
	if err != nil {
		return err
	}
	if condition == 1 {
		selectPostQuery := `UPDATE post SET likes = (SELECT likes FROM post WHERE id = $1) + 1 WHERE id = $2`
		_, err = tx.ExecContext(ctxt, selectPostQuery, post_id, post_id)
		fmt.Println(err)
	} else if condition == 2 {
		selectPostQuery := `UPDATE post SET dislikes = (SELECT dislikes FROM post WHERE id = $1) + 1 WHERE id = $2`
		_, err = tx.ExecContext(ctxt, selectPostQuery, post_id, post_id)
	} else if condition == 3 {
		selectPostQuery := `UPDATE post SET likes = (SELECT likes FROM post WHERE id = $1) - 1 WHERE id = $2`
		_, err = tx.ExecContext(ctxt, selectPostQuery, post_id, post_id)
	} else if condition == 4 {
		selectPostQuery := `UPDATE post SET dislikes = (SELECT dislikes FROM post WHERE id = $1) - 1 WHERE id = $2`
		_, err = tx.ExecContext(ctxt, selectPostQuery, post_id, post_id)
	}
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("create Emotion() -> update posts counter: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("create Emotion() -> commit update posts counter: %v", err)
	}
	return nil
}
