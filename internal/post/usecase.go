package post

import (
	"context"
	"forum/models"
)

type UseCase interface {
	CreatePost(ctx context.Context, title, author, content, author_id string)
	GetAllPosts(ctx context.Context) []models.Post
	GetPost(ctx context.Context, id string) models.Post
	GetLikedPosts(ctx context.Context)
	GetUnlikedPosts(ctx context.Context)
	GetMyPosts(ctx context.Context, author_id string) []models.Post
	CreateEmotion(ctx context.Context, post string, user_id int, like, dislike bool) error
}

// type CommentUseCase interface {
// 	CreateComment(ctx context.Context)
// 	UpdateComment(ctx context.Context)
// 	DeleteComment(ctx context.Context)
// }
