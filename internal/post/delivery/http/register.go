package http

import (
	"forum/internal/middleware"
	"forum/internal/post"
	"net/http"

	"github.com/go-redis/redis"
)

func RegisterHTTPEndpoints(router *http.ServeMux, auc post.UseCase, redis *redis.Client) {
	h := NewHandler(auc)
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))
	router.HandleFunc("/logout", h.Logout)
	router.HandleFunc("/", middleware.Handle(http.HandlerFunc(h.MainPage), redis))
	router.HandleFunc("/create-post", middleware.Handle(http.HandlerFunc(h.Create), redis))
	router.HandleFunc("/article/", middleware.Handle(http.HandlerFunc(h.Post), redis))
	router.HandleFunc("/my-posts", middleware.Handle(http.HandlerFunc(h.MyPosts), redis))
	router.HandleFunc("/my-posts", middleware.Handle(http.HandlerFunc(h.LikedPosts), redis))
}
