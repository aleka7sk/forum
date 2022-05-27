package repository

import (
	"context"
	"database/sql"
	"forum/internal/post"
)

type Repo struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) post.Repository {
	return &Repo{
		db: db,
	}
}

func (pr Repo) CreatePost(ctx context.Context) {
}

func (pr Repo) GetAllPosts(ctx context.Context) {
}

func (pr Repo) GetLikedPosts(ctx context.Context) {
}

func (pr Repo) GetUnlikedPosts(ctx context.Context) {
}
