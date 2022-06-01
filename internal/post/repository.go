package post

import (
	"context"
	"forum/internal/postmodels"
	"forum/models"
)

type Repository interface {
	CreatePost(ctx context.Context, title, author, content string, category, author_id int) error
	GetAllPosts(ctx context.Context) []postmodels.Post
	GetPost(ctx context.Context, post_id, user_id int) (postmodels.Post, error)
	GetLikedPosts(ctx context.Context, user_id int) ([]postmodels.Post, error)
	GetUnlikedPosts(ctx context.Context, user_id int) ([]postmodels.Post, error)
	GetMyPosts(ctx context.Context, author_id string) []models.Post
	CreateVote(ctx context.Context, post_id, user_id, condition int) error
}
