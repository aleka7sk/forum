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
	User       *models.User
	Posts      []postmodels.Post
	Post       postmodels.Post
	Vote       models.Vote
	Votes      []models.Vote
	Comments   []models.Comment
	Categories []models.Category
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
		log.Printf("Get votes error: %v", err)
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
	categories, err := h.usecase.GetCategoryName(r.Context())
	if err != nil {
		log.Printf("Get category name error: %v", err)
	}
	if r.Method == "POST" {
		if right != nil {
			r.ParseForm()
			category := r.Form.Get("category")
			category_id, err := strconv.Atoi(category)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			title := r.Form.Get("title")
			author := right.(middleware.UserInfo).Username
			content := r.Form.Get("content")
			author_id := right.(middleware.UserInfo).Id

			h.usecase.CreatePost(r.Context(), title, author, content, category_id, author_id)

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
		tmpl.Execute(w, IsAuth{IsAuth: true, Categories: categories})
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
	comments, err := h.usecase.GetComments(r.Context(), id_int)
	if err != nil {
		log.Printf("Get comments error: %v", err)
	}
	if r.Method == "POST" && right.(middleware.UserInfo).Rights {
		user_id := right.(middleware.UserInfo).Id
		// r.ParseForm()
		r.ParseForm()
		var vote string
		for key := range r.Form {
			vote = key
		}

		if vote == "like" {
			if err := h.usecase.CreateVote(r.Context(), id_int, user_id, 1); err != nil {
				log.Fatalf("Create vote error: %v", err)
			}
			post, err := h.usecase.GetPost(r.Context(), id_int, user_id)
			if err != nil {
				log.Printf("Get post error: %v", err)
			}
			tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: true, Post: post, Comments: comments})

		} else if vote == "dislike" {
			if err := h.usecase.CreateVote(r.Context(), id_int, user_id, 2); err != nil {
				log.Fatalf("Create vote error: %v", err)
			}
			post, err := h.usecase.GetPost(r.Context(), id_int, user_id)
			if err != nil {
				log.Printf("Get post error: %v", err)
			}
			tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: true, Post: post, Comments: comments})
		}
		return
	}

	if r.Method == "GET" && right.(middleware.UserInfo).Rights {
		user_id := right.(middleware.UserInfo).Id
		post, err := h.usecase.GetPost(r.Context(), id_int, user_id)
		if err != nil {
			log.Printf("Get post error: %v", err)
		}
		tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: true, Post: post, Comments: comments})
		return
	}
	post, err := h.usecase.GetPost(r.Context(), id_int, 0)
	if err != nil {
		log.Printf("Get post error noauth: %v", err)
	}
	tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: false, Post: post, Comments: comments})
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
	tmpl.ExecuteTemplate(w, "post.html", IsAuth{IsAuth: true, Posts: nil})
}

func (h *Handler) LikedPosts(w http.ResponseWriter, r *http.Request) {
	right := r.Context().Value("info")
	if !right.(middleware.UserInfo).Rights {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles(RenderTemplate("index.html")...)
	if err != nil {
		log.Printf("Error parse main page index: %v", tmpl)
	}

	posts, err := h.usecase.GetLikedPosts(r.Context(), right.(middleware.UserInfo).Id)
	if err != nil {
		log.Printf("Liked Posts error: %v", err)
	}

	tmpl.ExecuteTemplate(w, "index.html", IsAuth{IsAuth: true, Posts: posts})
}

func (h *Handler) UnlikedPosts(w http.ResponseWriter, r *http.Request) {
	right := r.Context().Value("info")
	if !right.(middleware.UserInfo).Rights {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl, err := template.ParseFiles(RenderTemplate("index.html")...)
	if err != nil {
		log.Printf("Error parse main page index: %v", tmpl)
	}

	posts, err := h.usecase.GetDislikedPosts(r.Context(), right.(middleware.UserInfo).Id)
	if err != nil {
		log.Printf("Liked Posts error: %v", err)
	}

	tmpl.ExecuteTemplate(w, "index.html", IsAuth{IsAuth: true, Posts: posts})
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("comment create")
	right := r.Context().Value("info")
	tmpl, err := template.ParseFiles(RenderTemplate("post.html")...)
	if err != nil {
		log.Fatalf("Error parse main page index: %v", tmpl)
	}

	id := strings.TrimPrefix(r.URL.Path, "/create-comment/")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == "POST" && right.(middleware.UserInfo).Rights {
		r.ParseForm()
		content := r.Form.Get("comment")
		if err := h.usecase.CreateComment(r.Context(), id_int, right.(middleware.UserInfo).Id, content); err != nil {
			log.Printf("Create comment error: %v", err)
		}
		http.Redirect(w, r, "/article/"+id, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
