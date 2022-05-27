package http

import (
	"forum/internal/auth"
	"net/http"
)

func RegisterHTTPEndpoints(router *http.ServeMux, auc auth.UseCase) {
	fs := http.FileServer(http.Dir("./static"))
	h := NewHandler(auc)
	router.Handle("/static/", http.StripPrefix("/static/", fs))
	router.HandleFunc("/", h.Index)
	router.HandleFunc("/sign-up", h.SignUp)
	router.HandleFunc("/sign-in", h.SignIn)
	router.HandleFunc("/private-post", h.Private)
}
