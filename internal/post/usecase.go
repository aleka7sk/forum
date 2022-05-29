package post

import (
	"context"
	"forum/models"
)

type UseCase interface {
	CreatePost(ctx context.Context, title, author, content string)
	GetAllPosts(ctx context.Context) []models.Post
	GetLikedPosts(ctx context.Context)
	GetUnlikedPosts(ctx context.Context)
	ParseToken(ctx context.Context, accessToken string) (*models.User, error)
}

// type CommentUseCase interface {
// 	CreateComment(ctx context.Context)
// 	UpdateComment(ctx context.Context)
// 	DeleteComment(ctx context.Context)
// }
