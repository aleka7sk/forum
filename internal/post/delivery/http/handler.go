package http

import (
	"fmt"
	"forum/internal/middleware"
	"forum/internal/post"
	"forum/internal/postmodels"
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
	User     *models.User
	Posts    []postmodels.Post
	Post     postmodels.Post
	Emotion  models.Emotion
	Emotions []models.Emotion
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
	fmt.Println(posts)
	// emotions, err := h.usecase.GetEmotions(r.Context())
	if err != nil {
		log.Printf("Get emotions error: %v", err)
	}

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
	id_int, err := strconv.Atoi(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" && right.(middleware.UserInfo).Rights {
		user_id := right.(middleware.UserInfo).Id
		// r.ParseForm()
		r.ParseForm()
		var emotion string
		for key := range r.Form {
			emotion = key
		}
		if emotion == "like" {
			if err := h.usecase.CreateEmotion(r.Context(), id_int, user_id, true, false); err != nil {
				log.Fatalf("Create emotion error: %v", err)
			}
			post, err := h.usecase.GetPost(r.Context(), id_int, user_id)
			if err != nil {
				log.Printf("Get post error: %v", err)
			}
			tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: true, Post: post})

		} else if emotion == "dislike" {
			if err := h.usecase.CreateEmotion(r.Context(), id_int, user_id, false, true); err != nil {
				log.Fatalf("Create emotion error: %v", err)
			}
			post, err := h.usecase.GetPost(r.Context(), id_int, user_id)
			if err != nil {
				log.Printf("Get post error: %v", err)
			}
			tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: true, Post: post})

		}
		return
		// json.NewDecoder(r.Body).Decode(&groups)
	}

	if r.Method == "GET" && right.(middleware.UserInfo).Rights {
		user_id := right.(middleware.UserInfo).Id
		post, err := h.usecase.GetPost(r.Context(), id_int, user_id)
		if err != nil {
			log.Printf("Get post error: %v", err)
		}
		tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: true, Post: post})
		return
	}
	post, err := h.usecase.GetPost(r.Context(), id_int, 0)
	if err != nil {
		log.Printf("Get post error noauth: %v", err)
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
	fmt.Println(posts)
	tmpl.ExecuteTemplate(w, "my-post.html", IsAuth{IsAuth: true, Posts: nil})
}

func (h *Handler) LikedPosts(w http.ResponseWriter, r *http.Request) {
	right := r.Context().Value("info")
	if !right.(middleware.UserInfo).Rights {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles(RenderTemplate("post.html")...)
	if err != nil {
		log.Printf("Error parse main page index: %v", tmpl)
	}

	posts, err := h.usecase.GetLikedPosts(r.Context(), right.(middleware.UserInfo).Id)
	if err != nil {
		log.Printf("Liked Posts error: %v", err)
	}
	fmt.Println(posts)
	tmpl.ExecuteTemplate(w, "my-post.html", IsAuth{IsAuth: true, Posts: nil})
}

func (h *Handler) UnlikedPosts(w http.ResponseWriter, r *http.Request) {
	right := r.Context().Value("info")
	if !right.(middleware.UserInfo).Rights {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles(RenderTemplate("post.html")...)
	if err != nil {
		log.Printf("Error parse main page index: %v", tmpl)
	}

	posts, err := h.usecase.GetUnlikedPosts(r.Context(), right.(middleware.UserInfo).Id)
	if err != nil {
		log.Printf("Unliked Posts error: %v", err)
	}
	fmt.Println(posts)
	tmpl.ExecuteTemplate(w, "my-post.html", IsAuth{IsAuth: true, Posts: nil})
}
