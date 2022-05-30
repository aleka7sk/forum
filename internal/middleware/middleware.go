package middleware

import (
	"context"
	"forum/internal/middleware/utils"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
)

type UserInfo struct {
	Rights   bool
	Id       int
	Username string
}

func Handle(next http.Handler, redis *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cook, err := r.Cookie("token")
		if err != nil {
			userInfo := UserInfo{}
			ctx := context.WithValue(r.Context(), "info", userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		user, err := utils.ParseToken(r.Context(), cook.Value)
		if err != nil {
			userInfo := UserInfo{}
			ctx := context.WithValue(r.Context(), "info", userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		redis_value, err := redis.Get(strconv.Itoa(user.Id)).Result()
		if err != nil {
			userInfo := UserInfo{}
			ctx := context.WithValue(r.Context(), "info", userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		if redis_value == cook.Value {
			userInfo := UserInfo{
				Rights:   true,
				Id:       user.Id,
				Username: user.Username,
			}
			ctx := context.WithValue(r.Context(), "info", userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			userInfo := UserInfo{}
			ctx := context.WithValue(r.Context(), "info", userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	}
}
