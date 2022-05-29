package utils

import (
	"context"
	"fmt"
	"forum/internal/auth"
	"forum/models"
	"log"

	"github.com/dgrijalva/jwt-go/v4"
)

type AuthClaims struct {
	jwt.StandardClaims
	User *models.User `json:"user"`
}

func ParseToken(ctx context.Context, accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("signing_key"), nil
	})
	if err != nil {
		log.Fatalf("Parse token error: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, auth.ErrInvalidAccessToken
}
