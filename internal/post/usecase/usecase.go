package post

import (
	"context"
	"fmt"
	"forum/internal/post"
	"forum/models"
	"log"
	"strconv"
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

func (h *service) GetAllPosts(ctx context.Context) []models.Post {
	posts := h.repository.GetAllPosts(ctx)
	return posts
}

func (h *service) GetPost(ctx context.Context, id string) models.Post {
	post := h.repository.GetPost(ctx, id)
	return post
}

func (h *service) GetLikedPosts(ctx context.Context) {
}

func (h *service) GetUnlikedPosts(ctx context.Context) {
}

func (h *service) GetMyPosts(ctx context.Context, author_id string) []models.Post {
	posts := h.repository.GetMyPosts(ctx, author_id)
	return posts
}

func (h *service) CreateEmotion(ctx context.Context, post string, user_id int, like, dislike bool) error {
	PostId, err := strconv.Atoi(post)
	fmt.Println(like)
	fmt.Println(dislike)
	if err != nil {
		return err
	}
	var LikeInt int
	var DisLikeInt int
	if like == true {
		LikeInt = 1
		DisLikeInt = 0
	}
	if dislike == true {
		LikeInt = 0
		DisLikeInt = 1
	}
	err = h.repository.CreateEmotion(ctx, PostId, user_id, LikeInt, DisLikeInt)
	if err != nil {
		return err
	}
	return nil
}
