package http

import (
	"forum/internal/middleware"
	"forum/internal/post"
	"net/http"
)

func RegisterHTTPEndpoints(router *http.ServeMux, auc post.UseCase) {
	h := NewHandler(auc)
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))
	router.HandleFunc("/", middleware.Handle(http.HandlerFunc(h.MainPage)))
	router.HandleFunc("/logout", h.Logout)
	router.HandleFunc("/create-post", middleware.Handle(http.HandlerFunc(h.Create)))
	router.HandleFunc("/article", middleware.Handle(http.HandlerFunc(h.Post)))
}
