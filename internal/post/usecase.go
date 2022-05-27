package post

import (
	"context"
	"forum/models"
)

type UseCase interface {
	CreatePost(ctx context.Context)
	GetAllPosts(ctx context.Context) *[]models.Post
	GetLikedPosts(ctx context.Context)
	GetUnlikedPosts(ctx context.Context)
}

// type CommentUseCase interface {
// 	CreateComment(ctx context.Context)
// 	UpdateComment(ctx context.Context)
// 	DeleteComment(ctx context.Context)
// }
