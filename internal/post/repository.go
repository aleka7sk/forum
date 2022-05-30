package post

import (
	"context"
	"forum/models"
)

type Repository interface {
	CreatePost(ctx context.Context, title, author, content, author_id string) error
	GetAllPosts(ctx context.Context) []models.Post
	GetPost(ctx context.Context, id string) models.Post
	GetLikedPosts(ctx context.Context)
	GetUnlikedPosts(ctx context.Context)
	GetMyPosts(ctx context.Context, author_id string) []models.Post
}
