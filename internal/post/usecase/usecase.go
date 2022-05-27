package post

import (
	"context"
	"forum/internal/post"
	"forum/models"
)

type service struct {
	repository post.Repository
}

func NewService(repository post.Repository) post.UseCase {
	return &service{repository: repository}
}

func (h *service) CreatePost(ctx context.Context) {
}

func (h *service) GetAllPosts(ctx context.Context) *[]models.Post {
	str1 := models.Post{Id: 1, Title: "Global Climate", Author: "ALish", Content: "Some news"}

	str2 := models.Post{Id: 2, Title: "News", Author: "Hasan", Content: "Nice boy"}
	return &[]models.Post{str1, str2}
}

func (h *service) GetLikedPosts(ctx context.Context) {
}

func (h *service) GetUnlikedPosts(ctx context.Context) {
}
