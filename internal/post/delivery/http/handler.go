package http

import (
	"forum/internal/post"
	"log"
	"net/http"
	"text/template"
)

type Handler struct {
	usecase post.UseCase
}

func NewHandler(usecase post.UseCase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		log.Printf("%v", err)
	}
	data := h.usecase.GetAllPosts(r.Context())
	tmpl.Execute(w, data)
}
