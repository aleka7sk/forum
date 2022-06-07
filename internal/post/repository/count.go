package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/postmodels"
	"forum/models"
	"log"
)

func postchecker(db *sql.DB, rows *sql.Rows) []postmodels.Post {
	array := make([]int, 0)
	for rows.Next() {
		vote := models.Vote{}
		if err := rows.Scan(&vote.Id, &vote.Condition, &vote.PostId, &vote.UserId); err != nil {
			fmt.Print(err)
		}
		array = append(array, vote.PostId)
	}
	posts := []postmodels.Post{}

	sqlQuery := `SELECT * from post WHERE id = $1`
	for _, id := range array {
		post := postmodels.Post{}
		count_likes, count_dislikes := counting(db, id)
		if err := db.QueryRow(sqlQuery, id).Scan(&post.Id, &post.Title, &post.Author, &post.Content, &post.AuthorId, &post.Condition); err != nil {
			log.Printf("Postchecker() error: %v", err)
		}
		post.Likes = count_likes
		post.Dislikes = count_dislikes
		posts = append(posts, post)

	}
	return posts
}

func countingLikesDislikes(db *sql.DB, rows *sql.Rows) []postmodels.Post {
	posts_client := []postmodels.Post{}

	for rows.Next() {

		post := models.Post{}

		if err := rows.Scan(&post.Id, &post.Title, &post.Author, &post.Content, &post.AuthorId, &post.CategoryId); err != nil {
			return posts_client
		}
		count_like, count_dislike := counting(db, post.Id)
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

func counting(db *sql.DB, post_id int) (int, int) {
	count_like := 0
	count_dislike := 0
	sqlQueryEmotionLike := `SELECT * FROM vote WHERE post_id = $1 AND condition = 1`
	rows2, err := db.Query(sqlQueryEmotionLike, post_id)
	if err != nil {
		fmt.Printf("Not likes in post: %v", err)
	}

	defer rows2.Close()
	for rows2.Next() {
		count_like++
	}

	sqlQueryEmotionDislike := `SELECT * FROM vote WHERE post_id = $1 AND condition = 2`
	rows3, err := db.Query(sqlQueryEmotionDislike, post_id)
	if err != nil {
		fmt.Printf("Not dislikes in post: %v", err)
	}
	defer rows3.Close()
	for rows3.Next() {
		count_dislike++
	}
	return count_like, count_dislike
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
