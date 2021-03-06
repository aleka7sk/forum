package auth

import (
	"context"
	"crypto/sha1"
	"fmt"
	"forum/internal/auth"
	"forum/models"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type service struct {
	repository     auth.Repository
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

type AuthClaims struct {
	jwt.StandardClaims
	User *models.User `json:"user"`
}

func NewService(repository auth.Repository, hashSalt string, signingKey []byte, tokenTTLSecond time.Duration) *service {
	return &service{
		repository:     repository,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSecond,
	}
}

func (s *service) SignUp(ctx context.Context, username, password string) error {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(s.hashSalt))

	user := &models.User{
		Username: username,
		Password: fmt.Sprintf("%x", pwd.Sum(nil)),
	}

	return s.repository.CreateUser(ctx, user)
}

func (s *service) SignIn(ctx context.Context, username, password string) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(s.hashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := s.repository.GetUser(ctx, username, password)
	if err != nil {
		fmt.Println("UserNotFound")
		return "", auth.ErrUserNotFound
	}

	claims := AuthClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(s.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	our_token, err := token.SignedString(s.signingKey)
	s.repository.SaveRedis(our_token, user.Id)
	return our_token, err
}
