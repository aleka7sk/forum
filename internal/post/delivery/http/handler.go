package http

import (
	"forum/internal/post"
	"forum/models"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Handler struct {
	usecase post.UseCase
}

type IsAuth struct {
	IsAuth bool
	// Data   interface{}
	User  *models.User
	Posts []models.Post
}

func NewHandler(usecase post.UseCase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func RenderTemplate(tmpl string) []string {
	files := []string{}
	files = append(files, "templates/"+tmpl)
	static_tmpl := []string{"templates/layout/layout.html"}
	files = append(files, static_tmpl...)
	return files
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Error"))
		return
	}
	right := r.Context().Value("rights")

	tmpl, err := template.ParseFiles(RenderTemplate("index.html")...)
	if err != nil {
		log.Printf("Error parse main page index: %v", tmpl)
	}

	posts := h.usecase.GetAllPosts(r.Context())
	data := IsAuth{}
	if right == "auth" {
		data = IsAuth{IsAuth: true, Posts: posts}
	} else {
		data = IsAuth{IsAuth: false, Posts: posts}
	}

	tmpl.ExecuteTemplate(w, "index.html", data)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "index.html", 200)
		return
	}
	c := &http.Cookie{
		Name:    "token",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, c)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	right := r.Context().Value("rights")
	if r.Method == "POST" {
		if right == "auth" {
			r.ParseForm()
			title := r.Form.Get("title")
			author := r.Form.Get("author")
			content := r.Form.Get("content")
			h.usecase.CreatePost(r.Context(), title, author, content)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	if right == "auth" {
		tmpl, err := template.ParseFiles(RenderTemplate("create_post.html")...)
		if err != nil {
			log.Printf("Error parse main page index: %v", tmpl)
		}
		tmpl.Execute(w, IsAuth{IsAuth: true})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Недоделанный ---> доделать
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	right := r.Context().Value("rights")
	if r.Method == "POST" {
		if right == "auth" {
			r.ParseForm()
			title := r.Form.Get("title")
			author := r.Form.Get("author")
			content := r.Form.Get("content")
			h.usecase.CreatePost(r.Context(), title, author, content)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	if right == "auth" {
		tmpl, err := template.ParseFiles(RenderTemplate("create_post.html")...)
		if err != nil {
			log.Printf("Error parse main page index: %v", tmpl)
		}
		tmpl.Execute(w, IsAuth{IsAuth: true})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
