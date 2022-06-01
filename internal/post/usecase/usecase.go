package post

import (
	"context"
	"fmt"
	"forum/internal/post"
	"forum/internal/postmodels"
	"forum/models"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type service struct {
	repository     post.Repository
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

type AuthClaims struct {
	jwt.StandardClaims
	User *models.User `json:"user"`
}

func NewService(repository post.Repository, hashSalt string, signingKey []byte, tokenTTLSecond time.Duration) *service {
	return &service{
		repository:     repository,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSecond,
	}
}

func (h *service) CreatePost(ctx context.Context, title, author, content, author_id string) {
	err := h.repository.CreatePost(ctx, title, author, content, author_id)
	if err != nil {
		log.Printf("Error")
	}
}

func (h *service) GetAllPosts(ctx context.Context) []postmodels.Post {
	posts := h.repository.GetAllPosts(ctx)
	return posts
}

func (h *service) GetPost(ctx context.Context, post_id, user_id int) (postmodels.Post, error) {
	post, err := h.repository.GetPost(ctx, post_id, user_id)
	fmt.Printf("Usecase GetPost() -> POST_ID: %d, USER_ID: %d\n", post_id, user_id)
	if err != nil {
		return postmodels.Post{}, err
	}
	return post, nil
}

func (h *service) GetLikedPosts(ctx context.Context, user_id int) ([]postmodels.Post, error) {
	return nil, nil
}

func (h *service) GetUnlikedPosts(ctx context.Context, user_id int) ([]postmodels.Post, error) {
	return nil, nil
}

func (h *service) GetMyPosts(ctx context.Context, author_id string) []models.Post {
	posts := h.repository.GetMyPosts(ctx, author_id)
	return posts
}

func (h *service) CreateEmotion(ctx context.Context, post_id, user_id int, like, dislike bool) error {
	var LikeInt int
	var DisLikeInt int
	if like {
		LikeInt = 1
		DisLikeInt = 0
	}
	if dislike {
		LikeInt = 0
		DisLikeInt = 1
	}
	err := h.repository.CreateEmotion(ctx, post_id, user_id, LikeInt, DisLikeInt)
	if err != nil {
		return err
	}
	return nil
}
