package http

import (
	"forum/internal/post"
	"net/http"
)

func RegisterHTTPEndpoints(router *http.ServeMux, auc post.UseCase) {
	h := NewHandler(auc)
	router.HandleFunc("/post", h.Index)
}
