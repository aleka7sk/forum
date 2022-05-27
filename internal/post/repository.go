package post

import "context"

type Repository interface {
	CreatePost(ctx context.Context)
	GetAllPosts(ctx context.Context)
	GetLikedPosts(ctx context.Context)
	GetUnlikedPosts(ctx context.Context)
}
