package post

import (
	"context"
	"forum/internal/postmodels"
	"forum/models"
)

type Repository interface {
	CreatePost(ctx context.Context, title, author, content, author_id string) error
	GetAllPosts(ctx context.Context) []postmodels.Post
	GetPost(ctx context.Context, post_id, user_id int) (postmodels.Post, error)
	GetLikedPosts(ctx context.Context, user_id int) ([]postmodels.Post, error)
	GetUnlikedPosts(ctx context.Context, user_id int) ([]postmodels.Post, error)
	GetMyPosts(ctx context.Context, author_id string) []models.Post
	CreateEmotion(ctx context.Context, post_id, user_id, like, dislike int) error
}
