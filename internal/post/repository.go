package post

import (
	"context"
	"forum/models"
)

type Repository interface {
	CreatePost(ctx context.Context, title, author, content string) error
	GetAllPosts(ctx context.Context) []models.Post
	GetLikedPosts(ctx context.Context)
	GetUnlikedPosts(ctx context.Context)
}
