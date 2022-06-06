package http

import (
	"forum/internal/auth"
	"forum/internal/middleware"
	"net/http"

	"github.com/go-redis/redis"
)

func RegisterHTTPEndpoints(router *http.ServeMux, auc auth.UseCase, redis *redis.Client) {
	h := NewHandler(auc)
	router.HandleFunc("/sign-up", middleware.Handle(http.HandlerFunc(h.SignUp), redis))
	router.HandleFunc("/sign-in", middleware.Handle(http.HandlerFunc(h.SignIn), redis))
}
