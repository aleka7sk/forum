package http

import (
	"fmt"
	"forum/internal/middleware"
	"forum/internal/post"
	"forum/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type Handler struct {
	usecase post.UseCase
}

type IsAuth struct {
	IsAuth bool
	// Data   interface{}
	User    *models.User
	Posts   []models.Post
	Post    models.Post
	Emotion models.Emotion
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
	right := r.Context().Value("info")

	tmpl, err := template.ParseFiles(RenderTemplate("index.html")...)
	if err != nil {
		log.Fatalf("Error parse main page index: %v", tmpl)
	}

	posts := h.usecase.GetAllPosts(r.Context())
	data := IsAuth{}
	if right.(middleware.UserInfo).Rights {
		data = IsAuth{IsAuth: true, Posts: posts}
	} else {
		data = IsAuth{IsAuth: false, Posts: posts}
	}

	tmpl.ExecuteTemplate(w, "index.html", data)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	right := r.Context().Value("info")
	if r.Method == "POST" {
		if right != nil {
			r.ParseForm()
			title := r.Form.Get("title")
			author := right.(middleware.UserInfo).Username
			content := r.Form.Get("content")
			author_id := right.(middleware.UserInfo).Id
			fmt.Println(author_id)
			h.usecase.CreatePost(r.Context(), title, author, content, strconv.Itoa(author_id))

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	if right.(middleware.UserInfo).Rights {
		tmpl, err := template.ParseFiles(RenderTemplate("create_post.html")...)
		if err != nil {
			log.Printf("Error parse main page index: %v", tmpl)
		}
		tmpl.Execute(w, IsAuth{IsAuth: true})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	right := r.Context().Value("info")
	tmpl, err := template.ParseFiles(RenderTemplate("post.html")...)
	if err != nil {
		log.Fatalf("Error parse main page index: %v", tmpl)
	}

	id := strings.TrimPrefix(r.URL.Path, "/article/")

	post := h.usecase.GetPost(r.Context(), id)

	if r.Method == "POST" && right.(middleware.UserInfo).Rights {
		// r.ParseForm()
		r.ParseForm()
		var emotion string
		for key := range r.Form {
			emotion = key
		}
		if emotion == "like" {
			if err := h.usecase.CreateEmotion(r.Context(), id, right.(middleware.UserInfo).Id, true, false); err != nil {
				log.Fatalf("Create emotion error: %v", err)
			}
			tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: true, Post: post, Emotion: models.Emotion{Like: 1, Dislike: 0}})
		} else if emotion == "dislike" {
			if err := h.usecase.CreateEmotion(r.Context(), id, right.(middleware.UserInfo).Id, false, true); err != nil {
				log.Fatalf("Create emotion error: %v", err)
			}
			tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: true, Post: post, Emotion: models.Emotion{Like: 0, Dislike: 1}})
		}

		// json.NewDecoder(r.Body).Decode(&groups)
	} else if r.Method == "GET" && right.(middleware.UserInfo).Rights {
		tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: true, Post: post})
		return
	}

	tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: false, Post: post})
}

func (h *Handler) MyPosts(w http.ResponseWriter, r *http.Request) {
	right := r.Context().Value("info")
	if !right.(middleware.UserInfo).Rights {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles(RenderTemplate("my-post.html")...)
	if err != nil {
		log.Printf("Error parse main page index: %v", tmpl)
	}
	fmt.Println(strconv.Itoa(right.(middleware.UserInfo).Id))
	posts := h.usecase.GetMyPosts(r.Context(), strconv.Itoa(right.(middleware.UserInfo).Id))

	tmpl.ExecuteTemplate(w, "my-post.html", IsAuth{IsAuth: true, Posts: posts})
}
