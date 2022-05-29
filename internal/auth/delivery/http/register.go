package http

import (
	"forum/internal/auth"
	"forum/internal/middleware"
	"net/http"
)

func RegisterHTTPEndpoints(router *http.ServeMux, auc auth.UseCase) {
	h := NewHandler(auc)
	router.HandleFunc("/sign-up", middleware.Handle(http.HandlerFunc(h.SignUp)))
	router.HandleFunc("/sign-in", middleware.Handle(http.HandlerFunc(h.SignIn)))
	// router.HandleFunc("/private", middleware.Handle(http.HandlerFunc(h.Private)))
}
