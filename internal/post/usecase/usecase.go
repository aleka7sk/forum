package post

import (
	"context"
	"fmt"
	"forum/internal/auth"
	"forum/internal/post"
	"forum/models"
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

func (h *service) CreatePost(ctx context.Context, title, author, content string) {
	h.repository.CreatePost(ctx, title, author, content)
}

func (h *service) GetAllPosts(ctx context.Context) []models.Post {
	posts := h.repository.GetAllPosts(ctx)
	return posts
}

func (h *service) GetLikedPosts(ctx context.Context) {
}

func (h *service) GetUnlikedPosts(ctx context.Context) {
}

func (h *service) ParseToken(ctx context.Context, accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return h.signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, auth.ErrInvalidAccessToken
}
