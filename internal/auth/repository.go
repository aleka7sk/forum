package auth

import (
	"context"
	"forum/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, username, password string) (*models.User, error)
	SaveRedis(token string, id int)
}
